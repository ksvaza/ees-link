package db

import (
	"context"
	"errors"
	"time"

	"github.com/ksvaza/ees-link/models"
)

/*
type RaceInstance struct {
	RaceName  string `json:"raceName"`
	AttemptNr int    `json:"attempt"`
}

type RaceStartFinishTableEntry struct {
	StartTime  time.Time `json:"startTime"`
	FinishTime time.Time `json:"finishTime"`
}

type RaceStartFinishTable map[int] */ /*CarID*/ /* map[RaceInstance]RaceStartFinishTableEntry
 */
var (
	raceStartFinishTable models.RaceStartFinishTable = make(models.RaceStartFinishTable)
	currentRaceInfo      map[int]models.RaceInstance = make(map[int]models.RaceInstance)
)

func (DB *RealDB) StartRace(ctx context.Context, raceStarts []models.RaceStart) error {
	startTime := time.Now()
	for _, raceStart := range raceStarts {
		carID := raceStart.CarID

		if _, exists := raceStartFinishTable[carID]; !exists {
			raceStartFinishTable[carID] = make(map[models.RaceInstance]models.RaceStartFinishTableEntry)
		}

		raceInstance := models.RaceInstance{
			RaceName:  raceStart.RaceName,
			AttemptNr: raceStart.AttemptNr,
		}
		raceStartFinishTable[carID][raceInstance] = models.RaceStartFinishTableEntry{
			StartTime:  startTime,
			FinishTime: time.Time{}, // zero value, means not finished yet
		}
		currentRaceInfo[carID] = raceInstance
	}
	return nil
}

func (DB *RealDB) FinishCar(ctx context.Context, carID int) error {
	finishTime := time.Now()
	var vidus map[models.RaceInstance]models.RaceStartFinishTableEntry
	vidus, existsV := raceStartFinishTable[carID]
	if !existsV {
		return errors.New("car not found")
	}
	// for raceInstance, entry := range raceStartFinishTable[carID] {
	// 	if entry.FinishTime.IsZero() { // means this race instance is not finished
	// 		entry.FinishTime = finishTime
	// 		raceStartFinishTable[carID][raceInstance] = entry
	// 		break // assuming a car can only have one active race instance at a time
	// 	}
	// }

	if ri, ok := currentRaceInfo[carID]; ok {
		if entry, exists := vidus[ri]; exists {
			if entry.FinishTime.IsZero() {
				entry.FinishTime = finishTime
				raceStartFinishTable[carID][ri] = entry
				currentRaceInfo[carID] = models.RaceInstance{
					RaceName:  "",
					AttemptNr: 0,
				}
			} else {
				return errors.New("car has already finished")
			}
		} else {
			return errors.New("car has race start entry")
		}
	} else {
		return errors.New("car has no current race")
	}

	// triggers the calculation of race results
	return nil
}
