package data

import (
	"context"

	"github.com/ksvaza/ees-link/models"
)

type Database interface {
	// Health check
	TestHealthiness(ctx context.Context) error

	// Applications, registration form data
	RegisterNewApplication(ctx context.Context, newApplicant models.RegistrationFormData) error
	GetAllApplications(ctx context.Context) ([]models.RegistrationFormData, error)
	GetApplicationByID(ctx context.Context, id string) (*models.RegistrationFormData, error)
	UpdateApplication(ctx context.Context, updatedApplicant models.RegistrationFormData) error
	GetApplicationByTeamName(ctx context.Context, teamName string) (*models.RegistrationFormData, error)

	// User accounts
	RegisterNewAccount(ctx context.Context, newAccount models.Account) error
	GetAccountByFullname(ctx context.Context, fullname string) (*models.Account, error)
	GetAccountByUsername(ctx context.Context, username string) (*models.Account, error)
	GetAccountByDateOfBirth(ctx context.Context, dateOfBirth string) (*models.Account, error)
	GetAccountByKey(ctx context.Context, key string) (*models.Account, error)
	GetAccounts(ctx context.Context) ([]models.Account, error)
	UpdateAccount(ctx context.Context, updatedAccount models.Account) error

	// Accounts applications for user creation
	GetAccountApplications(ctx context.Context) ([]models.AccountApplication, error)
	RegisterNewAccountApplication(ctx context.Context, newAccountApplication models.AccountApplication) error
	GetAccountApplicationByKey(ctx context.Context, key string) (*models.AccountApplication, error)
	DeleteAccountApplicationByKey(ctx context.Context, key string) error

	// Team data
	GetTeamDataByKey(ctx context.Context, key string) (*models.TeamData, error)
	UpdateTeamDataByKey(ctx context.Context, key string, updatedData models.TeamData) error
	RegisterTeamData(ctx context.Context, newTeamData models.TeamData) error
	AssignAccountsToTeamDataByUsername(ctx context.Context, teamKey string, accountUsernames []string) error

	// Admin accounts
	RegisterNewAdmin(ctx context.Context, newAccount models.AdminAccount) error
	GetAllAdmins(ctx context.Context) ([]models.AdminAccount, error)
	GetAdminAccountByUsername(ctx context.Context, username string) (*models.AdminAccount, error)
	GetAdminAccountByKey(ctx context.Context, key string) (*models.AdminAccount, error)
	UpdateAdminAccount(ctx context.Context, updatedAccount models.AdminAccount) error

	// MQTT logging
	SaveMQTTLog(ctx context.Context, log *models.MqttLogEntry) error
	GetMQTTLogs(ctx context.Context, limit int) ([]models.MqttLogEntry, error)

	// Car telemetry
	SaveCarTelemetry(ctx context.Context, telemetry models.CarTelemetry) error

	// Live race data
	GetLiveRaceData(ctx context.Context) ([]models.LiveRaceData, error)

	// Admin settings
	GetAdminSettings(ctx context.Context) (models.AdminSettings, error)
	UpdateAdminSettings(ctx context.Context, newSettings models.AdminSettings) error

	// Car parameters
	GetAllCarParameters(ctx context.Context) ([]models.CarParameters, error)
	ReplaceAllCarParameters(ctx context.Context, carParams []models.CarParameters) error
}
