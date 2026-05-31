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

// PointPostPoints

func PointPostPoints(r *http.Request, ps httprouter.Params) (*httpResult, error) {
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

	var points models.Points
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read body")
	}

	if err := json.Unmarshal(body, &points); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	err = realDB.SavePoints(r.Context(), points)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save points")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         `{"status":"points saved"}`,
	}, nil
}
