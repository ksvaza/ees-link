package myqtt

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ksvaza/ees-link/data"
	"github.com/ksvaza/ees-link/db"
	"github.com/ksvaza/ees-link/models"
)

var realDB data.Database = &db.RealDB{}

func mqttGeneralHandler(c mqtt.Client, m mqtt.Message) {
	entry := models.MqttMessage{
		Topic:   m.Topic(),
		Payload: m.Payload(),
	}
	realDB.SaveMQTTLog(context.Background(), &entry)
}
