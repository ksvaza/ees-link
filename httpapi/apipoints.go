package httpapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ksvaza/ees-link/data"
	"github.com/ksvaza/ees-link/fakedb"
	"github.com/ksvaza/ees-link/websockets"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func errorHandler(w http.ResponseWriter, err error, code int) {
	logrus.WithError(err).Error("Error")
	//http.Error(w, err.Error(), code)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(`{"error": "notika kļūme"}`))
}

type httpResult struct {
	ResponseType int
	Headers      map[string]string
	Body         interface{}
}

func test(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	if r.Method != http.MethodGet {
		return nil, errors.New("method not allowed")
	}

	//account := GetAccount(r.Context())
	//admin := GetAdminAccount(r.Context())
	// if account != nil {
	// 	logrus.Infof("Authenticated account: %+v", account)
	// }
	// if admin != nil {
	// 	logrus.Infof("Authenticated admin account: %+v", admin)
	// }
	// if account.Cilveks.Role != "admin" {
	// 	return nil, errors.New("forbidden")
	// }

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         "Sveika, pasaule!",
	}, nil

}

func Handler(fn func(r *http.Request, ps httprouter.Params) (*httpResult, error)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		result, err := fn(Authenticate(r), ps)
		if err != nil {
			errorHandler(w, err, http.StatusInternalServerError)
			return
		}

		var body []byte
		body, err = json.Marshal(result.Body)
		if err != nil {
			errorHandler(w, err, http.StatusInternalServerError)
			return
		}

		for k, v := range result.Headers {
			w.Header().Set(k, v)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(result.ResponseType)
		w.Write(body)
	}
}

// Sveika pasaule tipa situācija
func TestPointTestAPI(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("TestPointTestAPI called %+v", ps)

	if r.Method != http.MethodGet {
		errorHandler(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	// Read the request body and return it as the response
	_, err := io.ReadAll(r.Body)
	if err != nil {
		errorHandler(w, errors.Wrap(err, "failed to read request body"), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Hello, this is a test API!</h1>\n"))
}

func TestDatabase(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	logrus.Infof("TestPointTestAPI called %+v", ps)

	if r.Method != http.MethodGet {
		return nil, errors.New("method not allowed")
	}

	var realDB data.Database
	realDB = &fakedb.FSDatabase{}
	err := realDB.TestHealthiness(r.Context())
	if err != nil {
		return nil, errors.Wrap(err, "failed to test database health")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         "Database is healthy",
	}, nil
}

// tīri tests formu saņemšanai, vēlāk varētu būt noderīgi, lai saņemtu datus no frontend formas
func TestPointReceiveForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("TestPointReceiveForm called %+v", ps)

	fmt.Printf("Received form data: %+v", ps)

	if r.Method != http.MethodPost {
		errorHandler(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorHandler(w, err, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	logrus.Infof("Received form data: %s", string(body))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))

}

func sendNotImplemented(w http.ResponseWriter) {
	errorHandler(w, errors.New("not implemented"), http.StatusNotImplemented)
}

func requireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		errorHandler(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
		return false
	}
	return true
}

// --------------------------------------------------------------------------------------------------------------------------------

// --------------------------------------------------------------------------------------------------------------------------------

// User
// ----

func PointGetUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointGetUsers called %+v", ps)
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	// TODO: implement get users
	sendNotImplemented(w)
}

// --------------------------------------------------------------------------------------------------------------------------------

// Competitors
// -----------

// --------------------------------------------------------------------------------------------------------------------------------

// Races
// -----

func PointRaceStart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointRaceStart called %+v", ps)
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	// TODO: implement race start
	sendNotImplemented(w)
}

func PointCarFinish(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointCarFinish called %+v", ps)
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	// TODO: implement car finish
	sendNotImplemented(w)
}

func PointGetRaces(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointGetRaces called %+v", ps)
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	// TODO: implement get races
	sendNotImplemented(w)
}

func PointPostRaces(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointPostRaces called %+v", ps)
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	// TODO: implement post races
	sendNotImplemented(w)
}

func PointGetRaceConfig(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointGetRaceConfig called %+v", ps)
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	// TODO: implement get race config
	sendNotImplemented(w)
}

func PointGetCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointGetCategories called %+v", ps)
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	// TODO: implement get categories
	sendNotImplemented(w)
}

// --------------------------------------------------------------------------------------------------------------------------------

// Results
// -------

func PointGetRaceResults(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointGetRaceResults called %+v", ps)
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	// TODO: implement get race results
	sendNotImplemented(w)
}

// --------------------------------------------------------------------------------------------------------------------------------

// Points
// ------

func PointPostPoints(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointPostPoints called %+v", ps)
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	// TODO: implement post points
	sendNotImplemented(w)
}

func PointDeletePoints(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointDeletePoints called %+v", ps)
	if !requireMethod(w, r, http.MethodDelete) {
		return
	}
	// TODO: implement delete points
	sendNotImplemented(w)
}

// --------------------------------------------------------------------------------------------------------------------------------

// Leaderboard
// -----------

func PointGetLeaderboard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointGetLeaderboard called %+v", ps)
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	// TODO: implement get leaderboard
	sendNotImplemented(w)
}

// --------------------------------------------------------------------------------------------------------------------------------

// Admin
// -----

// --------------------------------------------------------------------------------------------------------------------------------

// News
// ----

func PointGetNews(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointGetNews called %+v", ps)
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	// TODO: implement get news
	sendNotImplemented(w)
}

func PointPostNews(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointPostNews called %+v", ps)
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	// TODO: implement post news
	sendNotImplemented(w)
}

func PointGetNewsByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointGetNewsByID called %+v", ps)
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	// TODO: implement get news by id
	sendNotImplemented(w)
}

func PointPutNewsByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointPutNewsByID called %+v", ps)
	if !requireMethod(w, r, http.MethodPut) {
		return
	}
	// TODO: implement put news by id
	sendNotImplemented(w)
}

func PointDeleteNewsByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointDeleteNewsByID called %+v", ps)
	if !requireMethod(w, r, http.MethodDelete) {
		return
	}
	// TODO: implement delete news by id
	sendNotImplemented(w)
}

func PointPostNewsAttachments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointPostNewsAttachments called %+v", ps)
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	// TODO: implement post news attachments
	sendNotImplemented(w)
}

// --------------------------------------------------------------------------------------------------------------------------------

// Events
// ------

func PointGetEvents(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointGetEvents called %+v", ps)
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	// TODO: implement get events
	sendNotImplemented(w)
}

func PointPostEvents(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointPostEvents called %+v", ps)
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	// TODO: implement post events
	sendNotImplemented(w)
}

func PointGetEventByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointGetEventByID called %+v", ps)
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	// TODO: implement get event by id
	sendNotImplemented(w)
}

func PointPutEventByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointPutEventByID called %+v", ps)
	if !requireMethod(w, r, http.MethodPut) {
		return
	}
	// TODO: implement put event by id
	sendNotImplemented(w)
}

func PointDeleteEventByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logrus.Infof("PointDeleteEventByID called %+v", ps)
	if !requireMethod(w, r, http.MethodDelete) {
		return
	}
	// TODO: implement delete event by id
	sendNotImplemented(w)
}

// --------------------------------------------------------------------------------------------------------------------------------

// --------------------------------------------------------------------------------------------------------------------------------

// WebSocket
// ---------

func PointWebSocket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := websockets.UpgradeWebSocket(w, r, ps)
	if err != nil {
		logrus.Errorf("Error upgrading WebSocket: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "notika kļūme"}`))
	}
}
