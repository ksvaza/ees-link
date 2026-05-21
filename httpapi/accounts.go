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
			Key:        adminaccount.Key,
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
			// Key: account.Key,
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

func registerAccountByApplication(ctx context.Context, realDB data.Database, app models.AccountApplication) error { // applicant guarranteed != "team_leader"
	// Šeit var sarakstīt, kas ir un nav obliātie lauki.
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
	var teamID string = ""
	newAccount := models.Account{
		Cilveks: models.TeamMember{
			Key:         app.Key,
			FullName:    app.FullName,
			DateOfBirth: app.DateOfBirth,
			Role:        app.Role,
			ID:          "",
		},
		Password:               app.Password,
		Username:               app.Username,
		Email:                  app.Email,
		PhoneNumber:            app.PhoneNumber,
		Salt:                   app.Salt,
		PendingTeamID:          "",
		Verified:               false,
		EducationalInstitution: "",
		ClassOrYear:            "",
	}

	admin := GetAdminAccount(ctx)
	account := GetAccount(ctx)
	if newAccount.Cilveks.FullName == account.Cilveks.FullName && newAccount.Cilveks.DateOfBirth == account.Cilveks.DateOfBirth {
		// lai iet pāris mājas tālāk un nelien, kur nevajag
		return errors.New("account registration logic error: applicant information matches currently authenticated account - cannot register")
	}

	team, err := realDB.GetApplicationByTeamName(ctx, app.TeamName)
	if err != nil {
		return errors.Wrap(err, "failed to get application by team name")
	}
	if team == nil {
		if app.TeamName == "" {
			logrus.Infof("Registering account without team association")
			newAccount.PendingTeamID = ""
		} else {
			return errors.New("team not found")
		}
	} else {
		teamID = team.ID
		newAccount.PendingTeamID = teamID

		if account != nil && account.Cilveks.Role == "team_leader" {
			if account.Cilveks.ID != teamID {
				return errors.New("account registration logic error: authenticated team leader does not match team application - cannot register")
			}
		}

		if (admin != nil && admin.Superadmin) || (account != nil && account.Cilveks.Role == "team_leader") {
			for _, member := range team.Members {
				if member.FullName == app.FullName && member.DateOfBirth == app.DateOfBirth {
					newAccount.EducationalInstitution = member.EducationalInstitution
					newAccount.ClassOrYear = member.ClassOrYear
					logrus.Infof("Found matching team member in application for account registration: %s (team ID: %s)", app.FullName, teamID)
					break
				}
			}

			logrus.Infof("Registering account with team association: %s (team ID: %s)", app.TeamName, teamID)
		}
	}

	err = realDB.RegisterNewAccount(ctx, newAccount)
	if err != nil {
		return errors.Wrap(err, "failed to register new account")
	}

	if newAccount.Cilveks.Role != "team_leader" {
		if (admin != nil && admin.Superadmin) || (account != nil && account.Cilveks.Role == "team_leader") {
			//
			//
			// add account to team
			//
			// id/pendingteamid maģija
			//

			// jāatjauno konta info
			if newAccount.Cilveks.ID != newAccount.PendingTeamID {
				return errors.New("account registration logic error: pending team ID does not match team member ID")
			}

			newAccount.Verified = true

			err = realDB.UpdateAccount(ctx, newAccount)
		} else {
			logrus.Infof("Account registered with team association but pending verification: %s (team ID: %s)", app.FullName, teamID)
		}
	} else {
		// For team leaders, we require manual verification before they are fully registered and associated with the team, so we do not automatically update the account to registered=true here. The verification process will handle that.
		logrus.Infof("Team leader account registered but pending verification: %s (team ID: %s)", app.FullName, teamID)
	}

	return nil
}

func registerTeamLeaderAccountByApplication(ctx context.Context, realDB data.Database, app models.AccountApplication, teamIDcheck string) error {
	if app.Role != "team_leader" {
		return errors.New("application role is not team_leader")
	}

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
	var teamID string = ""
	newAccount := models.Account{
		Cilveks: models.TeamMember{
			Key:         app.Key,
			FullName:    app.FullName,
			DateOfBirth: app.DateOfBirth,
			Role:        app.Role,
			ID:          "",
		},
		Password:               app.Password,
		Username:               app.Username,
		Email:                  app.Email,
		PhoneNumber:            app.PhoneNumber,
		Salt:                   app.Salt,
		PendingTeamID:          "",
		Verified:               false,
		EducationalInstitution: "",
		ClassOrYear:            "",
	}

	// Get the team application by name
	teamApp, err := realDB.GetApplicationByTeamName(ctx, app.TeamName)
	if err != nil {
		return errors.Wrap(err, "failed to get team application")
	}
	if teamApp == nil {
		return errors.New("team application not found")
	}

	admin := GetAdminAccount(ctx)
	if admin == nil || !admin.Superadmin {
		// Validate that the team application ID matches the expected ID
		if teamApp.ID != teamIDcheck {
			return errors.New("team ID mismatch - verification failed")
		}
	}
	teamID = teamApp.ID
	newAccount.PendingTeamID = teamID

	for _, member := range teamApp.Members {
		if member.FullName == app.FullName && member.DateOfBirth == app.DateOfBirth {
			if admin == nil || !admin.Superadmin {
				if member.Role != "team_leader" {
					return errors.New("member is not team leader - cannot register as team leader")
				}
			}
			newAccount.EducationalInstitution = member.EducationalInstitution
			newAccount.ClassOrYear = member.ClassOrYear
			logrus.Infof("Found matching team member in application for account registration: %s (team ID: %s)", app.FullName, teamID)
			break
		}
	}

	logrus.Infof("Registering team leader account with team association: %s (team ID: %s)", app.TeamName, teamApp.ID)

	err = realDB.RegisterNewAccount(ctx, newAccount)
	if err != nil {
		return errors.Wrap(err, "failed to register new account")
	}

	//
	//
	// add account to team
	//
	// id/pendingteamid maģija
	//

	// jāatjauno konta info
	if newAccount.Cilveks.ID != newAccount.PendingTeamID {
		return errors.New("account registration logic error: pending team ID does not match team member ID")
	}

	newAccount.Verified = true

	err = realDB.UpdateAccount(ctx, newAccount)

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

	if pieteikums.Role == "team_leader" {
		// https://e-es.lv/api/register/1234567890
		preid := ps.ByName("uniqueID")
		id, err := url.PathUnescape(preid)
		if err != nil {
			logrus.WithError(err).Warnf("Failed to unescape team leader verification ID: %v", err)
		}
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

	err = registerAccountByApplication(r.Context(), realDB, pieteikums)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register account by application")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         `{"status":"success"}`,
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

	var teamExists bool
	var memberFound bool
	var dobMatches bool
	var roleMatches bool
	var teamLeaderFound bool

	if app.TeamName != "" {
		teamApp, err := realDB.GetApplicationByTeamName(ctx, app.TeamName)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get team application")
		}
		teamExists = teamApp != nil
		if teamExists {
			member := teamApp.FindTeamMemberByFullName(app.FullName)
			memberFound = member != nil
			if memberFound {
				dobMatches = member.DateOfBirth == app.DateOfBirth
				roleMatches = member.Role == app.Role
				teamLeaderFound = member.Role == "team_leader"
			}
		}
	}

	// Determine whether the application can be registered in principle
	if criteria.RequiredFieldsPresent && criteria.UsernameAvailable {
		if app.TeamName == "" {
			criteria.CanRegister = app.Role != "team_leader"
		} else if app.Role == "team_leader" {
			criteria.CanRegister = teamExists && memberFound && teamLeaderFound && dobMatches
		} else {
			criteria.CanRegister = teamExists
		}
	}

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
			criteria.TeamApplicationExists = &teamExists
			criteria.TeamMemberFound = &memberFound
			if memberFound {
				criteria.TeamLeaderDateOfBirthMatch = &dobMatches
				criteria.TeamMemberRoleMatches = &roleMatches

				isTeamLeaderRole := app.Role == "team_leader"
				criteria.TeamLeaderRole = &isTeamLeaderRole

				if isTeamLeaderRole {
					criteria.TeamLeaderMemberFound = &teamLeaderFound
				}
			} else {
				noMatch := teamExists
				criteria.NoMatchingTeamMember = &noMatch
			}
		}
	}

	return criteria, nil
}

func PointGetAccountVerificationInfo(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	logrus.Infof("PointGetAccountVerificationInfo called %+v", ps)
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

	logrus.Infof("Received application data for verification info: %+v", pieteikums)

	var realDB data.Database
	realDB = &db.RealDB{}

	criteria, err := buildVerificationCriteria(r.Context(), realDB, pieteikums)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build verification criteria")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         criteria,
	}, nil
}

func PointGetAccountVerificationInfoByKey(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	if r.Method != http.MethodGet {
		return nil, errors.New("method not allowed")
	}

	prekey := ps.ByName("key")
	if prekey == "" {
		return nil, errors.New("missing application key")
	}

	key, err := url.PathUnescape(prekey)
	if err != nil {
		return nil, errors.Wrap(err, "Unescape application key")
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

func PointVerifyAccountApplicationByKey(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed")
	}

	prekey := ps.ByName("key")
	if prekey == "" {
		return nil, errors.New("missing application key")
	}

	key, err := url.PathUnescape(prekey)
	if err != nil {
		return nil, errors.Wrap(err, "Unescape application key")
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

	adminAccount := GetAdminAccount(ctx)
	if acapp.Role == "team_leader" {
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
		}
		return &httpResult{
			ResponseType: http.StatusForbidden,
			Body:         `{"error":"only superadmin can verify team leader applications"}`,
		}, nil
	} else {
		account := GetAccount(ctx)
		if (account != nil && account.Cilveks.Role == "team_leader") || (adminAccount != nil && adminAccount.Superadmin) {
			if adminAccount != nil && adminAccount.Superadmin {
				logrus.Infof("Superadmin account found in context: %s", adminAccount.Username)
			} else {
				logrus.Infof("Team leader account found in context: %s", account.Username)
			}
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
		}
		return &httpResult{
			ResponseType: http.StatusForbidden,
			Body:         `{"error":"only team leaders can verify member applications"}`,
		}, nil
	}
} // vairs neizmantots

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
	accounts, err := realDB.GetAccounts(r.Context())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get accounts")
	}

	if account != nil && account.Cilveks.Role == "team_leader" {
		logrus.Infof("Team leader account found in context: %s", account.Username)

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
	} else if account != nil {
		logrus.Infof("Regular account found in context: %s", account.Username)
		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         []models.Account{*account},
		}, nil
	}

	return &httpResult{
		ResponseType: http.StatusForbidden,
		Body:         `{"error":"forbidden"}`,
	}, nil
}

func PointGetUnregisteredAccounts(r *http.Request, ps httprouter.Params) (*httpResult, error) {
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
		unvaccounts := []models.Account{}
		for _, acc := range accounts {
			if !acc.Verified {
				unvaccounts = append(unvaccounts, acc)
			}
		}
		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         unvaccounts,
		}, nil
	}

	account := GetAccount(r.Context())
	accounts, err := realDB.GetAccounts(r.Context())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get accounts")
	}
	unvaccounts := []models.Account{}
	for _, acc := range accounts {
		if !acc.Verified {
			unvaccounts = append(unvaccounts, acc)
		}
	}

	if account != nil && account.Cilveks.Role == "team_leader" {
		logrus.Infof("Team leader account found in context: %s", account.Username)

		var teamAccounts []models.Account
		for _, acc := range unvaccounts {
			if acc.Cilveks.ID == account.Cilveks.ID {
				teamAccounts = append(teamAccounts, acc)
			}
		}
		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         teamAccounts,
		}, nil
	} else if account != nil {
		logrus.Infof("Regular account found in context: %s", account.Username)
		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         []models.Account{*account},
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

	preKey := ps.ByName("key")
	if preKey == "" {
		return nil, errors.New("missing account key")
	}

	key, err := url.PathUnescape(preKey)
	if err != nil {
		return nil, errors.Wrap(err, "Unescape account key")
	}

	logrus.Infof("Patching account with key: %s", key)

	var realDB data.Database
	realDB = &db.RealDB{}

	// Get account by key from database
	existingAccount, err := realDB.GetAccountByKey(r.Context(), key)
	if err != nil {
		return nil, errors.Wrap(err, "GetAccountByKey")
	}

	if existingAccount == nil {
		return nil, errors.New("account not found")
	}

	prevUsername := existingAccount.Username

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read body")
	}

	if err := json.Unmarshal(body, existingAccount); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}
	existingAccount.Cilveks.Key = key

	adminaccount := GetAdminAccount(r.Context())
	if adminaccount != nil && adminaccount.Superadmin && adminaccount.Username != "" && adminaccount.Password != "" && adminaccount.Salt != "" {
		logrus.Infof("Admin account found in context: %+v", adminaccount)

		logrus.Infof("Received account update (superadmin): %+v", existingAccount)

		if existingAccount.Username != prevUsername {
			duplicateAccount, err := realDB.GetAccountByUsername(r.Context(), existingAccount.Username)
			if err != nil {
				return nil, errors.Wrap(err, "failed to check account username")
			}
			if duplicateAccount != nil {
				return &httpResult{
					ResponseType: http.StatusConflict,
					Body:         `{"error":"username already exists"}`,
				}, nil
			}
		}

		err = realDB.UpdateAccount(r.Context(), *existingAccount)
		if err != nil {
			return nil, errors.Wrap(err, "UpdateAccount")
		}
		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         "Account updated",
		}, nil
	}

	account := GetAccount(r.Context())
	if account != nil && account.Cilveks.Key == key {
		logrus.Infof("Account found in context: %+v", account)

		logrus.Infof("Received account update (account owner): %+v", existingAccount)

		if existingAccount.Username != prevUsername {
			duplicateAccount, err := realDB.GetAccountByUsername(r.Context(), existingAccount.Username)
			if err != nil {
				return nil, errors.Wrap(err, "failed to check account username")
			}
			if duplicateAccount != nil {
				return &httpResult{
					ResponseType: http.StatusConflict,
					Body:         `{"error":"username already exists"}`,
				}, nil
			}
		}

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

func PointVerifyAccountByKey(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed")
	}

	prekey := ps.ByName("key")
	if prekey == "" {
		return nil, errors.New("missing application key")
	}

	key, err := url.PathUnescape(prekey)
	if err != nil {
		return nil, errors.Wrap(err, "Unescape application key")
	}

	var realDB data.Database
	realDB = &db.RealDB{}

	acapp, err := realDB.GetAccountByKey(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account by key")
	}
	if acapp == nil {
		return &httpResult{
			ResponseType: http.StatusNotFound,
			Body:         `{"error":"account not found"}`,
		}, nil
	}

	adminAccount := GetAdminAccount(ctx)
	if acapp.Cilveks.Role == "team_leader" {
		if adminAccount != nil && adminAccount.Superadmin {
			logrus.Infof("Superadmin account found in context: %s", adminAccount.Username)

			//
			//
			// add account to team
			//
			// id/pendingteamid maģija
			//

			// jāatjauno konta info
			if acapp.Cilveks.ID != acapp.PendingTeamID {
				return &httpResult{
					ResponseType: http.StatusBadRequest,
					Body:         `{"error":"failed to add account to team"}`,
				}, nil
			}

			acapp.Verified = true

			err = realDB.UpdateAccount(ctx, *acapp)
			if err != nil {
				return nil, errors.Wrap(err, "failed to update account")
			}

			return &httpResult{
				ResponseType: http.StatusOK,
				Body:         `{"status":"account verified and registered"}`,
			}, nil
		}
		return &httpResult{
			ResponseType: http.StatusForbidden,
			Body:         `{"error":"only superadmin can verify team leader applications"}`,
		}, nil
	} else {
		account := GetAccount(ctx)
		if (account != nil && account.Cilveks.Role == "team_leader") || (adminAccount != nil && adminAccount.Superadmin) {
			if adminAccount != nil && adminAccount.Superadmin {
				logrus.Infof("Superadmin account found in context: %s", adminAccount.Username)
			} else {
				logrus.Infof("Team leader account found in context: %s", account.Username)
			}

			//
			//
			// add account to team
			//
			// id/pendingteamid maģija
			//

			// jāatjauno konta info
			if acapp.Cilveks.ID != acapp.PendingTeamID {
				return &httpResult{
					ResponseType: http.StatusBadRequest,
					Body:         `{"error":"failed to add account to team"}`,
				}, nil
			}

			acapp.Verified = true

			err = realDB.UpdateAccount(ctx, *acapp)

			return &httpResult{
				ResponseType: http.StatusOK,
				Body:         `{"status":"account verified and registered"}`,
			}, nil
		}
		return &httpResult{
			ResponseType: http.StatusForbidden,
			Body:         `{"error":"only team leaders or superadmins can verify member applications"}`,
		}, nil
	}
}

func PointRegisterAccount(r *http.Request, ps httprouter.Params) (*httpResult, error) {
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

	err = realDB.RegisterNewAccountApplication(r.Context(), pieteikums)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register new account application")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         `{"status":"pending verification"}`,
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
	admins, err := realDB.GetAllAdmins(r.Context())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get admin accounts")
	}

	if adminAccount != nil && adminAccount.Superadmin {
		logrus.Infof("Superadmin account found in context: %s", adminAccount.Username)

		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         admins,
		}, nil
	} else if adminAccount != nil {
		logrus.Infof("Admin account found in context: %s", adminAccount.Username)
		// Return only the authenticated admin's own account info
		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         []models.AdminAccount{*adminAccount},
		}, nil
	}

	return &httpResult{
		ResponseType: http.StatusForbidden,
		Body:         `{"error":"forbidden"}`,
	}, nil
}

func PointPatchAdminAccountByKey(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	if r.Method != http.MethodPatch {
		return nil, errors.New("method not allowed")
	}

	preKey := ps.ByName("key")
	if preKey == "" {
		return nil, errors.New("missing admin key")
	}

	key, err := url.PathUnescape(preKey)
	if err != nil {
		return nil, errors.Wrap(err, "Unescape admin key")
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
	prevUsername := existingAdmin.Username

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read body")
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &existingAdmin); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}
	existingAdmin.Key = key

	adminAccount := GetAdminAccount(r.Context())
	if adminAccount != nil && adminAccount.Superadmin {
		logrus.Infof("Admin update requested by superadmin: %s", adminAccount.Username)

		if existingAdmin.Username != prevUsername {
			duplicateAdmin, err := realDB.GetAdminAccountByUsername(r.Context(), existingAdmin.Username)
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

		err = realDB.UpdateAdminAccount(r.Context(), *existingAdmin)
		if err != nil {
			return nil, errors.Wrap(err, "failed to update admin account")
		}

		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         `{"status":"admin account updated"}`,
		}, nil
	} else if adminAccount != nil && adminAccount.Key == key {
		logrus.Infof("Admin update requested by account owner: %s", adminAccount.Username)
		if existingAdmin.Username != prevUsername {
			duplicateAdmin, err := realDB.GetAdminAccountByUsername(r.Context(), existingAdmin.Username)
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
		err = realDB.UpdateAdminAccount(r.Context(), *existingAdmin)
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
