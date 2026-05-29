package myqtt

import (
	"context"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ksvaza/ees-link/logeris"
	"github.com/ksvaza/ees-link/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	mqttClient   mqtt.Client
	mqttConfig   models.MqttConfig
	mqttInstance = "ees-link"
	mqttUrl      = ""
	readyForInit = false
)

func InitMQTT(config models.MqttConfig) {
	mqttConfig = config
	mqttUrl = fmt.Sprintf("tcp://%s:%d", mqttConfig.Host, mqttConfig.Port)
	mqttInstance = fmt.Sprintf("%s-%d", mqttInstance, time.Now().Unix())
	readyForInit = true
}

// MQTT client
func StartMQTTHost(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				logrus.WithError(errors.New(fmt.Sprintf("%v", r))).Error("Panic")
			}
		}()

		for !readyForInit {
			time.Sleep(time.Millisecond * 100)
		}

		logentry := logrus.WithField("component", "MQTT")
		logentry.Info("Starting MQTT host")

		mqtt.CRITICAL = logeris.NewLogrusLogger(logentry, logrus.FatalLevel)
		mqtt.ERROR = logeris.NewLogrusLogger(logentry, logrus.ErrorLevel)
		mqtt.WARN = logeris.NewLogrusLogger(logentry, logrus.WarnLevel)
		mqtt.DEBUG = logeris.NewLogrusLogger(logentry, logrus.DebugLevel)

		mqttClient = createMqttClient(ctx, mqttConfig)
		if mqttClient != nil {
			subscribeMqtt(ctx)
			logrus.Info("MQTT client connected")

			<-ctx.Done()
			mqttClient.Disconnect(uint(2 * time.Second))

			logrus.Info("MQTT client stopped")
		}
	}()
}

func createMqttClient(ctx context.Context, config models.MqttConfig) mqtt.Client {
	for {
		logrus.Debugf("MQTT client (%s) connecting to %s", mqttInstance, mqttUrl)

		opts := mqtt.NewClientOptions()
		opts.AddBroker(mqttUrl)
		if config.Username != "" && config.Password != "" {
			opts.SetUsername(config.Username)
			opts.SetPassword(config.Password)
		}
		opts.SetClientID(mqttInstance)
		opts.AutoReconnect = true
		opts.ResumeSubs = true
		opts.Order = false // Message receive order is not important. Disable feature to improve performance
		opts.OnConnect = func(c mqtt.Client) {
			logrus.Debugf("MQTT client (%s) connected to %s", mqttInstance, mqttUrl)
		}
		opts.OnConnectionLost = func(c mqtt.Client, err error) {
			logrus.WithError(errors.Wrap(err, "MQTT")).Debug("MQTT connection lost")
		}

		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			logrus.WithError(errors.Wrap(token.Error(), "MQTT")).Error("Can not connect to MQTT broker")

			sleep, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			<-sleep.Done()

			if sleep.Err() != nil && sleep.Err() == context.DeadlineExceeded {
				continue
			}
			return nil
		} else {
			return client
		}
	}
}

func subscribeMqtt(ctx context.Context) {
	roottopic := "ees-link/"

	// Vispārējais apstrādātājs logošanai
	token := mqttClient.Subscribe(roottopic+"#", 1, mqttGeneralHandler)
	token.Wait()
	if err := token.Error(); err != nil {
		logrus.WithError(errors.Wrap(err, "MQTT")).Error("Subscribe error")
	}
}
