package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ksvaza/ees-link/db"
	"github.com/ksvaza/ees-link/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func PointRegister(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	logrus.Infof("PointRegister called %+v", ps)
	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed")
	}

	//id := ps.ByName("uniqueID") // e-es.lv/api/register?uniqueID=1234

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

	realDB := db.RealDB{}

	err = realDB.RegisterNewAccountApplication(context.Background(), pieteikums)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register new account application")
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
	salt, err := GenerateSalt(16)
	if err != nil {
		return errors.Wrap(err, "failed to generate salt")
	}

	realDB := db.RealDB{}

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

	fmt.Printf("%+v\n", admin)

	err = realDB.RegisterNewAdmin(context.Background(), admin)
	if err != nil {
		return errors.Wrap(err, "failed to register new admin")
	}

	return nil
}

func PointLogin(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	logrus.Infof("PointLogin called %+v", ps)
	if r.Method != http.MethodPost {
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
			Cilveks:   account.Cilveks,
			Username:  account.Username,
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

func PointGetAccountApplications(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	if r.Method != http.MethodGet {
		return nil, errors.New("method not allowed")
	}

	realDB := db.RealDB{}

	applications, err := realDB.GetAllApplications(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all applications")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         applications,
	}, nil
}

func PointReceiveVerifiedAccounts(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed")
	}

	var verifiedAccounts []models.Account

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read request body")
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &verifiedAccounts); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal request body")
	}

	realDB := db.RealDB{}
	for _, account := range verifiedAccounts {
		err := realDB.RegisterNewAccount(context.Background(), account)
		if err != nil {
			return nil, errors.Wrap(err, "failed to save verified accounts")
		}
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         "Accounts inserted into DB",
	}, nil
}
