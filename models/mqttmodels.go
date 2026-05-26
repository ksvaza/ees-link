package models

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
