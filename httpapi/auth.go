package httpapi

import (
	"context"
	"encoding/json"
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

	if pieteikums.TeamName == "" {
		if pieteikums.Username == "" || pieteikums.Password == "" || pieteikums.Email == "" || pieteikums.FullName == "" || pieteikums.DateOfBirth == "" {
			return &httpResult{
				ResponseType: http.StatusBadRequest,
				Body:         `{"error":"missing required fields"}`,
			}, nil
		}

		existing, err := realDB.GetAccountByUsername(r.Context(), pieteikums.Username)
		if err != nil {
			return nil, errors.Wrap(err, "failed to check username")
		}
		if existing != nil {
			return &httpResult{
				ResponseType: http.StatusConflict,
				Body:         `{"error":"username already exists"}`,
			}, nil
		}

		salt, err := GenerateSalt(16)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate salt")
		}

		account := models.Account{
			Cilveks: models.TeamMember{
				FullName:    pieteikums.FullName,
				DateOfBirth: pieteikums.DateOfBirth,
			},
			Username:    pieteikums.Username,
			Password:    hashPassword(pieteikums.Password, salt),
			Email:       pieteikums.Email,
			PhoneNumber: pieteikums.PhoneNumber,
			Salt:        salt,
		}

		err = realDB.RegisterNewAccount(r.Context(), account)
		if err != nil {
			return nil, errors.Wrap(err, "failed to register new account")
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
			team, err := realDB.GetApplicationByID(r.Context(), id)
			if err != nil {
				logrus.WithError(err).Error("Failed to find team by ID")
			} else {
				member := team.FindTeamMemberByFullName(pieteikums.FullName)
				if member != nil {
					if member.Role == "team_leader" {
						if pieteikums.Username == "" || pieteikums.Password == "" || pieteikums.Email == "" {
							return &httpResult{
								ResponseType: http.StatusBadRequest,
								Body:         `{"error":"missing required fields for team leader"}`,
							}, nil
						}

						existing, err := realDB.GetAccountByUsername(r.Context(), pieteikums.Username)
						if err != nil {
							return nil, errors.Wrap(err, "failed to check username")
						}
						if existing != nil {
							return &httpResult{
								ResponseType: http.StatusConflict,
								Body:         `{"error":"username already exists"}`,
							}, nil
						}

						if pieteikums.DateOfBirth != member.DateOfBirth {
							logrus.Infof("Nesakrīt dzimšanas dienas datumi kontam ar pieteikumu \"%s\" pret \"%s\"", pieteikums.DateOfBirth, member.DateOfBirth)
						} else {
							var account models.Account
							account.Cilveks = *member
							account.Salt, err = GenerateSalt(16)
							if err != nil {
								logrus.WithError(err).Error("Neizdevās saģenerēt sāli.")
							} else {
								// paroli nevajadzētu sūtīt kā parastu teksu, bet tas tā
								account.Password = hashPassword(pieteikums.Password, account.Salt)
								account.Username = pieteikums.Username
								account.Email = pieteikums.Email
								account.PhoneNumber = pieteikums.PhoneNumber

								err = realDB.RegisterNewAccount(r.Context(), account)
								if err != nil {
									logrus.WithError(err).Error("Neizdevās piereģistrēt komandas līderi automātiski")
								} else {
									return &httpResult{
										ResponseType: http.StatusOK,
										Body:         `{"status":"accepted"}`,
									}, nil
								}
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

	err = realDB.RegisterNewAccountApplication(r.Context(), pieteikums)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register new account application")
	}

	return &httpResult{
		ResponseType: http.StatusOK,
		Body:         `{"status":"pending verification"}`,
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
			Cilveks:  account.Cilveks,
			Username: account.Username,
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
