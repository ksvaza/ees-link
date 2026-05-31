package models

import "time"

type PointEntry struct {
	CarID   int    `json:"car_id"`
	Points  int    `json:"points"`
	Attempt string `json:"attempt_nr"`
}

type RacePointEntry struct {
	RaceName string `json:"race_name"`
	Attempt  string `json:"attempt_nr"`
	Points   int    `json:"points"`
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
	RaceName string `json:"race_name"` // filtrs

	CarID    int    `json:"car_id"`    // ID
	TeamName string `json:"team_name"` // Komanda

	UsedTime        time.Duration `json:"used_time"`        // Patērētais laiks 			// EkoRace, MainRace, ShuttleRun, SteeringTest
	AverageSpeed    float64       `json:"average_speed"`    // Vid. ātrums 					// EkoRace, MainRace
	Penalties       time.Duration `json:"penalties"`        // Sodi (s) 					// EkoRace, MainRace, ShuttleRun, SteeringTest
	TotalTime       time.Duration `json:"total_time"`       // Kopējais laiks 				// EkoRace, MainRace, DragRace, ShuttleRun, SteeringTest
	Placement       int           `json:"placement"`        // Vieta 						// EkoRace, MainRace
	AttemptNr       int           `json:"attempt_nr"`       // Mēģinājuma Nr. 				// EkoRace, DragRace, ShuttleRun, SteeringTest, BrakingTest
	Useful          string        `json:"useful"`           // Derīgs 						// EkoRace, DragRace, ShuttleRun, SteeringTest, BrakingTest
	UsedEnergy      float64       `json:"used_energy"`      // Patērētā enerģija (kWh) 		// EkoRace
	Efficiency      float64       `json:"efficiency"`       // Efektivitāte (Wh/kg) 		// EkoRace
	MaxSpeed        float64       `json:"max_speed"`        // Maks. ātrums (km/h) 			// DragRace
	BrakingDistance float64       `json:"braking_distance"` // Bremzēšanas distance (m) 	// BrakingTest
	DriveinSpeed    float64       `json:"drivein_speed"`    // Iebraukšanas ātrums (km/h)	// BrakingTest
}

// Mašīnu konfigurācijas datu struktūras

// Sacensību konfigurācijas datu struktūras

type AdminSettings struct {
	PowerCoefficient float64 `json:"PowerCoef"`
	MaximumSpeed     float64 `json:"MaxSpd"`
}
