package db

import (
	"context"

	"github.com/ksvaza/ees-link/models"
	"github.com/sirupsen/logrus"
)

/*
	type ResultsEntry struct {
		// filtri
		RaceName  string `json:"race_name"`
		AttemptNr int    `json:"attempt_nr"` // Mēģinājuma Nr. 				// EkoRace, DragRace, ShuttleRun, SteeringTest, BrakingTest

		// mašīnas dati
		CarID    int    `json:"car_id"`    // ID
		TeamName string `json:"team_name"` // Komanda

		// Manual entry
		Penalties       time.Duration `json:"penalties"`        // Sodi (s) 					// EkoRace, MainRace, ShuttleRun, SteeringTest
		UsedTime        time.Duration `json:"used_time"`        // Patērētais laiks 			// EkoRace, MainRace, ShuttleRun, SteeringTest
		MaxSpeed        float64       `json:"max_speed"`        // Maks. ātrums (km/h) 			// DragRace
		Valid           bool          `json:"valid"`            // Derīgs 						// EkoRace, DragRace, ShuttleRun, SteeringTest, BrakingTest
		BrakingDistance float64       `json:"braking_distance"` // Bremzēšanas distance (m) 	// BrakingTest
		DriveinSpeed    float64       `json:"drivein_speed"`    // Iebraukšanas ātrums (km/h)	// BrakingTest

		// Calculated fields
		AverageSpeed    float64       `json:"average_speed"`    // Vid. ātrums = distance / TotalTime 					// EkoRace, MainRace
		TotalTime       time.Duration `json:"total_time"`       // Kopējais laiks = patērētais laiks + sodi 			// EkoRace, MainRace, DragRace, ShuttleRun, SteeringTest
		Placement       int           `json:"placement"`        // Vieta = mainrace:totaltime, ekorace: skatās secību 	// EkoRace, MainRace
		UsedEnergy      float64       `json:"used_energy"`      // Patērētā enerģija (Wh) 								// EkoRace
		Efficiency      float64       `json:"efficiency"`       // Efektivitāte (Wh/kg) 								// EkoRace
		ShellEfficiency float64       `json:"shell_efficiency"` // "Shell" efektivitāte (km/kWh) 						// EkoRace
	}
*/
func (DB *RealDB) SaveResultsEntry(ctx context.Context, resultsEntry models.ResultsEntry) error {
	_, err := Pool.Exec(ctx, ResultsEntryWriteRequest,
		resultsEntry.RaceName,
		resultsEntry.AttemptNr,
		resultsEntry.CarID,
		resultsEntry.TeamName,
		resultsEntry.StartTime,
		resultsEntry.FinishTime,
		resultsEntry.Penalties,
		resultsEntry.UsedTime,
		resultsEntry.MaxSpeed,
		resultsEntry.Valid,
		resultsEntry.BrakingDistance,
		resultsEntry.DriveinSpeed,
		resultsEntry.AverageSpeed,
		resultsEntry.TotalTime,
		resultsEntry.Placement,
		resultsEntry.UsedEnergy,
		resultsEntry.Efficiency,
		resultsEntry.ShellEfficiency,
	)

	if err != nil {
		logrus.WithError(err).Error("Failed to write results entry")
		return err
	}

	return nil
}

// UpdateResultsEntry updates an existing results entry in the database
func (DB *RealDB) UpdateResultsEntry(ctx context.Context, resultsEntry models.ResultsEntry) error {
	_, err := Pool.Exec(ctx, ResultsEntryUpdateRequest,
		resultsEntry.RaceName,
		resultsEntry.AttemptNr,
		resultsEntry.CarID,
		resultsEntry.TeamName,
		resultsEntry.StartTime,
		resultsEntry.FinishTime,
		resultsEntry.Penalties,
		resultsEntry.UsedTime,
		resultsEntry.MaxSpeed,
		resultsEntry.Valid,
		resultsEntry.BrakingDistance,
		resultsEntry.DriveinSpeed,
		resultsEntry.AverageSpeed,
		resultsEntry.TotalTime,
		resultsEntry.Placement,
		resultsEntry.UsedEnergy,
		resultsEntry.Efficiency,
		resultsEntry.ShellEfficiency,
	)

	if err != nil {
		logrus.WithError(err).
			WithField("race_name", resultsEntry.RaceName).
			WithField("attempt_nr", resultsEntry.AttemptNr).
			WithField("car_id", resultsEntry.CarID).
			Error("Failed to update results entry")
		return err
	}

	logrus.WithField("race_name", resultsEntry.RaceName).
		WithField("attempt_nr", resultsEntry.AttemptNr).
		WithField("car_id", resultsEntry.CarID).
		Info("Results entry updated successfully")

	return nil
}

// SaveOrUpdateResultsEntry inserts a new results entry or updates if it already exists
func (DB *RealDB) SaveOrUpdateResultsEntry(ctx context.Context, resultsEntry models.ResultsEntry) error {
	_, err := Pool.Exec(ctx, ResultsEntryUpsertRequest,
		resultsEntry.RaceName,
		resultsEntry.AttemptNr,
		resultsEntry.CarID,
		resultsEntry.TeamName,
		resultsEntry.StartTime,
		resultsEntry.FinishTime,
		resultsEntry.Penalties,
		resultsEntry.UsedTime,
		resultsEntry.MaxSpeed,
		resultsEntry.Valid,
		resultsEntry.BrakingDistance,
		resultsEntry.DriveinSpeed,
		resultsEntry.AverageSpeed,
		resultsEntry.TotalTime,
		resultsEntry.Placement,
		resultsEntry.UsedEnergy,
		resultsEntry.Efficiency,
		resultsEntry.ShellEfficiency,
	)

	if err != nil {
		logrus.WithError(err).
			WithField("race_name", resultsEntry.RaceName).
			WithField("attempt_nr", resultsEntry.AttemptNr).
			WithField("car_id", resultsEntry.CarID).
			Error("Failed to save or update results entry")
		return err
	}

	logrus.WithField("race_name", resultsEntry.RaceName).
		WithField("attempt_nr", resultsEntry.AttemptNr).
		WithField("car_id", resultsEntry.CarID).
		Info("Results entry saved or updated successfully")

	return nil
}

// fake
func (DB *RealDB) GetResultsEntriesByRaceName(ctx context.Context, raceName string) ([]models.ResultsEntry, error) {
	// var entry models.ResultsEntry
	// err := Pool.QueryRow(ctx, ResultsEntriesReadByRaceNameRequest, raceName).Scan(
	// 	&entry.RaceName,
	// 	&entry.AttemptNr,
	// 	&entry.CarID,
	// 	&entry.TeamName,
	// 	&entry.StartTime,
	// 	&entry.FinishTime,
	// 	&entry.Penalties,
	// 	&entry.UsedTime,
	// 	&entry.MaxSpeed,
	// 	&entry.Valid,
	// 	&entry.BrakingDistance,
	// 	&entry.DriveinSpeed,
	// 	&entry.AverageSpeed,
	// 	&entry.TotalTime,
	// 	&entry.Placement,
	// 	&entry.UsedEnergy,
	// 	&entry.Efficiency,
	// 	&entry.ShellEfficiency,
	// )
	// if err != nil {
	// 	if err == pgx.ErrNoRows {
	// 		return nil, nil
	// 	}
	// 	logrus.WithError(err).WithField("race_name", raceName).Error("Failed to query result_entries table")
	// 	return nil, err
	// }
	// defer rows.Close()

	// var res []models.ResultsEntry
	// for rows.Next() {
	// 	var entry models.ResultsEntry
	// 	err = rows.Scan(ResultsEntryUpdateRequest,
	// 		&entry.RaceName,
	// 		&entry.AttemptNr,
	// 		&entry.CarID,
	// 		&entry.TeamName,
	// 		&entry.StartTime,
	// 		&entry.FinishTime,
	// 		&entry.Penalties,
	// 		&entry.UsedTime,
	// 		&entry.MaxSpeed,
	// 		&entry.Valid,
	// 		&entry.BrakingDistance,
	// 		&entry.DriveinSpeed,
	// 		&entry.AverageSpeed,
	// 		&entry.TotalTime,
	// 		&entry.Placement,
	// 		&entry.UsedEnergy,
	// 		&entry.Efficiency,
	// 		&entry.ShellEfficiency,
	// 	)
	// 	if err != nil {
	// 		logrus.WithError(err).Error("Failed to scan result_entries row")
	// 		return nil, err
	// 	}
	// 	res = append(res, entry)
	// }

	// if err = rows.Err(); err != nil {
	// 	logrus.WithError(err).Error("Failed iterating result_entries rows")
	// 	return nil, err
	// }

	// if len(res) == 0 {
	return nil, nil
	// }

	// return res, nil
}

// fake
func (DB *RealDB) GetResultEntryByCarID(ctx context.Context, carID int) (*models.ResultsEntry, error) {
	// var points models.Points
	// var pointsData []byte
	// err := Pool.QueryRow(ctx, ResultsEntryReadByCarID, raceName).Scan(&points.RaceName, &pointsData)
	// if err != nil {
	// 	logrus.WithError(err).Error("Failed to query points by race name")
	// 	return nil, err
	// }
	// err = json.Unmarshal(pointsData, &points.Points)
	// if err != nil {
	// 	logrus.WithError(err).Error("Failed to unmarshal points data")
	// 	return nil, err
	// }
	// return &points, nil

	return nil, nil
}
