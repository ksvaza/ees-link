package db

import (
	"context"
	"time"

	"github.com/ksvaza/ees-link/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (DB *RealDB) SaveCarTelemetry(ctx context.Context, telemetry models.CarTelemetry) error {
	var (
		attempt  int    = 0
		raceName string = ""
	)

	if intermediate, exists := currentRaceInfo[telemetry.ID]; exists {
		attempt = intermediate.AttemptNr
		raceName = intermediate.RaceName
	}

	// Update livedata - super fake and inefficient
	if err := DB.UpdateLiveData(ctx, telemetry); err != nil {
		return errors.Wrap(err, "failed to update live telemetry data")
	}

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
		attempt,
		raceName,
	)
	if err != nil {
		logrus.WithError(err).Error("Failed to save Car Telemetry")
		return err
	}
	return nil
}
