package httpapi

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ksvaza/ees-link/data"
	"github.com/ksvaza/ees-link/db"
	"github.com/ksvaza/ees-link/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func PointGetCars(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	if r.Method != http.MethodGet {
		return nil, errors.New("method not allowed")
	}

	var realDB data.Database
	realDB = &db.RealDB{}

	adminAccount := GetAdminAccount(r.Context())
	if adminAccount == nil || !adminAccount.Superadmin {
		return &httpResult{
			ResponseType: http.StatusForbidden,
			Body:         `{"error":"forbidden"}`,
		}, nil
	}
	logrus.Infof("Superadmin account found in context: %s", adminAccount.Username)

	params, err := realDB.GetAllCarParameters(r.Context())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get car parameters")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         params,
	}, nil
}

func PointPostCars(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed")
	}

	var realDB data.Database
	realDB = &db.RealDB{}

	adminAccount := GetAdminAccount(r.Context())
	if adminAccount == nil || !adminAccount.Superadmin {
		return &httpResult{
			ResponseType: http.StatusForbidden,
			Body:         `{"error":"forbidden"}`,
		}, nil
	}
	logrus.Infof("Superadmin account found in context: %s", adminAccount.Username)

	var newCarParams []models.CarParameters
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read body")
	}

	if err := json.Unmarshal(body, &newCarParams); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	err = realDB.ReplaceAllCarParameters(r.Context(), newCarParams)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update car parameters")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         "Car parameters replaced successfully",
	}, nil
}

func PointRaceStart(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed")
	}

	var realDB data.Database
	realDB = &db.RealDB{}
	adminAccount := GetAdminAccount(r.Context())
	if adminAccount == nil || !adminAccount.Superadmin {
		return &httpResult{
			ResponseType: http.StatusForbidden,
			Body:         `{"error":"forbidden"}`,
		}, nil
	}
	logrus.Infof("Superadmin account found in context: %s", adminAccount.Username)

	var raceStarts []models.RaceStart
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read body")
	}

	if err := json.Unmarshal(body, &raceStarts); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	err = realDB.StartRace(r.Context(), raceStarts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start race")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         "Race started successfully",
	}, nil
}

func PointCarFinish(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed")
	}

	var realDB data.Database
	realDB = &db.RealDB{}

	var finishData models.CarFinish
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read body")
	}

	if err := json.Unmarshal(body, &finishData); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	err = realDB.FinishCar(r.Context(), finishData.CarID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save car finish data")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         "Car finish data saved successfully",
	}, nil
}
