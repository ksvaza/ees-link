package myqtt

import (
	"github.com/ksvaza/ees-link/models"
	"github.com/pkg/errors"
)

func SendMQTTMessage(message models.MqttMessage) error {
	if mqttClient == nil || !mqttClient.IsConnected() {
		return errors.New("MQTT client not connected")
	}

	token := mqttClient.Publish(message.Topic, 0, false, message.Payload)
	token.Wait()
	return token.Error()
}
