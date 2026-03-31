package main

import (
	"context"
	"fmt"

	"github.com/ksvaza/ees-link/db"
	"github.com/ksvaza/ees-link/envreader"
	"github.com/ksvaza/ees-link/httpapi"
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

	ctx := context.Background()

	dbURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		envreader.GetEnvString("POSTGRES_HOST"),
		envreader.GetEnvString("POSTGRES_PORT"),
		envreader.GetEnvString("POSTGRES_USER"),
		envreader.GetEnvString("POSTGRES_PASSWORD"),
		envreader.GetEnvString("POSTGRES_DB"))

	e := db.MigrateUp(ctx, dbURL)
	if e != nil {
		logrus.WithError(err).Error("Failed to migrate up")
		return
	}

	logrus.Info("\nSveika, http aplikācija!\n")

	err = httpapi.SetupHTTPAPI()
	if err != nil {
		logrus.WithError(errors.Wrap(err, "HTTP")).Error("Error")
	}

	for {

	}
}
