package db

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Event struct
type Event struct {
	ID         int
	InstanceID int
	Name       string
	Value      float64
	Timestamp  time.Time
}

// DB holds the connection pool
var Pool *pgxpool.Pool

// Init connects to Postgres and prepares the pool
func Init(dbURL string) error {
	logrus.Info("Connecting to DB...")
	ctx := context.Background()
	var err error
	Pool, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		logrus.WithError(err).Error("Failed to create DB pool")
		return err
	}

	// Create table if not exists

	table_query := []string{`CREATE TABLE IF NOT EXISTS events (
        id SERIAL PRIMARY KEY,
        instance_id INT NOT NULL,
        name TEXT NOT NULL,
        value DOUBLE PRECISION NOT NULL,
        timestamp TIMESTAMPTZ NOT NULL
    );`, `CREATE TABLE IF NOT EXISTS dalibnieki (
        id SERIAL PRIMARY KEY,
        instance_id INT NOT NULL,
        name TEXT NOT NULL,
        value DOUBLE PRECISION NOT NULL,
        timestamp TIMESTAMPTZ NOT NULL
    );`,
	}

	for _, q := range table_query {
		_, err := Pool.Exec(ctx, q)
		if err != nil {
			logrus.WithError(err).Error("Failed to create tables")
			return err
		}
	}

	return nil
}

// InsertEvent inserts a single Event and sets its ID
func InsertEvent(e *Event) error {
	ctx := context.Background()
	insertSQL := `
    INSERT INTO events (instance_id, name, value, timestamp)
    VALUES ($1, $2, $3, $4)
    RETURNING id`

	return Pool.QueryRow(ctx, insertSQL, e.InstanceID, e.Name, e.Value, e.Timestamp).Scan(&e.ID)
}

func InsertEventsBatchToTable(events []*Event, table string) (error) {
	ctx := context.Background()
	fmt.Printf("asaaa\n")
	tx, err := Pool.Begin(ctx)

	if err != nil {

		logrus.WithError(err).Error("Failed to begin transaction")
		return err
	}

	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}

	for _, e := range events {
		batch.Queue(
			"INSERT INTO "+table+" (instance_id, name, value, timestamp) VALUES ($1, $2, $3, $4) RETURNING id",
			e.InstanceID, e.Name, e.Value, e.Timestamp,
		)
	}

	br := tx.SendBatch(ctx, batch)
	for i, _ := range events {
		err := br.QueryRow().Scan(&events[i].ID)
		if err != nil {
			logrus.WithError(err).Error("Failed to execute batch insert")
			br.Close()
			return err
		}
	}
	br.Close()
	return tx.Commit(ctx)
}

func GetAllEventsByName(name string) ([]*Event, error) {
	ctx := context.Background()

	rows, err := Pool.Query(ctx, `
		SELECT id, name, instance_id, value, timestamp
		FROM dalibnieki
		WHERE name = $1
	`, name)
	if err != nil {
		logrus.WithError(err).Error("Failed to query events")
		return nil, err
	}
	defer rows.Close()

	var events []*Event

	for rows.Next() {
		var e Event
		err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.InstanceID,
			&e.Value,
			&e.Timestamp,
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to scan event row")
			return nil, err
		}
		events = append(events, &e)
	}

	return events, nil
}

// Close closes the DB pool
func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
