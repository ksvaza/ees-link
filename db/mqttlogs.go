package db

import (
	"context"

	"github.com/ksvaza/ees-link/models"
	"github.com/sirupsen/logrus"
)

func (db *RealDB) SaveMQTTLog(ctx context.Context, log *models.MqttMessage) error {
	_, err := Pool.Exec(ctx, MQTTLogSaveRequest,
		log.Topic,
		log.Payload,
	)
	if err != nil {
		logrus.WithError(err).Error("Failed to save MQTT log")
		return err
	}
	return nil
}

func (db *RealDB) GetMQTTLogs(ctx context.Context, limit int) ([]models.MqttMessage, error) {
	rows, err := Pool.Query(ctx, MQTTLogsReadRequest, limit)
	if err != nil {
		logrus.WithError(err).Error("Failed to query MQTT logs table")
		return nil, err
	}
	defer rows.Close()

	var logs []models.MqttMessage
	for rows.Next() {
		var log models.MqttMessage
		err = rows.Scan(
			&log.Topic,
			&log.Payload,
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to scan MQTT log row")
			return nil, err
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		logrus.WithError(err).Error("Failed iterating MQTT log rows")
		return nil, err
	}

	if len(logs) == 0 {
		return nil, nil
	}

	return logs, nil
}
