package myqtt

import (
	"context"
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ksvaza/ees-link/data"
	"github.com/ksvaza/ees-link/db"
	"github.com/ksvaza/ees-link/models"
	"github.com/sirupsen/logrus"
)

func mqttCarTelemetryHandler(c mqtt.Client, m mqtt.Message) {
	logrus.Infof("Received car telemetry data: %+v", m)

	var telemetry models.CarTelemetry
	if err := json.Unmarshal(m.Payload(), &telemetry); err != nil {
		logrus.WithError(err).Error("Failed to unmarshal received car telemetry data")
	}

	logrus.Debugf("Car telemetry data structure %+v", telemetry)

	// Update livedata
	//websockets.UpdateLiveData(telemetry)

	var realDB data.Database = &db.RealDB{}
	if err := realDB.SaveCarTelemetry(context.Background(), telemetry); err != nil {
		logrus.WithError(err).Error("Failed to save car telemetry data to database")
	}
}
