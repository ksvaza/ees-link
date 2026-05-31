package httpapi

import (
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

// ----------------------------------------------------------------

// Race management
// ----------------------------------------------------------------

func PointGetAdminSettings(r *http.Request, ps httprouter.Params) (*httpResult, error) {
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

	settings, err := realDB.GetAdminSettings(r.Context())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get admin settings")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         settings,
	}, nil
}

func PointPostAdminSettings(r *http.Request, ps httprouter.Params) (*httpResult, error) {
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

	var newSettings models.AdminSettings
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read body")
	}

	if err := json.Unmarshal(body, &newSettings); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	err = realDB.UpdateAdminSettings(r.Context(), newSettings)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update admin settings")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         `{"status":"admin settings updated"}`,
	}, nil
}
