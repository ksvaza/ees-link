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

// Health check
func (db *FSDatabase) TestHealthiness(ctx context.Context) error {
	return testExistance()
}

// Applications, registration form data
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

func (db *FSDatabase) GetAllApplications(ctx context.Context) ([]models.RegistrationFormData, error) {
	return getAllApplications()
}

func (db *FSDatabase) GetApplicationByID(ctx context.Context, id string) (*models.RegistrationFormData, error) {
	app, err := getApplicationByID(id)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, nil
	}
	return app, nil
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

func (db *FSDatabase) GetApplicationByTeamName(ctx context.Context, teamName string) (*models.RegistrationFormData, error) {
	apps, err := getAllApplications()
	if err != nil {
		return nil, err
	}
	for _, app := range apps {
		if app.TeamName == teamName {
			return &app, nil
		}
	}
	return nil, nil
}

// User accounts
func (db *FSDatabase) RegisterNewAccount(ctx context.Context, newAccount models.Account) error {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil
}

func (db *FSDatabase) GetAccountByFullname(ctx context.Context, fullname string) (*models.Account, error) {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil, nil
}

func (db *FSDatabase) GetAccountByUsername(ctx context.Context, username string) (*models.Account, error) {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil, nil
}

func (db *FSDatabase) GetAccountByDateOfBirth(ctx context.Context, dateOfBirth string) (*models.Account, error) {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil, nil
}

func (db *FSDatabase) GetAccountByKey(ctx context.Context, key string) (*models.Account, error) {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil, nil
}

func (db *FSDatabase) GetAccounts(ctx context.Context) ([]models.Account, error) {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return []models.Account{}, nil
}

func (db *FSDatabase) UpdateAccount(ctx context.Context, updatedAccount models.Account) error {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil
}

// Account applications for user creation
func (db *FSDatabase) GetAccountApplications(ctx context.Context) ([]models.AccountApplication, error) {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return []models.AccountApplication{}, nil
}

func (db *FSDatabase) RegisterNewAccountApplication(ctx context.Context, newAccountApplication models.AccountApplication) error {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil
}

func (db *FSDatabase) GetAccountApplicationByKey(ctx context.Context, key string) (*models.AccountApplication, error) {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil, nil
}

func (db *FSDatabase) DeleteAccountApplicationByKey(ctx context.Context, key string) error {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil
}

// Admin accounts
func (db *FSDatabase) RegisterNewAdmin(ctx context.Context, newAccount models.AdminAccount) error {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil
}

func (db *FSDatabase) GetAllAdmins(ctx context.Context) ([]models.AdminAccount, error) {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return []models.AdminAccount{}, nil
}

func (db *FSDatabase) GetAdminAccountByUsername(ctx context.Context, username string) (*models.AdminAccount, error) {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil, nil
}

func (db *FSDatabase) GetAdminAccountByKey(ctx context.Context, key string) (*models.AdminAccount, error) {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil, nil
}

func (db *FSDatabase) UpdateAdminAccount(ctx context.Context, updatedAccount models.AdminAccount) error {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil
}

func (db *FSDatabase) UpdateTeamDataByKey(ctx context.Context, key string, newTeamData models.TeamData) error {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil
}

func (db *FSDatabase) RegisterTeamData(ctx context.Context, newTeamData models.TeamData) error {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil
}

func (db *FSDatabase) GetTeamDataByKey(ctx context.Context, key string) (*models.TeamData, error) {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil, nil
}

func (db *FSDatabase) AssignAccountsToTeamDataByUsername(ctx context.Context, teamKey string, accountUsernames []string) error {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil
}

func (db *FSDatabase) SaveMQTTLog(ctx context.Context, log *models.MqttLogEntry) error {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return nil
}

func (db *FSDatabase) GetMQTTLogs(ctx context.Context, limit int) ([]models.MqttLogEntry, error) {
	logrus.Info("fake db called")
	fmt.Printf("fake db called")
	return []models.MqttLogEntry{}, nil
}
