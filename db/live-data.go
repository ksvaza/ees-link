package db

import (
	"context"

	"github.com/ksvaza/ees-link/models"
)

func (db *RealDB) GetLiveRaceData(ctx context.Context) ([]models.LiveRaceData, error) {
	return []models.LiveRaceData{
		{
			Key:          "1",
			ID:           1,
			Username:     "team1",
			Avatar:       "https://www.gravatar.com/avatar/00000000000000000000000000000000?d=mp&f=y",
			Status:       "online",
			Position:     1,
			Lat:          56.9496,
			Lon:          24.1052,
			Spd:          10.5,
			Power:        100.0,
			Acceleration: 2.5,
			Voltage:      12.0,
			UpdatedAt:    "2024-06-01T12:00:00Z",
		},
		{
			Key:          "2",
			ID:           2,
			Username:     "team2",
			Avatar:       "https://www.gravatar.com/avatar/00000000000000000000000000000000?d=mp&f=y",
			Status:       "offline",
			Position:     3,
			Lat:          56.9501,
			Lon:          24.1060,
			Spd:          9.8,
			Power:        98.7,
			Acceleration: 2.3,
			Voltage:      11.9,
			UpdatedAt:    "2024-06-01T12:05:00Z",
		},
	}, nil
}
