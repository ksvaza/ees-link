package db

import (
	"context"
	"encoding/json"

	"github.com/ksvaza/ees-link/models"
	"github.com/sirupsen/logrus"
)

func (DB *RealDB) SavePoints(ctx context.Context, points models.Points) error {
	// check whether points have already been saved for the given race_name(?, if yes then update instead of insert?)
	points_data, err := json.Marshal(points.Points)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal points data")
		return err
	}
	_, err = Pool.Exec(ctx, PointsSaveRequest,
		points.RaceName,
		points_data,
	)
	if err != nil {
		logrus.WithError(err).Error("Failed to save points")
		return err
	}
	return nil
}

func (DB *RealDB) GetAllPoints(ctx context.Context) ([]models.Points, error) {
	rows, err := Pool.Query(ctx, PointsReadAllRequest)
	if err != nil {
		logrus.WithError(err).Error("Failed to query points table")
		return nil, err
	}
	defer rows.Close()

	var allPoints []models.Points
	for rows.Next() {
		var points models.Points
		var pointsData []byte
		err = rows.Scan(&points.RaceName, &pointsData)
		if err != nil {
			logrus.WithError(err).Error("Failed to scan points row")
			return nil, err
		}
		err = json.Unmarshal(pointsData, &points.Points)
		if err != nil {
			logrus.WithError(err).Error("Failed to unmarshal points data")
			return nil, err
		}
		allPoints = append(allPoints, points)
	}

	if err = rows.Err(); err != nil {
		logrus.WithError(err).Error("Failed iterating points rows")
		return nil, err
	}

	if len(allPoints) == 0 {
		return nil, nil
	}

	return allPoints, nil
}

func (DB *RealDB) GetPointsByRaceName(ctx context.Context, raceName string) (*models.Points, error) {
	var points models.Points
	var pointsData []byte
	err := Pool.QueryRow(ctx, PointsReadByRaceNameRequest, raceName).Scan(&points.RaceName, &pointsData)
	if err != nil {
		logrus.WithError(err).Error("Failed to query points by race name")
		return nil, err
	}
	err = json.Unmarshal(pointsData, &points.Points)
	if err != nil {
		logrus.WithError(err).Error("Failed to unmarshal points data")
		return nil, err
	}
	return &points, nil
}

func (DB *RealDB) UpdatePointsByRaceName(ctx context.Context, raceName string, points models.Points) error {
	points_data, err := json.Marshal(points.Points)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal points data")
		return err
	}
	_, err = Pool.Exec(ctx, PointsUpdateByRaceNameRequest,
		points_data,
		raceName,
	)
	if err != nil {
		logrus.WithError(err).Error("Failed to update points")
		return err
	}
	return nil
}
