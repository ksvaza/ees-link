package db

import (
	"context"

	"github.com/ksvaza/ees-link/models"
)

func (DB *RealDB) BuildLeaderboard(ctx context.Context, ageGroup string) ([]models.LeaderboardEntry, error) {
	if ageGroup == "" {
		return nil, nil
	}

	var leaderboard []models.LeaderboardEntry

	if ageGroup != "all" {
		// get all car parameters for the age group
		carParams, err := DB.GetCarParametersByAgeGroup(ctx, ageGroup)
		if err != nil {
			return nil, err
		}
		leaderboard = make([]models.LeaderboardEntry, len(carParams))

		for i, carParam := range carParams {
			leaderboard[i] = models.LeaderboardEntry{
				AgeGroup:    ageGroup,
				CarID:       carParam.CarID,
				Avatar:      carParam.Avatar,
				TeamName:    carParam.TeamName,
				RaceEntries: make([]models.RacePointEntry, 0),
				TotalPoints: 0,
			}
		}

		// get all points from the points table
		points, err := DB.GetAllPoints(ctx)
		if err != nil {
			return nil, err
		}

		// fill out the leaderboard entries with the points for each car id
		for _, p := range points { // each is a different race
			for _, pointEntry := range p.Points {
				for _, le := range leaderboard {
					if pointEntry.CarID == le.CarID {
						rpe := models.RacePointEntry{
							RaceName: p.RaceName,
							Points:   pointEntry.Points,
						}
						le.RaceEntries = append(le.RaceEntries, rpe)
						break
					}
				}
			}
		}

		// calculate total points
		for _, le := range leaderboard {
			for _, re := range le.RaceEntries {
				le.TotalPoints += re.Points
			}
		}

		return leaderboard, nil
	} else {
		// get all car parameters for the age group
		carParams, err := DB.GetAllCarParameters(ctx)
		if err != nil {
			return nil, err
		}
		leaderboard = make([]models.LeaderboardEntry, len(carParams))

		for i, carParam := range carParams {
			leaderboard[i] = models.LeaderboardEntry{
				AgeGroup:    ageGroup,
				CarID:       carParam.CarID,
				Avatar:      carParam.Avatar,
				TeamName:    carParam.TeamName,
				RaceEntries: make([]models.RacePointEntry, 0),
				TotalPoints: 0,
			}
		}

		// get all points from the points table
		points, err := DB.GetAllPoints(ctx)
		if err != nil {
			return nil, err
		}

		// fill out the leaderboard entries with the points for each car id
		for _, p := range points { // each is a different race
			for _, pointEntry := range p.Points {
				for _, le := range leaderboard {
					if pointEntry.CarID == le.CarID {
						rpe := models.RacePointEntry{
							RaceName: p.RaceName,
							Points:   pointEntry.Points,
						}
						le.RaceEntries = append(le.RaceEntries, rpe)
						break
					}
				}
			}
		}

		// calculate total points
		for _, le := range leaderboard {
			for _, re := range le.RaceEntries {
				le.TotalPoints += re.Points
			}
		}

		return leaderboard, nil
	}
}
