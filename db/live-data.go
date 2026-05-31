package db

import (
	"context"
	"math"
	"time"

	"github.com/ksvaza/ees-link/models"
)

var (
	localLiveData map[int]models.LiveRaceData = make(map[int]models.LiveRaceData)
)

func (db *RealDB) GetLiveRaceData(ctx context.Context) ([]models.LiveRaceData, error) {
	var arr []models.LiveRaceData = make([]models.LiveRaceData, 0)
	for _, e := range localLiveData {
		arr = append(arr, e)
	}
	return arr, nil
}

func (DB *RealDB) UpdateLiveData(ctx context.Context, telemetry models.CarTelemetry) error {
	carID := telemetry.ID

	val, err := DB.GetAllCarParameters(ctx)
	if err != nil {
		return err
	}

	var params models.CarParameters
	for _, v := range val {
		if v.CarID == carID {
			params = v
		}
	}

	//DB.GetApplicationByTeamName()
	lrd := models.LiveRaceData{
		Key:          "",
		ID:           telemetry.ID,
		Username:     params.TeamName,
		Avatar:       params.Avatar,
		Status:       "online",
		Position:     0,
		Lat:          float64(telemetry.GPSData.Latitude),
		Lon:          float64(telemetry.GPSData.Longitute),
		Spd:          float32(telemetry.GPSData.Speed),
		Power:        float32(telemetry.PSUData.PowerOut) / float32(100),
		Acceleration: float32(math.Sqrt(float64(telemetry.AccelData.X)*float64(telemetry.AccelData.X) + float64(telemetry.AccelData.Y)*float64(telemetry.AccelData.Y) + float64(telemetry.AccelData.Z)*float64(telemetry.AccelData.Z))),
		Voltage:      float32(telemetry.PSUData.VoltageOut) / float32(100),
		UpdatedAt:    time.Now(),
	}

	localLiveData[carID] = lrd
	return nil
}
