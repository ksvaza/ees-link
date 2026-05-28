package models

import "time"

type MqttConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type MqttMessage struct {
	Topic   string `json:"topic"`
	Payload []byte `json:"payload"`
}

type MqttLogEntry struct {
	Message    MqttMessage `json:"message"`
	ReceivedAt time.Time   `json:"received_at"`
}
