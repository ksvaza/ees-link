package httpapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/ksvaza/ees-link/db"
	"github.com/ksvaza/ees-link/fakedb"
	"github.com/ksvaza/ees-link/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func PointGetApplications(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	logrus.Infof("PointGetApplications called %+v", ps)
	if r.Method != http.MethodGet {
		return nil, errors.New("method not allowed")
	}

	ctx := r.Context()

	realDB := db.RealDB{}
	applications, err := realDB.GetAllApplications(ctx)

	adminaccount := GetAdminAccount(r.Context())
	if adminaccount != nil && adminaccount.Username != "" && adminaccount.Password != "" && adminaccount.Salt != "" {
		logrus.Infof("Admin account found in context: %+v", adminaccount)
		if err != nil {
			return nil, errors.Wrap(err, "GetAllApplications")
		}

		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         applications,
		}, nil
	}

	var a models.RegistrationFormData
	account := GetAccount(r.Context())
	if account == nil {
		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         "Unauthorized",
		}, nil
	}

	for _, application := range applications {
		if application.ID == account.Cilveks.ID {
			a = application
			break
		}
	}

	if a.ID == "" {
		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         "Application not found",
		}, nil
	}

	if account != nil && account.Cilveks.Role == "team_leader" && account.Username != "" && account.Password != "" && account.Salt != "" {
		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         a,
		}, nil
	} else {
		applicationRestricted := models.RegistrationFormDataRestricted{
			TeamName:    a.TeamName,
			Institution: a.Institution,
			MemberCount: len(a.Members),
			AppliedAt:   a.AppliedAt,
			Status:      a.Status,
		}

		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         applicationRestricted,
		}, nil
	}
}

func PointPostApplications(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	logrus.Infof("PointPostApplications called %+v", ps)
	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed")
	}

	var newApplicant models.RegistrationFormData
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read body")
	}

	if err := json.Unmarshal(body, &newApplicant); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}
	newApplicant.ID = uuid.New().String() // Assign a new unique ID to the application

	for i := range newApplicant.Members {
		newApplicant.Members[i].ID = newApplicant.ID
	}

	newApplicant.AppliedAt = time.Now()
	realDB := db.RealDB{}

	fmt.Printf("Registering new applicant:\n")

	// FS Backup
	err = fakedb.BackupApplication(newApplicant)
	if err != nil {
		return nil, errors.Wrap(err, "BackupApplication")
	}

	// New application handling logic here (e.g., save to database)

	logrus.Infof("Received new application: %+v", newApplicant)

	//realDB := fakedb.FSDatabase{}

	err = realDB.RegisterNewApplication(r.Context(), newApplicant)
	if err != nil {
		return nil, errors.Wrap(err, "AddApplication")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         "Application submitted",
	}, nil
}

func PointPatchApplicationByID(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	logrus.Infof("PointPatchApplicationByID called %+v", ps)
	if r.Method != http.MethodPatch {
		return nil, errors.New("method not allowed")
	}

	ID := ps.ByName("id")
	if ID == "" {
		return nil, errors.New("missing application ID")
	}

	logrus.Infof("Patching application with ID: %s", ID)

	realDB := db.RealDB{}

	// Get application by ID from database
	existingApplication, err := realDB.GetApplicationByID(r.Context(), ID)
	if err != nil {
		return nil, errors.Wrap(err, "GetApplicationByID")
	}

	if existingApplication == nil {
		return nil, errors.New("application not found")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read body")
	}

	adminaccount := GetAdminAccount(r.Context())
	if adminaccount != nil && adminaccount.Superadmin && adminaccount.Username != "" && adminaccount.Password != "" && adminaccount.Salt != "" {
		logrus.Infof("Admin account found in context: %+v", adminaccount)

		if err := json.Unmarshal(body, existingApplication); err != nil {
			return nil, errors.Wrap(err, "Unmarshal")
		}
		existingApplication.ID = ID

		logrus.Infof("Received application update (superadmin): %+v", existingApplication)
		err = realDB.UpdateApplication(r.Context(), *existingApplication)
		if err != nil {
			return nil, errors.Wrap(err, "UpdateApplication")
		}

		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         "Application updated",
		}, nil
	}

	account := GetAccount(r.Context())
	if account != nil && account.Username != "" && account.Password != "" && account.Salt != "" {
		if account.Cilveks.ID == ID {
			logrus.Infof("Received application update from team member: %+v", account)
			if err := json.Unmarshal(body, existingApplication); err != nil {
				return nil, errors.Wrap(err, "Unmarshal")
			}
			existingApplication.ID = ID

			err = realDB.UpdateApplication(r.Context(), *existingApplication)
			if err != nil {
				return nil, errors.Wrap(err, "UpdateApplication")
			}
			return &httpResult{
				ResponseType: http.StatusOK,
				Body:         "Application updated",
			}, nil
		}
	}

	return &httpResult{
		ResponseType: http.StatusForbidden,
		Body:         `{"error":"forbidden"}`,
	}, nil
}

// func PointGetApplicationsRestricted(r *http.Request, ps httprouter.Params) (*httpResult, error) {
// 	logrus.Infof("PointGetApplicationsRestricted called %+v", ps)
// 	if r.Method != http.MethodGet {
// 		return nil, errors.New("method not allowed")
// 	}

// 	//realDB := fakedb.FSDatabase{}

// 	var realDB data.Database
// 	realDB = &db.RealDB{}

// 	//db.GetAllApplicants()

// 	applications, err := realDB.GetAllApplications(r.Context())
// 	if err != nil {
// 		return nil, errors.Wrap(err, "GetAllApplications")
// 	}

// 	restrictedApplications := make([]models.RegistrationFormDataRestricted, len(applications))
// 	for i, app := range applications {
// 		restrictedApplications[i] = models.RegistrationFormDataRestricted{
// 			TeamName:    app.TeamName,
// 			Institution: app.Institution,
// 			MemberCount: len(app.Members),
// 			AppliedAt:   app.AppliedAt,
// 			Status:      app.Status,
// 		}
// 	}

// 	return &httpResult{
// 		ResponseType: http.StatusOK,
// 		Body:         restrictedApplications,
// 	}, nil
// }
