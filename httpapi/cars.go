package httpapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
	"github.com/ksvaza/ees-link/data"
	"github.com/ksvaza/ees-link/db"
	"github.com/ksvaza/ees-link/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// func PointGetCars(r *http.Request, ps httprouter.Params) (*httpResult, error) {

// }

// //
// 	//
// 	// add account to team
// 	//
// 	// id/pendingteamid maģija
// 	//

// // jāatjauno konta info
//
//	if newAccount.Cilveks.ID != newAccount.PendingTeamID {
//		return errors.New("account registration logic error: pending team ID does not match team member ID")
//	}
func addOrRemoveAccountToTeam(ctx context.Context, account *models.Account) error {
	var realDB data.Database
	realDB = &db.RealDB{}

	if account.Verified || (account.PendingTeamID == account.Cilveks.ID) {
		return nil
	}

	if account.PendingTeamID != "" {
		err := realDB.AssignAccountsToTeamDataByUsername(ctx, account.PendingTeamID, []string{account.Username})
		if err != nil {
			return errors.Wrap(err, "Failed to add Account to team")
		}

		// sanāca

		account.Cilveks.ID = account.PendingTeamID
		account.Verified = true

		err = realDB.UpdateAccount(ctx, *account)
		if err != nil {
			return errors.Wrap(err, "Failed to update account")
		}
	} else {
		// izdzēst no pierakstītās komandas
		// nav implementēts
	}

	return nil
}

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

	prekey := ps.ByName("key") // key ļoti vajadzētu sakrist ar teamID
	if prekey == "" {
		return nil, errors.New("missing team key")
	}

	key, err := url.PathUnescape(prekey)
	if err != nil {
		return nil, errors.Wrap(err, "Unescape team key")
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

	// izdevās; jāatjauno pieteikuma statuss
	app, err := realDB.GetApplicationByID(r.Context(), key)
	if err != nil {
		return nil, err
	}
	app.Status = "accepted"
	err = realDB.UpdateApplication(r.Context(), *app)
	if err != nil {
		return nil, err
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

	prekey := ps.ByName("key")
	if prekey == "" {
		return nil, errors.New("missing team key")
	}

	key, err := url.PathUnescape(prekey)
	if err != nil {
		return nil, errors.Wrap(err, "Unescape team key")
	}

	// key, err := GenerateSalt(16)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "Generate salt")
	// }

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

	// Assign accounts to team if provided
	if len(newTeamData.Accounts) > 0 {
		var teamDataUsernames []string
		for _, account := range newTeamData.Accounts {
			teamDataUsernames = append(teamDataUsernames, account.Username)
		}

		err = realDB.AssignAccountsToTeamDataByUsername(r.Context(), key, teamDataUsernames)
		if err != nil {
			return nil, errors.Wrap(err, "Assign accounts to team data by username")
		}
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         "Team data updated",
	}, nil
}

func PointGetTeamDataByKey(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	logrus.Infof("PointGetTeamDataByKey called %+v", ps)
	if r.Method != http.MethodGet {
		return nil, errors.New("method not allowed")
	}

	prekey := ps.ByName("key")
	if prekey == "" {
		return nil, errors.New("missing team key")
	}

	key, err := url.PathUnescape(prekey)
	if err != nil {
		return nil, errors.Wrap(err, "Unescape team key")
	}

	adminaccount := GetAdminAccount(r.Context())
	if adminaccount == nil || !adminaccount.Superadmin || adminaccount.Username == "" || adminaccount.Password == "" || adminaccount.Salt == "" {
		// tad paskatīties; iespējams, komandas līderis ir autentificējies
		account := GetAccount(r.Context())
		if account == nil || account.Username == "" || account.Password == "" || account.Salt == "" {
			return &httpResult{
				ResponseType: http.StatusForbidden,
				Body:         `{"error":"forbidden"}`,
			}, nil
		} else {
			logrus.Infof("Account found in context: %s", account.Username)
		}
		if account.Cilveks.Role != "team_leader" || account.Cilveks.ID != key {
			return &httpResult{
				ResponseType: http.StatusForbidden,
				Body:         `{"error":"forbidden"}`,
			}, nil
		}
	} else {
		logrus.Infof("Admin account found in context: %s", adminaccount.Username)
	}

	// key, err := GenerateSalt(16)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "Generate salt")
	// }

	var realDB data.Database
	realDB = &db.RealDB{}

	teamData, err := realDB.GetTeamDataByKey(r.Context(), key)
	if err != nil {
		return nil, err
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         teamData,
	}, nil
}
