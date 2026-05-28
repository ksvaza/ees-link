package websockets

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/ksvaza/ees-link/data"
	"github.com/ksvaza/ees-link/db"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	writeWait      = 10 * time.Second    // max time to write one message
	pongWait       = 60 * time.Second    // must hear from client within this
	pingPeriod     = (pongWait * 9) / 10 // must be < pongWait
	broadcastEvery = 500 * time.Millisecond
	sendBuffer     = 16
)

var upgrader = websocket.Upgrader{}

type client struct {
	conn   *websocket.Conn
	send   chan []byte
	closed chan struct{}
	once   sync.Once
}

func (c *client) shutdown() {
	c.once.Do(func() {
		close(c.closed)
		_ = c.conn.Close()
	})
}

var (
	connectionMutex sync.RWMutex
	clients         = map[*client]struct{}{}
)

func removeClient(c *client) {
	logrus.Debugf("Removing websocket client, currently %d clients", len(clients))
	connectionMutex.Lock()
	delete(clients, c)
	connectionMutex.Unlock()
	c.shutdown()
}

func StartWebSocketServer() {
	var realDB data.Database = &db.RealDB{}

	go func() {
		ticker := time.NewTicker(broadcastEvery)
		defer ticker.Stop()

		for range ticker.C {
			liveData, err := realDB.GetLiveRaceData(context.Background())
			if err != nil {
				logrus.WithError(err).Error("Failed to get live race data from database")
				continue
			}
			msg, err := json.Marshal(liveData)
			if err != nil {
				logrus.WithError(err).Error("Failed to marshal live race data")
				continue
			}

			// snapshot under lock, send outside lock
			connectionMutex.RLock()
			snapshot := make([]*client, 0, len(clients))
			for c := range clients {
				snapshot = append(snapshot, c)
			}
			connectionMutex.RUnlock()

			for _, c := range snapshot {
				select {
				case c.send <- msg:
				default:
					// buffer full => client can't keep up, drop it
					logrus.Warn("Dropping slow websocket client")
					removeClient(c)
				}
			}
		}
	}()
}

// writePump writes queued messages and emits pings on a timer.
func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		removeClient(c)
	}()

	for {
		select {
		case <-c.closed:
			return
		case msg := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				logrus.WithError(err).Debug("write failed, dropping client")
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logrus.WithError(err).Debug("ping failed, dropping client")
				return
			}
		}
	}
}

// readPump's real job is detecting inactivity. The pong handler bumps the
// read deadline; if no pong arrives within pongWait, ReadMessage errors
// and we tear the client down.
func (c *client) readPump() {
	defer removeClient(c)

	c.conn.SetReadLimit(1024)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			return
		}
	}
}

func UpgradeWebSocket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	logrus.Infof("New websocket connection, current clients: %d", len(clients)+1)
	if r.Method != http.MethodGet {
		return errors.New("method not allowed")
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.WithError(err).Error("Error")
		return errors.Wrap(err, "Failed to upgrade websocket connection :/")
	}

	c := &client{
		conn:   conn,
		send:   make(chan []byte, sendBuffer),
		closed: make(chan struct{}),
	}

	connectionMutex.Lock()
	clients[c] = struct{}{}
	connectionMutex.Unlock()

	go c.writePump()
	go c.readPump()

	return nil
}
