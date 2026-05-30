package db

import (
	"context"

	"github.com/ksvaza/ees-link/models"
	"github.com/sirupsen/logrus"
)

func (db *RealDB) SaveMQTTLog(ctx context.Context, log *models.MqttLogEntry) error {
	_, err := Pool.Exec(ctx, MQTTLogSaveRequest,
		log.Message.Topic,
		log.Message.Payload,
		log.ReceivedAt,
	)
	if err != nil {
		logrus.WithError(err).Error("Failed to save MQTT log")
		return err
	}
	return nil
}

func (db *RealDB) GetMQTTLogs(ctx context.Context, limit int) ([]models.MqttLogEntry, error) {
	rows, err := Pool.Query(ctx, MQTTLogsReadRequest, limit)
	if err != nil {
		logrus.WithError(err).Error("Failed to query MQTT logs table")
		return nil, err
	}
	defer rows.Close()

	var logs []models.MqttLogEntry
	for rows.Next() {
		var log models.MqttLogEntry
		err = rows.Scan(
			&log.Message.Topic,
			&log.Message.Payload,
			&log.ReceivedAt,
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
