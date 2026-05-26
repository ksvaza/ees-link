package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/ksvaza/ees-link/db"
	"github.com/ksvaza/ees-link/envreader"
	"github.com/ksvaza/ees-link/httpapi"
	"github.com/ksvaza/ees-link/logeris"
	"github.com/ksvaza/ees-link/models"
	"github.com/ksvaza/ees-link/myqtt"
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

	RealDB := db.RealDB{}

	RealDB.TestHealthiness(ctx)

	// DB impregnācija ar testu datiem

	if envreader.GetEnvBool("CLEARSEEDDATA") {
		logrus.Info("Clearing seeded test data from database...")
		err = db.DeleteSeedData(ctx)
		if err != nil {
			logrus.WithError(errors.Wrap(err, "DB Clear Seed Data")).Error("Error")
			return
		}
	}

	// visadministratora atgūšana no vides mainīgajiem
	var superadmin models.AdminAccount
	superadmin.Username = envreader.GetEnvString("SUPERADMIN_USERNAME")
	superadmin.Password = envreader.GetEnvString("SUPERADMIN_PASSWORD")
	superadmin.Superadmin = true
	//logrus.Infof("Superadmin credentials from environment: username='%s', password='%s'", superadmin.Username, superadmin.Password)
	err = httpapi.RegisterAdminAccount(superadmin)
	if err != nil {
		logrus.WithError(err).Error("Failed to register superadmin account")
		return
	}
	// ----

	if envreader.GetEnvBool("SEEDDATABASE") {
		logrus.Info("Seeding database with test data...")
		err = db.SeedTestData(ctx)
		if err != nil {
			logrus.WithError(errors.Wrap(err, "DB Seed")).Error("Error")
			return
		}
	}

	// MQTT vides mainīgo iegūšana un inicializācija
	mqttConfig := models.MqttConfig{
		Host:     envreader.GetEnvString("MQTT_HOST"),
		Port:     envreader.GetEnvInt("MQTT_PORT"),
		Username: envreader.GetEnvString("MQTT_USERNAME"),
		Password: envreader.GetEnvString("MQTT_PASSWORD"),
	}
	myqtt.InitMQTT(mqttConfig)
	// ----

	// wg un pavedienu uzdevumu inicializācija

	wg := &sync.WaitGroup{}

	myqtt.StartMQTTHost(ctx, wg)

	// ----

	logrus.Info("\nSveika, http aplikācija!\n")

	err = httpapi.SetupHTTPAPI()
	if err != nil {
		logrus.WithError(errors.Wrap(err, "HTTP")).Error("Error")
	}

	for {

	}
}
