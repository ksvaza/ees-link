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

	if pieteikums.Role == "team_leader" {
		id := ps.ByName("uniqueID") // http://e-es.lv/api/register?uniqueID=1234567890
		if id != "" {
			team, err := realDB.GetApplicationByID(context.Background(), id)
			if err != nil {
				logrus.WithError(err).Error("Failed to find team by ID")
			} else {
				member := team.FindTeamMemberByFullName(pieteikums.FullName)
				if member != nil {
					if member.Role == "team_leader" {
						if pieteikums.DateOfBirth != member.DateOfBirth {
							logrus.Infof("Nesakrīt dzimšanas dienas datumi kontam ar pieteikumu \"%s\" pret \"%s\"", pieteikums.DateOfBirth, member.DateOfBirth)
						} else {
							var account models.Account
							account.Cilveks = *member
							account.Salt, err = GenerateSalt(16)
							if err != nil {
								logrus.WithError(err).Error("Neizdevās saģenerēt sāli.")
							}
							account.Password = hashPassword(pieteikums.Password, account.Salt)
							account.Username = pieteikums.Username
							account.Email = pieteikums.Email
							account.PhoneNumber = pieteikums.PhoneNumber

							err = realDB.RegisterNewAccount(context.Background(), account)
							if err != nil {
								logrus.WithError(err).Error("Neizdevās piereģistrēt komandas līderi automātiski")
							} else {
								return &httpResult{
									ResponseType: http.StatusOK,
									Body:         `{"status":"accepted"}`,
								}, nil
							}
						}
					} else {
						logrus.Infof("Piesakoties kā līderim, šeit \"%+v\" cilvēka vārds nesakrīt ar komandas \"%s\" līdera vārdu", pieteikums, member.FullName)
					}
				} else {
					logrus.Infof("Neizdevās atrast komandā \"%s\" cilvēku ar vārdu \"%s\"", team.TeamName, pieteikums.FullName)
				}
			}
		}
	}

	err = realDB.RegisterNewAccountApplication(context.Background(), pieteikums)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register new account application")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         `{"status":"waiting"}`,
	}, nil
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
		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         adminaccount,
		}, nil
	}
	account := GetAccount(r.Context())
	if account != nil && account.Username != "" && account.Password != "" && account.Salt != "" {
		logrus.Infof("Authenticated account: %+v", account)
		return &httpResult{
			ResponseType: http.StatusOK,
			Body:         account,
		}, nil
	}

	return &httpResult{
		ResponseType: http.StatusUnauthorized,
		Body:         `{"error":"unauthorized"}`,
	}, nil
}

func PointGetAccountApplications(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	logrus.Infof("PointGetAccountApplications called %+v", ps)
	if r.Method != http.MethodGet {
		return nil, errors.New("method not allowed")
	}

	realDB := db.RealDB{}
	accountApplications, err := realDB.GetAccountApplications(context.Background())

	if err != nil {
		return nil, errors.Wrap(err, "failed to get account applications")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         accountApplications,
	}, nil
}

func PointReceiveVerifiedAccounts(r *http.Request, ps httprouter.Params) (*httpResult, error) {
	logrus.Infof("PointReceiveVerifiedAccounts called %+v", ps)
	if r.Method != http.MethodPost {
		return nil, errors.New("method not allowed")
	}

	var accounts []models.Account
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read body")
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &accounts); err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	realDB := db.RealDB{}
	for _, account := range accounts {
		err := realDB.RegisterNewAccount(context.Background(), account)
		if err != nil {
			return nil, errors.Wrap(err, "failed to save verified accounts")
		}
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         accounts,
	}, nil
}
