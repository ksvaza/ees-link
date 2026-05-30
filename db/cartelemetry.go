package db

import (
	"context"
	"time"

	"github.com/ksvaza/ees-link/models"
	"github.com/sirupsen/logrus"
)

func (DB *RealDB) SaveCarTelemetry(ctx context.Context, telemetry models.CarTelemetry) error {
	_, err := Pool.Exec(ctx, CarTelemetrySaveRequest,
		telemetry.ID,
		telemetry.RSSI,
		time.Now(),
		telemetry.AccelData.X,
		telemetry.AccelData.Y,
		telemetry.AccelData.Z,
		telemetry.GPSData.Latitude,
		telemetry.GPSData.Longitute,
		telemetry.GPSData.Speed,
		telemetry.GPSData.SatC,
		telemetry.PSUData.VoltageOut,
		telemetry.PSUData.CurrentOut,
		telemetry.PSUData.PowerOut,
		telemetry.PSUData.VoltageIn,
		telemetry.PSUData.WattHours,
		telemetry.SYSData.VoltageBat,
		telemetry.SYSData.BatConnected > 0,
		telemetry.SYSData.ErrorCode,
		telemetry.Meginajums,
	)
	if err != nil {
		logrus.WithError(err).Error("Failed to save Car Telemetry")
		return err
	}
	return nil
}
