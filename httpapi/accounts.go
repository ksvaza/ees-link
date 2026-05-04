package httpapi

import (
	"context"
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

func RegisterAdminAccount(admin models.AdminAccount) error {
	salt, err := GenerateSalt(16)
	if err != nil {
		return errors.Wrap(err, "failed to generate salt")
	}

	var realDB data.Database
	realDB = &db.RealDB{}

	admins, err := realDB.GetAllAdmins(context.Background())

	if err != nil {
		return errors.Wrap(err, "failed to get all admins")
	}

	for _, a := range admins {
		if a.Username == admin.Username {
			logrus.Infof("Admin account with username '%s' already exists", admin.Username)
			return nil
		}
	}

	admin.Salt = salt
	admin.Password = hashPassword(admin.Password, salt)
	admin.Key, err = GenerateSalt(16)
	if err != nil {
		return errors.Wrap(err, "failed to generate salt for admin key")
	}

	err = realDB.RegisterNewAdmin(context.Background(), admin)
	if err != nil {
		return errors.Wrap(err, "failed to register new admin")
	}

	return nil
}

func PointLogin(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	logrus.Infof("PointLogin called %+v", ps)
	if r.Method != http.MethodGet {
		return nil, errors.New("method not allowed")
	}
	adminaccount := GetAdminAccount(r.Context())
	if adminaccount != nil && adminaccount.Username != "" && adminaccount.Password != "" && adminaccount.Salt != "" {
		logrus.Infof("Authenticated admin account: %+v", adminaccount)
		adminaccountInfo := models.AdminAccount{
			Superadmin: adminaccount.Superadmin,
			Username:   adminaccount.Username}
		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         adminaccountInfo,
		}, nil
	}
	account := GetAccount(r.Context())
	if account != nil && account.Username != "" && account.Password != "" && account.Salt != "" {
		logrus.Infof("Authenticated account: %+v", account)
		accountInfo := models.Account{
			Cilveks:     account.Cilveks,
			Username:    account.Username,
			Email:       account.Email,
			PhoneNumber: account.PhoneNumber,
		}
		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         accountInfo,
		}, nil
	}

	return &httpResult{
		ResponseType: http.StatusUnauthorized,
		Body:         `{"error":"unauthorized"}`,
	}, nil
}

// ----------------------------------------------------------------

// Account application management
// ----------------------------------------------------------------

func registerAccountByApplication(ctx context.Context, realDB data.Database, app models.AccountApplication) error {
	// Šeit var sarakstīt, kas ir un nav obliātie lauki.
	if app.Key == "" || app.FullName == "" || app.DateOfBirth == "" || app.Role == "" || app.Password == "" || app.Username == "" || app.Email == "" || app.PhoneNumber == "" {
		return errors.New("missing required fields")
	}

	existing, err := realDB.GetAccountByUsername(ctx, app.Username)
	if err != nil {
		return errors.Wrap(err, "failed to check username")
	}
	if existing != nil {
		return errors.New("account by username already exists")
	}
	var teamID string = ""

	team, err := realDB.GetApplicationByTeamName(ctx, app.TeamName)
	if err != nil {
		return errors.Wrap(err, "failed to get application by team name")
	}
	if team == nil {
		if app.TeamName == "" {
			logrus.Infof("Registering account without team association")
			teamID = ""
		} else {
			return errors.New("team not found")
		}
	} else {
		teamID = team.ID
		logrus.Infof("Registering account with team association: %s (team ID: %s)", app.TeamName, teamID)
	}

	salt, err := GenerateSalt(16)
	if err != nil {
		return errors.Wrap(err, "failed to generate salt")
	}

	account := models.Account{
		Cilveks: models.TeamMember{
			Key:         app.Key,
			FullName:    app.FullName,
			DateOfBirth: app.DateOfBirth,
			Role:        app.Role,
			ID:          teamID,
		},
		Password:    hashPassword(app.Password, salt),
		Username:    app.Username,
		Email:       app.Email,
		PhoneNumber: app.PhoneNumber,
		Salt:        salt,
	}

	err = realDB.RegisterNewAccount(ctx, account)
	if err != nil {
		return errors.Wrap(err, "failed to register new account")
	}

	return nil
}

func registerTeamLeaderAccountByApplication(ctx context.Context, realDB data.Database, app models.AccountApplication, teamIDcheck string) error {
	if app.TeamName == "" {
		return errors.New("missing team name for team leader application")
	}

	if app.Key == "" || app.FullName == "" || app.DateOfBirth == "" || app.Password == "" || app.Username == "" || app.Email == "" || app.PhoneNumber == "" {
		return errors.New("missing required fields")
	}

	existing, err := realDB.GetAccountByUsername(ctx, app.Username)
	if err != nil {
		return errors.Wrap(err, "failed to check username")
	}
	if existing != nil {
		return errors.New("account by username already exists")
	}

	// Get the team application by name
	teamApp, err := realDB.GetApplicationByTeamName(ctx, app.TeamName)
	if err != nil {
		return errors.Wrap(err, "failed to get team application")
	}
	if teamApp == nil {
		return errors.New("team application not found")
	}

	// Validate that the team application ID matches the expected ID
	if teamApp.ID != teamIDcheck {
		return errors.New("team ID mismatch - verification failed")
	}

	// Find the team leader member in the application
	member := teamApp.FindTeamMemberByFullName(app.FullName)
	if member == nil {
		return errors.New("team member not found in application")
	}

	// Verify the member is a team leader
	if member.Role != "team_leader" {
		return errors.New("member is not a team leader")
	}

	// Verify date of birth matches
	if app.DateOfBirth != member.DateOfBirth {
		return errors.Wrap(errors.New("date of birth mismatch"), "team member verification failed")
	}

	logrus.Infof("Registering team leader account with team association: %s (team ID: %s)", app.TeamName, teamApp.ID)

	salt, err := GenerateSalt(16)
	if err != nil {
		return errors.Wrap(err, "failed to generate salt")
	}

	account := models.Account{
		Cilveks: models.TeamMember{
			Key:         app.Key,
			FullName:    app.FullName,
			DateOfBirth: app.DateOfBirth,
			Role:        "team_leader",
			ID:          teamApp.ID,
		},
		Password:    hashPassword(app.Password, salt),
		Username:    app.Username,
		Email:       app.Email,
		PhoneNumber: app.PhoneNumber,
		Salt:        salt,
	}

	err = realDB.RegisterNewAccount(ctx, account)
	if err != nil {
		return errors.Wrap(err, "failed to register new account")
	}

	return nil
}

func PointRegister(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	logrus.Infof("PointRegister called %+v", ps)
	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed")
	}

	var pieteikums models.AccountApplication
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read body")
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &pieteikums); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	logrus.Infof("Received form data: %+v", pieteikums)

	pieteikums.Key, err = GenerateSalt(16)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate salt for random key")
	}

	var realDB data.Database
	realDB = &db.RealDB{}

	if pieteikums.TeamName == "" {
		err = registerAccountByApplication(r.Context(), realDB, pieteikums)
		if err != nil {
			return nil, errors.Wrap(err, "failed to register account by application")
		}

		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         `{"status":"accepted"}`,
		}, nil
	}

	if pieteikums.Role == "team_leader" {
		// http://e-es.lv/api/register?uniqueID=1234567890
		id := r.URL.Query().Get("uniqueID")
		if id != "" {
			err = registerTeamLeaderAccountByApplication(r.Context(), realDB, pieteikums, id)
			if err == nil {
				return &httpResult{
					ResponseType: http.StatusOK,
					Body:         `{"status":"accepted"}`,
				}, nil
			}
			// If verification or registration fails, fall through to save as pending application
			logrus.WithError(err).Warnf("Failed to automatically verify and register team leader: %v", err)
		}
	}

	pieteikums.Key, err = GenerateSalt(16)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate salt for random key")
	}

	pieteikums.Password = hashPassword(pieteikums.Password, pieteikums.Key)

	err = realDB.RegisterNewAccountApplication(r.Context(), pieteikums)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register new account application")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         `{"status":"pending verification"}`,
	}, nil
}

func PointGetAccountApplications(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	if r.Method != http.MethodGet {
		return nil, errors.New("method not allowed")
	}

	var realDB data.Database
	realDB = &db.RealDB{}

	adminAccount := GetAdminAccount(r.Context())
	if adminAccount != nil && adminAccount.Superadmin {
		logrus.Infof("Superadmin account found in context: %s", adminAccount.Username)
		applications, err := realDB.GetAccountApplications(r.Context())
		if err != nil {
			return nil, errors.Wrap(err, "failed to get account applications")
		}

		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         applications,
		}, nil
	}

	account := GetAccount(r.Context())
	if account != nil && account.Cilveks.Role == "team_leader" {
		logrus.Infof("Team leader account found in context: %s", account.Username)

		applications, err := realDB.GetAccountApplications(r.Context())
		if err != nil {
			return nil, errors.Wrap(err, "failed to get account applications")
		}

		registrationApplication, err := realDB.GetApplicationByID(r.Context(), account.Cilveks.ID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get application by ID")
		}

		if registrationApplication == nil {
			return &httpResult{
				ResponseType: http.StatusNotFound,
				Body:         `{"error":"registration application not found"}`,
			}, nil
		}

		var teamApplications []models.AccountApplication

		for _, app := range applications {
			if app.TeamName == registrationApplication.TeamName {
				teamApplications = append(teamApplications, app)
			}
		}

		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         teamApplications,
		}, nil
	}

	return &httpResult{
		ResponseType: http.StatusForbidden,
		Body:         `{"error":"forbidden"}`,
	}, nil
}

func buildVerificationCriteria(ctx context.Context, realDB data.Database, app models.AccountApplication) (*models.AccountVerificationCriteria, error) {
	criteria := &models.AccountVerificationCriteria{}

	// General requirements - visible to everyone
	criteria.RequiredFieldsPresent = app.Key != "" && app.FullName != "" && app.DateOfBirth != "" && app.Password != "" && app.Username != "" && app.Email != "" && app.PhoneNumber != "" && app.Role != ""

	existing, err := realDB.GetAccountByUsername(ctx, app.Username)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check username availability")
	}
	criteria.UsernameAvailable = existing == nil

	// Check user context
	adminAccount := GetAdminAccount(ctx)
	account := GetAccount(ctx)

	var isSuperadmin bool
	var isTeamLeader bool
	var teamLeaderTeamID string

	if adminAccount != nil && adminAccount.Superadmin {
		isSuperadmin = true
	} else if account != nil && account.Cilveks.Role == "team_leader" {
		isTeamLeader = true
		teamLeaderTeamID = account.Cilveks.ID
	}

	// If superadmin or team leader checking their own team's application, show detailed criteria
	if isSuperadmin || (isTeamLeader && app.TeamName != "" && teamLeaderTeamID != "") {
		// Check if this is for the team leader's team (if they're a team leader)
		if isTeamLeader {
			teamApp, err := realDB.GetApplicationByID(ctx, teamLeaderTeamID)
			if err != nil || teamApp == nil || teamApp.TeamName != app.TeamName {
				// Not their team, only show general criteria
				return criteria, nil
			}
		}

		// Standalone account criteria
		noTeamName := app.TeamName == ""
		notTeamLeader := app.Role != "team_leader"
		criteria.NoTeamNameProvided = &noTeamName
		criteria.NotTeamLeaderRole = &notTeamLeader

		// Team-related criteria
		teamNameProvided := app.TeamName != ""
		criteria.TeamNameProvided = &teamNameProvided

		if teamNameProvided {
			teamApp, err := realDB.GetApplicationByTeamName(ctx, app.TeamName)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get team application")
			}

			teamExists := teamApp != nil
			criteria.TeamApplicationExists = &teamExists

			if teamExists {
				member := teamApp.FindTeamMemberByFullName(app.FullName)
				memberFound := member != nil
				criteria.TeamMemberFound = &memberFound

				if memberFound {
					// Check if date of birth matches
					dobMatches := member.DateOfBirth == app.DateOfBirth
					criteria.TeamLeaderDateOfBirthMatch = &dobMatches

					// Check role match
					roleMatches := member.Role == app.Role
					criteria.TeamMemberRoleMatches = &roleMatches

					// Team leader specific checks
					isTeamLeaderRole := app.Role == "team_leader"
					criteria.TeamLeaderRole = &isTeamLeaderRole

					if isTeamLeaderRole {
						teamLeaderFound := member.Role == "team_leader"
						criteria.TeamLeaderMemberFound = &teamLeaderFound
					}
				} else {
					// No matching member found
					noMatch := true
					criteria.NoMatchingTeamMember = &noMatch
				}
			}
		}
	}

	return criteria, nil
}

func PointGetAccountVerificationInfo(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	if r.Method != http.MethodGet {
		return nil, errors.New("method not allowed")
	}

	key := ps.ByName("key")
	if key == "" {
		return nil, errors.New("missing application key")
	}

	var realDB data.Database
	realDB = &db.RealDB{}

	acapp, err := realDB.GetAccountApplicationByKey(r.Context(), key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account application by key")
	}
	if acapp == nil {
		return &httpResult{
			ResponseType: http.StatusNotFound,
			Body:         `{"error":"account application not found"}`,
		}, nil
	}

	criteria, err := buildVerificationCriteria(r.Context(), realDB, *acapp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build verification criteria")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         criteria,
	}, nil
}

func PointVerifyAccountApplication(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed")
	}

	key := ps.ByName("key")
	if key == "" {
		return nil, errors.New("missing application key")
	}
	var realDB data.Database
	realDB = &db.RealDB{}

	acapp, err := realDB.GetAccountApplicationByKey(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account application by key")
	}
	if acapp == nil {
		return &httpResult{
			ResponseType: http.StatusNotFound,
			Body:         `{"error":"account application not found"}`,
		}, nil
	}

	if acapp.Role == "team_leader" {
		adminAccount := GetAdminAccount(ctx)
		if adminAccount != nil && adminAccount.Superadmin {
			logrus.Infof("Superadmin account found in context: %s", adminAccount.Username)
			err = registerAccountByApplication(ctx, realDB, *acapp)
			if err != nil {
				return nil, errors.Wrap(err, "failed to register account by application")
			}
			err = realDB.DeleteAccountApplicationByKey(ctx, key)
			if err != nil {
				return nil, errors.Wrap(err, "failed to delete account application by key")
			}
			return &httpResult{
				ResponseType: http.StatusOK,
				Body:         `{"status":"account verified and registered"}`,
			}, nil
		} else {
			return &httpResult{
				ResponseType: http.StatusForbidden,
				Body:         `{"error":"only superadmin can verify team leader applications"}`,
			}, nil
		}
	} else {
		account := GetAccount(ctx)
		if account != nil && account.Cilveks.Role == "team_leader" {
			logrus.Infof("Team leader account found in context: %s", account.Username)
			err = registerAccountByApplication(ctx, realDB, *acapp)
			if err != nil {
				return nil, errors.Wrap(err, "failed to register account by application")
			}
			err = realDB.DeleteAccountApplicationByKey(ctx, key)
			if err != nil {
				return nil, errors.Wrap(err, "failed to delete account application by key")
			}
			return &httpResult{
				ResponseType: http.StatusOK,
				Body:         `{"status":"account verified and registered"}`,
			}, nil
		} else {
			return &httpResult{
				ResponseType: http.StatusForbidden,
				Body:         `{"error":"only team leaders can verify member applications"}`,
			}, nil
		}
	}
}

// ----------------------------------------------------------------

// Account management
// ----------------------------------------------------------------

func PointGetAccounts(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	if r.Method != http.MethodGet {
		return nil, errors.New("method not allowed")
	}

	var realDB data.Database
	realDB = &db.RealDB{}

	adminAccount := GetAdminAccount(r.Context())
	if adminAccount != nil && adminAccount.Superadmin {
		logrus.Infof("Superadmin account found in context: %s", adminAccount.Username)
		accounts, err := realDB.GetAccounts(r.Context())
		if err != nil {
			return nil, errors.Wrap(err, "failed to get accounts")
		}
		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         accounts,
		}, nil
	}

	account := GetAccount(r.Context())
	if account != nil && account.Cilveks.Role == "team_leader" {
		logrus.Infof("Team leader account found in context: %s", account.Username)

		accounts, err := realDB.GetAccounts(r.Context())
		if err != nil {
			return nil, errors.Wrap(err, "failed to get accounts")
		}
		var teamAccounts []models.Account
		for _, acc := range accounts {
			if acc.Cilveks.ID == account.Cilveks.ID {
				teamAccounts = append(teamAccounts, acc)
			}
		}
		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         teamAccounts,
		}, nil
	}

	return &httpResult{
		ResponseType: http.StatusForbidden,
		Body:         `{"error":"forbidden"}`,
	}, nil
}

func PointPatchAccountByKey(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	logrus.Infof("PointPatchAccountByKey called %+v", ps)
	if r.Method != http.MethodPatch {
		return nil, errors.New("method not allowed")
	}

	Key := ps.ByName("key")
	if Key == "" {
		return nil, errors.New("missing account key")
	}

	logrus.Infof("Patching account with key: %s", Key)

	var realDB data.Database
	realDB = &db.RealDB{}

	// Get account by key from database
	existingAccount, err := realDB.GetAccountByKey(r.Context(), Key)
	if err != nil {
		return nil, errors.Wrap(err, "GetAccountByKey")
	}

	if existingAccount == nil {
		return nil, errors.New("account not found")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read body")
	}

	if err := json.Unmarshal(body, existingAccount); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}
	existingAccount.Cilveks.Key = Key

	adminaccount := GetAdminAccount(r.Context())
	if adminaccount != nil && adminaccount.Superadmin && adminaccount.Username != "" && adminaccount.Password != "" && adminaccount.Salt != "" {
		logrus.Infof("Admin account found in context: %+v", adminaccount)

		logrus.Infof("Received account update (superadmin): %+v", existingAccount)
		err = realDB.UpdateAccount(r.Context(), *existingAccount)
		if err != nil {
			return nil, errors.Wrap(err, "UpdateAccount")
		}
		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         "Account updated",
		}, nil
	}

	return &httpResult{
		ResponseType: http.StatusForbidden,
		Body:         `{"error":"forbidden"}`,
	}, nil
}

// ----------------------------------------------------------------

// Admin account management
// ----------------------------------------------------------------

func PointGetAdminAccounts(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	if r.Method != http.MethodGet {
		return nil, errors.New("method not allowed")
	}

	var realDB data.Database
	realDB = &db.RealDB{}

	adminAccount := GetAdminAccount(r.Context())
	if adminAccount != nil && adminAccount.Superadmin {
		logrus.Infof("Superadmin account found in context: %s", adminAccount.Username)
		admins, err := realDB.GetAllAdmins(r.Context())
		if err != nil {
			return nil, errors.Wrap(err, "failed to get admin accounts")
		}

		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         admins,
		}, nil
	}

	return &httpResult{
		ResponseType: http.StatusForbidden,
		Body:         `{"error":"forbidden"}`,
	}, nil
}

func PointPatchAdminAccount(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	if r.Method != http.MethodPatch {
		return nil, errors.New("method not allowed")
	}

	key := ps.ByName("key")
	if key == "" {
		return nil, errors.New("missing admin key")
	}

	var realDB data.Database
	realDB = &db.RealDB{}

	existingAdmin, err := realDB.GetAdminAccountByKey(r.Context(), key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get admin account by key")
	}
	if existingAdmin == nil {
		return &httpResult{
			ResponseType: http.StatusNotFound,
			Body:         `{"error":"admin account not found"}`,
		}, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read body")
	}
	defer r.Body.Close()

	var updatedAdmin models.AdminAccount
	if err := json.Unmarshal(body, &updatedAdmin); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}
	updatedAdmin.Key = key

	adminAccount := GetAdminAccount(r.Context())
	if adminAccount != nil && adminAccount.Superadmin {
		logrus.Infof("Admin update requested by superadmin: %s", adminAccount.Username)

		if updatedAdmin.Username != existingAdmin.Username {
			duplicateAdmin, err := realDB.GetAdminAccountByUsername(r.Context(), updatedAdmin.Username)
			if err != nil {
				return nil, errors.Wrap(err, "failed to check admin username")
			}
			if duplicateAdmin != nil {
				return &httpResult{
					ResponseType: http.StatusConflict,
					Body:         `{"error":"username already exists"}`,
				}, nil
			}
		}

		err = realDB.UpdateAdminAccount(r.Context(), updatedAdmin)
		if err != nil {
			return nil, errors.Wrap(err, "failed to update admin account")
		}

		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         `{"status":"admin account updated"}`,
		}, nil
	}

	return &httpResult{
		ResponseType: http.StatusForbidden,
		Body:         `{"error":"forbidden"}`,
	}, nil
}
