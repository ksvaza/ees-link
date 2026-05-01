package httpapi

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ksvaza/ees-link/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func PointRegister(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	logrus.Infof("PointRegister called %+v", ps)
	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed")
	}

	id := ps.ByName("uniqueID")

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

	err = RegisterApplication(pieteikums, id)
	if err != nil {
		return nil, errors.Wrap(err, "Register account")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         `{"status":"success"}`,
	}, nil
}

func RegisterApplication(pieteikums models.AccountApplication, id string) error {
	// TODO: Implement application registration logic here

	//realDB := db.RealDB{}
	return nil
}

func RegisterAdminAccount(admin models.AdminAccount) error {
	//admin.Salt = GenerateSalt()
	// if admin exists return nil, otherwise create new admin account and return nil
	// if errror occurs, return error

	// Bračiņ, vajag funkciju datubāzē šim
	return nil
}

func PointLogin(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	logrus.Infof("PointLogin called %+v", ps)
	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed")
	}
	// TODO: implement login
	return nil, errors.New("not implemented")
}
