package db

import (
	"context"

	"github.com/ksvaza/ees-link/models"
	"github.com/sirupsen/logrus"
)

func (DB *RealDB) GetAllCarParameters(ctx context.Context) ([]models.CarParameters, error) {
	rows, err := Pool.Query(ctx, GetAllCarParametersRequest)
	if err != nil {
		logrus.WithError(err).Error("Failed to query car parameters table")
		return nil, err
	}
	defer rows.Close()

	var params []models.CarParameters
	for rows.Next() {
		var p models.CarParameters
		err = rows.Scan(
			&p.CarID,
			&p.TeamName,
			&p.SetVoltage,
			&p.CalculatedCurrent,
			&p.Mass,
			&p.AgeGroup,
			&p.Avatar,
			&p.FinishedAt,
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to scan car parameters row")
			return nil, err
		}
		params = append(params, p)
	}

	if err = rows.Err(); err != nil {
		logrus.WithError(err).Error("Failed iterating car parameters rows")
		return nil, err
	}

	if len(params) == 0 {
		return nil, nil
	}

	return params, nil
}

func (DB *RealDB) GetCarParametersByAgeGroup(ctx context.Context, ageGroup string) ([]models.CarParameters, error) {
	rows, err := Pool.Query(ctx, GetCarParametersByAgeGroupRequest, ageGroup)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to query car parameters by age group %s", ageGroup)
		return nil, err
	}
	defer rows.Close()

	var params []models.CarParameters
	for rows.Next() {
		var p models.CarParameters
		err = rows.Scan(
			&p.CarID,
			&p.TeamName,
			&p.SetVoltage,
			&p.CalculatedCurrent,
			&p.Mass,
			&p.AgeGroup,
			&p.Avatar,
			&p.FinishedAt,
		)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to scan car parameters row for age group %s", ageGroup)
			return nil, err
		}
		params = append(params, p)
	}

	if err = rows.Err(); err != nil {
		logrus.WithError(err).Errorf("Failed iterating car parameters rows for age group %s", ageGroup)
		return nil, err
	}

	if len(params) == 0 {
		return nil, nil
	}

	return params, nil
}

func (DB *RealDB) ReplaceAllCarParameters(ctx context.Context, carParams []models.CarParameters) error {
	tx, err := Pool.Begin(ctx)
	if err != nil {
		logrus.WithError(err).Error("Failed to begin transaction for replacing car parameters")
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, DeleteAllCarParametersRequest)
	if err != nil {
		logrus.WithError(err).Error("Failed to delete existing car parameters")
		return err
	}

	for _, p := range carParams {
		_, err = tx.Exec(ctx, InsertCarParametersRequest,
			p.CarID,
			p.TeamName,
			p.SetVoltage,
			p.CalculatedCurrent,
			p.Mass,
			p.AgeGroup,
			p.Avatar,
			p.FinishedAt,
		)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to insert car parameters for team %s", p.TeamName)
			return err
		}
	}

	return nil
}
