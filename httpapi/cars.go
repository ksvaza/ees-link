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

// func PointGetCars(r *http.Request, ps httprouter.Params) (*httpResult, error) {

// }

func PointPostTeamDataByKey(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	logrus.Infof("PointPostTeamDataByKey called %+v", ps)
	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed")
	}

	adminaccount := GetAdminAccount(r.Context())
	if adminaccount == nil || !adminaccount.Superadmin || adminaccount.Username == "" || adminaccount.Password == "" || adminaccount.Salt == "" {
		return &httpResult{
			ResponseType: http.StatusForbidden,
			Body:         `{"error":"forbidden"}`,
		}, nil
	}

	key := ps.ByName("key")
	if key == "" {
		return nil, errors.New("missing team key")
	}

	var realDB data.Database
	realDB = &db.RealDB{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read body")
	}

	var newTeamData models.TeamData
	if err := json.Unmarshal(body, &newTeamData); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	existingTeamData, err := realDB.GetTeamDataByKey(r.Context(), key)
	if err != nil {
		return nil, errors.Wrap(err, "Get existing team data by key")
	}
	if existingTeamData != nil {
		return nil, errors.New("team data already exists")
	}

	err = realDB.RegisterTeamData(r.Context(), newTeamData)
	if err != nil {
		return nil, errors.Wrap(err, "Register team data")
	}

	var teamDataUsernames []string

	for _, account := range newTeamData.Accounts {
		teamDataUsernames = append(teamDataUsernames, account.Username)
	}

	err = realDB.AssignAccountsToTeamDataByUsername(r.Context(), key, teamDataUsernames)
	if err != nil {
		return nil, errors.Wrap(err, "Assign accounts to team data by username")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         "Team data registered",
	}, nil
}

func PointPatchTeamDataByKey(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	logrus.Infof("PointPatchTeamDataByKey called %+v", ps)
	if r.Method != http.MethodPatch {
		return nil, errors.New("method not allowed")
	}

	adminaccount := GetAdminAccount(r.Context())
	if adminaccount == nil || !adminaccount.Superadmin || adminaccount.Username == "" || adminaccount.Password == "" || adminaccount.Salt == "" {
		return &httpResult{
			ResponseType: http.StatusForbidden,
			Body:         `{"error":"forbidden"}`,
		}, nil
	}

	logrus.Infof("Admin account found in context: %s", adminaccount.Username)

	// key := ps.ByName("key")
	// if key == "" {
	// 	return nil, errors.New("missing team key")
	// }

	key, err := GenerateSalt(16)
	if err != nil {
		return nil, errors.Wrap(err, "Generate salt")
	}

	var realDB data.Database
	realDB = &db.RealDB{}

	var newTeamData models.TeamData

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read body")
	}

	if err := json.Unmarshal(body, &newTeamData); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	err = realDB.UpdateTeamDataByKey(r.Context(), key, newTeamData)
	if err != nil {
		return nil, errors.Wrap(err, "Update team data by key")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         "Team data updated",
	}, nil
}
