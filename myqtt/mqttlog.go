package myqtt

import (
	"context"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ksvaza/ees-link/data"
	"github.com/ksvaza/ees-link/db"
	"github.com/ksvaza/ees-link/models"
)

func mqttGeneralHandler(c mqtt.Client, m mqtt.Message) {
	entry := models.MqttLogEntry{
		Message: models.MqttMessage{
			Topic:   m.Topic(),
			Payload: m.Payload(),
		},
		ReceivedAt: time.Now(),
	}
	var realDB data.Database = &db.RealDB{}
	realDB.SaveMQTTLog(context.Background(), &entry)
}
