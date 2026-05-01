package data

import (
	"context"

	"github.com/ksvaza/ees-link/models"
)

type Database interface {
	TestHealthiness(ctx context.Context) error

	GetAllApplications(ctx context.Context) ([]models.RegistrationFormData, error)
	GetApplicationByID(ctx context.Context, id string) (models.RegistrationFormData, error)

	RegisterNewApplication(ctx context.Context, newApplicant models.RegistrationFormData) error
	UpdateApplication(ctx context.Context, updatedApplicant models.RegistrationFormData) error

	RegisterNewAccount(ctx context.Context, newAccount models.Account) error

	RegisterNewAccountApplication(ctx context.Context, newAccountApplication models.AccountApplication) (models.AccountApplication, error)

	GetAccountByUsername(ctx context.Context, username string) (models.Account, error)
	GetAccountByDateOfBirth(ctx context.Context, dateOfBirth string) (models.Account, error)
	GetAccountByFullname(ctx context.Context, fullname string) (models.Account, error)

	RegisterNewAdmin(ctx context.Context, newAccount models.AdminAccount) error
	GetAccountApplications(ctx context.Context) ([]models.AccountApplication, error)
}
