package httpapi

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ksvaza/ees-link/data"
	"github.com/ksvaza/ees-link/db"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func PointGetLeaderboard(r *http.Request, ps httprouter.Params) (*httpResult, error) {
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

	ageGroup := ps.ByName("ageGroup")
	if ageGroup == "" {
		return nil, errors.New("missing age group")
	}

	leaderboard, err := realDB.BuildLeaderboard(r.Context(), ageGroup)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get leaderboard")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         leaderboard,
	}, nil
}
