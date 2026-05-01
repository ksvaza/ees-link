package fakedb

import (
	"context"
	"fmt"

	"github.com/ksvaza/ees-link/models"
	"github.com/sirupsen/logrus"
)

type FSDatabase struct{}

func BackupApplication(a models.RegistrationFormData) error {
	return createApplicationFile(a)
}

func (db *FSDatabase) TestHealthiness(ctx context.Context) error {
	return testExistance()
}

func (db *FSDatabase) GetAllApplications(ctx context.Context) ([]models.RegistrationFormData, error) {
	return getAllApplications()
}

func (db *FSDatabase) GetApplicationByID(ctx context.Context, id string) (models.RegistrationFormData, error) {
	app, err := getApplicationByID(id)
	if err != nil {
		return models.RegistrationFormData{}, err
	}
	if app == nil {
		return models.RegistrationFormData{}, fmt.Errorf("application with ID '%s' not found", id)
	}
	return *app, nil
}

func (db *FSDatabase) RegisterNewApplication(ctx context.Context, newApplicant models.RegistrationFormData) error {
	existing, err := getAllApplications()
	if err != nil {
		return err
	}
	for _, app := range existing {
		if app.TeamName == newApplicant.TeamName {
			return fmt.Errorf("application with team name '%s' already exists", newApplicant.TeamName)
		}
	}
	return addApplication(newApplicant)
}

func (db *FSDatabase) UpdateApplication(ctx context.Context, updatedApplicant models.RegistrationFormData) error {
	existing, err := getAllApplications()
	if err != nil {
		return err
	}
	found := false
	for _, app := range existing {
		if app.ID == updatedApplicant.ID {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("application with ID '%s' not found", updatedApplicant.ID)
	}
	return addApplication(updatedApplicant)
}

func (db *FSDatabase) RegisterNewAccount(ctx context.Context, newAccount models.Account) error {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil
}

func (db *FSDatabase) RegisterNewAccountApplication(ctx context.Context, newAccountApplication models.AccountApplication) (models.AccountApplication, error) {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return newAccountApplication, nil
}

func (db *FSDatabase) GetAccountByUsername(ctx context.Context, username string) (models.Account, error) {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return models.Account{}, nil
}

func (db *FSDatabase) GetAccountByDateOfBirth(ctx context.Context, dateOfBirth string) (models.Account, error) {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return models.Account{}, nil
}

func (db *FSDatabase) GetAccountByFullname(ctx context.Context, fullname string) (models.Account, error) {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return models.Account{}, nil
}

func (db *FSDatabase) RegisterNewAdmin(ctx context.Context, newAccount models.AdminAccount) error {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil
}

func (db *FSDatabase) GetAccountApplications(ctx context.Context) ([]models.AccountApplication, error) {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return []models.AccountApplication{}, nil
}
