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
}
