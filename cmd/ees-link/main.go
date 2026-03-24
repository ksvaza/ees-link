package main

import (
	"fmt"
	"time"

	"github.com/ksvaza/ees-link/db"
	"github.com/ksvaza/ees-link/envreader"
	"github.com/ksvaza/ees-link/logeris"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithError(errors.Errorf("%v", r)).Error("Panic")
		}
	}()

	err := envreader.SetupEnvreader(".env")
	if err != nil {
		logrus.WithError(err).Error("Failed to setup environment reader")
		return
	}

	logfile := envreader.GetEnvString("LOGFILE")
	clearlog := envreader.GetEnvBool("CLEARLOG")
	f, err := logeris.SetupLogger(logfile, clearlog)
	if err != nil {
		logrus.WithError(err).Error("Failed to setup logger")
		return
	}
	defer f.Close()

	logrus.Info("\nSveika, pasaule!\n")

	for i := 0; i < 5; i++ {
		logrus.Tracef("Trace %d", i)
		logrus.Debugf("Debug %d", i)
		logrus.Infof("Info %d", i)
		logrus.Warnf("Warn %d", i)
		logrus.Errorf("Error %d", i)
	}

	// sūdu aispildīšana

	// fmt.Printf("karukskigi\n")
	/*
		var events []*db.Event
		numEvents := 1000 // number of events to generate

		for i := 0; i < numEvents; i++ {
			e := &db.Event{
				InstanceID: rand.Intn(30) + 1, // random instance 1–30
				Name:       fmt.Sprintf("metric_%d", rand.Intn(1000)),
				Value:      rand.Float64() * 100.0, // random float 0–100
				Timestamp:  time.Now().Add(time.Duration(rand.Intn(1000)) * time.Millisecond),
			}
			events = append(events, e)
		}
	*/

	events := []*db.Event{
		{
			InstanceID: 1,
			Name:       "temp",
			Value:      23.5,
			Timestamp:  time.Now(),
		},
		{
			InstanceID: 2,
			Name:       "humidity",
			Value:      65.0,
			Timestamp:  time.Now(),
		},
	}

	db.Init("postgres://postgres:postgres@localhost:8086/demo?sslmode=disable")
	err = db.InsertEventsBatchToTable(events, "dalibnieki")
	if err != nil {
		logrus.WithError(err).Error("Failed to insert events batch")
		return
	}

	dabut, err := db.GetAllEventsByName("temp")
	if err != nil {
		logrus.WithError(err).Error("Failed to get all events")
		return
	}
	for _, e := range dabut {
		fmt.Printf("Event: ID=%d, InstanceID=%d, Name=%s, Value=%.2f, Timestamp=%s\n",
			e.ID, e.InstanceID, e.Name, e.Value, e.Timestamp.Format(time.RFC3339))
	}

	/*

		err1 := db.Init("postgres://postgres:postgres@localhost:8086/demo?sslmode=disable")
		if err1 != nil {
			log.Fatal(err1)
		}
		defer db.Close()

		e := db.Event{InstanceID: 1, Name: "temp", Value: 23.5, Timestamp: time.Now()}
		err = db.InsertEvent(&e)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Inserted event: %+v\n", e)

	*/

	for {
		//fmt.Printf("nig")
		//time.Sleep(2000) // 2000ms

	}
}
