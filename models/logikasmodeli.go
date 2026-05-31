package models

import "time"

// ----------------------------------------------------------------

// Punkti

// {"RaceName":"EkoRace","Points":[{"ID":1,"Points":12},{"ID":2,"Points":11},{"ID":3,"Points":10},{"ID":4,"Points":9},{"ID":5,"Points":8},{"ID":6,"Points":7},{"ID":7,"Points":6},{"ID":8,"Points":5},{"ID":9,"Points":4},{"ID":10,"Points":3},{"ID":11,"Points":2},{"ID":12,"Points":1}]}
type PointEntry struct {
	CarID  int `json:"ID"`
	Points int `json:"Points"`
}

type Points struct {
	RaceName string       `json:"RaceName"`
	Points   []PointEntry `json:"Points"`
}

// ----------------------------------------------------------------

type RacePointEntry struct {
	RaceName string `json:"race_name"`
	//(?)Attempt  string `json:"attempt_nr"`
	Points int `json:"points"`
}

type LeaderboardEntry struct {
	AgeGroup string `json:"age_group"` // nerādās līderu tabulā, bet tāpat vajadzīgs tabulu filtrēšanai

	CarID    int    `json:"car_id"`
	Avatar   []byte `json:"avatar"` // encoding/json marshals []byte as base64
	TeamName string `json:"team_name"`

	RaceEntries []RacePointEntry `json:"race_entries"`
	TotalPoints int              `json:"total_points"`
}

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

// Mašīnu konfigurācijas datu struktūras

// Sacensību konfigurācijas datu struktūras

type AdminSettings struct {
	PowerCoefficient float64 `json:"PowerCoef"`
	MaximumSpeed     float64 `json:"MaxSpd"`
}
