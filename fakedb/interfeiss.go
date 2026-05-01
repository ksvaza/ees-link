package fakedb

import (
	"context"
	"fmt"

	"github.com/ksvaza/ees-link/models"
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
	return nil
}

func (db *FSDatabase) RegisterNewAccountApplication(ctx context.Context, newAccountApplication models.AccountApplication) (models.AccountApplication, error) {
	return newAccountApplication, nil
}

func (db *FSDatabase) GetAccountByUsername(ctx context.Context, username string) (models.Account, error) {
	return models.Account{}, nil
}

func (db *FSDatabase) GetAccountByDateOfBirth(ctx context.Context, dateOfBirth string) (models.Account, error) {
	return models.Account{}, nil
}

func (db *FSDatabase) GetAccountByFullname(ctx context.Context, fullname string) (models.Account, error) {
	return models.Account{}, nil
}

func (db *FSDatabase) RegisterNewAdmin(ctx context.Context, newAccount models.AdminAccount) error {
	return nil
}

func (db *FSDatabase) GetAccountApplications(ctx context.Context) ([]models.AccountApplication, error) {
	return []models.AccountApplication{}, nil
}
