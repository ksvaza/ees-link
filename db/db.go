package db

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ksvaza/ees-link/models"
	"github.com/sirupsen/logrus"
)

type RealDB struct{}

//go:embed schema.sql
var schema string
var Pool *pgxpool.Pool

func MigrateUp(ctx context.Context, dbURL string) error {
	var err error
	Pool, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		logrus.WithError(err).Error("Failed to create DB pool")
		return err
	}

	_, err = Pool.Exec(ctx, schema)
	if err != nil {
		fmt.Println("Migration failed:")
		logrus.WithError(err).Error("Failed to run migration")
		return err
	}
	logrus.Info("Migration completed successfully")
	return nil
}

func MigrateDown(ctx context.Context, dbURL string) error {
	var err error
	Pool, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		logrus.WithError(err).Error("Failed to create DB pool")
		return err
	}
	return nil
}

func impregnateDB(ctx context.Context, dbURL string) error {
	var err error
	Pool, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		logrus.WithError(err).Error("Failed to create DB pool")
		return err
	}

	return nil
}

// Health check

func (DB *RealDB) TestHealthiness(ctx context.Context) error {
	err := Pool.Ping(ctx)
	if err != nil {
		logrus.WithError(err).Error("Failed to ping DB")
		return err
	}
	logrus.Info("DB connection is healthy")
	return nil
}

// Applications/teams, registration form data

func (DB *RealDB) RegisterNewApplication(ctx context.Context, newApplicant models.RegistrationFormData) error {
	a, err := DB.GetAllApplications(ctx)
	if err != nil {
		logrus.Warnf("Failed to read from applicants table")
		return err
	}

	for _, applicant := range a {
		if applicant.TeamName == newApplicant.TeamName {
			logrus.Warnf("Applicant with team name %q already exists, skipping insert", newApplicant.TeamName)
			return nil
		}
	}

	membersJSON, err := json.Marshal(newApplicant.Members)
	if err != nil {
		return fmt.Errorf("failed to marshal members: %w", err)
	}

	ResponsiblePersonJSON, e := json.Marshal(newApplicant.ResponsiblePerson)
	if e != nil {
		return fmt.Errorf("failed to marshal responsible person: %w", err)
	}

	_, er := Pool.Exec(ctx, ApplicantWriteRequest,
		newApplicant.ID,              // $1 - id
		newApplicant.TeamName,        // $2 - team_name
		newApplicant.AgeGroup,        // $3 - age_group
		newApplicant.Institution,     // $4 - institution
		newApplicant.CityOrRegion,    // $5 - city_or_region
		membersJSON,                  // $6 - members (assumes proper serialization)
		ResponsiblePersonJSON,        // $7 - responsible_person_full_nam   // $10 - responsible_person_status
		newApplicant.HowHeardAbout,   // $11 - how_heard_about
		newApplicant.Comments,        // $12 - comments
		newApplicant.ConfirmTruthful, // $13 - confirm_truthful
		newApplicant.ConfirmRules,    // $14 - confirm_rules
		newApplicant.ConfirmMedia,    // $15 - confirm_media
		newApplicant.AppliedAt,       // $16 - applied_at
		newApplicant.Status,
	)
	if er != nil {
		logrus.WithError(err).Error("Failed to write to applicants table")
	}

	logrus.Info("Applicant written to DB successfully")

	return err
}

func (DB *RealDB) GetAllApplications(ctx context.Context) ([]models.RegistrationFormData, error) {
	query := ApplicantReadRequest // Adjust the query to select all columns from the applicants table
	rows, err := Pool.Query(ctx, query)
	if err != nil {
		logrus.WithError(err).Error("Failed to query applicants table")
		return nil, err
	}
	defer rows.Close()

	var membersRaw []byte
	var ResponsiblePersonRaw []byte

	var applicants []models.RegistrationFormData
	for rows.Next() {
		var a models.RegistrationFormData
		err := rows.Scan(
			&a.ID,
			&a.TeamName,
			&a.AgeGroup,
			&a.Institution,
			&a.CityOrRegion,
			&membersRaw,
			&ResponsiblePersonRaw,
			&a.HowHeardAbout,
			&a.Comments,
			&a.ConfirmTruthful,
			&a.ConfirmRules,
			&a.ConfirmMedia,
			&a.AppliedAt,
			&a.Status,
		)

		if err != nil {
			logrus.WithError(err).Error("Failed to scan applicant row")
			return nil, err
		}
		err = json.Unmarshal(membersRaw, &a.Members)
		if err != nil {
			logrus.WithError(err).Error("Failed to unmarshal members JSON")
			return nil, err
		}

		err = json.Unmarshal(ResponsiblePersonRaw, &a.ResponsiblePerson)
		if err != nil {
			logrus.WithError(err).Error("Failed to unmarshal responsible person JSON")
			return nil, err
		}

		applicants = append(applicants, a)
	}

	if err = rows.Err(); err != nil {
		logrus.WithError(err).Error("Failed iterating applicant rows")
		return nil, err
	}

	if len(applicants) == 0 {
		return nil, nil
	}

	return applicants, nil
}

func (DB *RealDB) GetApplicationByID(ctx context.Context, id string) (*models.RegistrationFormData, error) {
	a := &models.RegistrationFormData{}
	var membersRaw []byte
	var ResponsiblePersonRaw []byte

	err := Pool.QueryRow(ctx, ApplicantReadRequestByID, id).Scan(
		&a.ID,
		&a.TeamName,
		&a.AgeGroup,
		&a.Institution,
		&a.CityOrRegion,
		&membersRaw,
		&ResponsiblePersonRaw,
		&a.HowHeardAbout,
		&a.Comments,
		&a.ConfirmTruthful,
		&a.ConfirmRules,
		&a.ConfirmMedia,
		&a.AppliedAt,
		&a.Status,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("Failed to scan applicant row")
		return nil, err
	}

	err = json.Unmarshal(ResponsiblePersonRaw, &a.ResponsiblePerson)
	if err != nil {
		logrus.WithError(err).Error("Failed to unmarshal responsible person JSON")
		return nil, err
	}

	err = json.Unmarshal(membersRaw, &a.Members)
	if err != nil {
		logrus.WithError(err).Error("Failed to unmarshal members JSON")
		return nil, err
	}

	return a, nil
}

func (DB *RealDB) UpdateApplication(ctx context.Context, updatedApplicant models.RegistrationFormData) error {
	membersJSON, err := json.Marshal(updatedApplicant.Members)
	if err != nil {
		return fmt.Errorf("failed to marshal members: %w", err)
	}

	responsiblePersonJSON, err := json.Marshal(updatedApplicant.ResponsiblePerson)
	if err != nil {
		return fmt.Errorf("failed to marshal responsible person: %w", err)
	}

	_, err = Pool.Exec(ctx, ApplicantUpdateRequest,
		updatedApplicant.TeamName,
		updatedApplicant.AgeGroup,
		updatedApplicant.Institution,
		updatedApplicant.CityOrRegion,
		membersJSON,
		responsiblePersonJSON,
		updatedApplicant.HowHeardAbout,
		updatedApplicant.Comments,
		updatedApplicant.ConfirmTruthful,
		updatedApplicant.ConfirmRules,
		updatedApplicant.ConfirmMedia,
		updatedApplicant.AppliedAt,
		updatedApplicant.Status,
		updatedApplicant.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update applicant: %w", err)
	}

	return nil
}

func (DB *RealDB) GetApplicationByTeamName(ctx context.Context, teamName string) (*models.RegistrationFormData, error) {
	a := &models.RegistrationFormData{}
	var membersRaw []byte
	var ResponsiblePersonRaw []byte
	err := Pool.QueryRow(ctx, ApplicantReadRequestByTeamName, teamName).Scan(
		&a.ID,
		&a.TeamName,
		&a.AgeGroup,
		&a.Institution,
		&a.CityOrRegion,
		&membersRaw,
		&ResponsiblePersonRaw,
		&a.HowHeardAbout,
		&a.Comments,
		&a.ConfirmTruthful,
		&a.ConfirmRules,
		&a.ConfirmMedia,
		&a.AppliedAt,
		&a.Status,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("Failed to scan applicant row")
		return nil, err
	}
	return a, nil
}

// User accounts

func (DB *RealDB) RegisterNewAccount(ctx context.Context, newAccount models.Account) error {
	_, err := Pool.Exec(ctx, AccountWriteRequest,
		newAccount.Cilveks.Key,
		newAccount.Cilveks.FullName,
		newAccount.Cilveks.DateOfBirth,
		newAccount.Cilveks.Role,
		newAccount.Cilveks.ID,
		newAccount.Password,
		newAccount.Username,
		newAccount.Email,
		newAccount.PhoneNumber,
		newAccount.Salt,
	)

	if err != nil {
		logrus.WithError(err).Error("Failed to write to accounts table")
		return err
	}

	return nil
}

func (DB *RealDB) GetAccountByFullname(ctx context.Context, fullname string) (*models.Account, error) {
	a := &models.Account{}
	err := Pool.QueryRow(ctx, AccountReadRequestByFullname, fullname).Scan(
		&a.Cilveks.Key,
		&a.Cilveks.FullName,
		&a.Cilveks.DateOfBirth,
		&a.Cilveks.Role,
		&a.Cilveks.ID,
		&a.Password,
		&a.Username,
		&a.Email,
		&a.PhoneNumber,
		&a.Salt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("Failed to scan account row")
		return nil, err
	}

	return a, nil
}

func (DB *RealDB) GetAccountByUsername(ctx context.Context, username string) (*models.Account, error) {
	a := &models.Account{}
	err := Pool.QueryRow(ctx, AccountReadRequestByUsername, username).Scan(
		&a.Cilveks.Key,
		&a.Cilveks.FullName,
		&a.Cilveks.DateOfBirth,
		&a.Cilveks.Role,
		&a.Cilveks.ID,
		&a.Password,
		&a.Username,
		&a.Email,
		&a.PhoneNumber,
		&a.Salt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("Failed to scan account row")
		return nil, err
	}

	return a, nil
}

func (DB *RealDB) GetAccountByDateOfBirth(ctx context.Context, dateOfBirth string) (*models.Account, error) {
	a := &models.Account{}
	err := Pool.QueryRow(ctx, AccountReadRequestByDateOfBirth, dateOfBirth).Scan(
		&a.Cilveks.Key,
		&a.Cilveks.FullName,
		&a.Cilveks.DateOfBirth,
		&a.Cilveks.Role,
		&a.Cilveks.ID,
		&a.Password,
		&a.Username,
		&a.Email,
		&a.PhoneNumber,
		&a.Salt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("Failed to scan account row")
		return nil, err
	}

	return a, nil
}

func (DB *RealDB) GetAccountByKey(ctx context.Context, key string) (*models.Account, error) {
	a := &models.Account{}
	err := Pool.QueryRow(ctx, AccountReadRequestByKey, key).Scan(
		&a.Cilveks.Key,
		&a.Cilveks.FullName,
		&a.Cilveks.DateOfBirth,
		&a.Cilveks.Role,
		&a.Cilveks.ID,
		&a.Password,
		&a.Username,
		&a.Email,
		&a.PhoneNumber,
		&a.Salt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("Failed to scan account row")
		return nil, err
	}
	return a, nil
}

func (DB *RealDB) GetAccounts(ctx context.Context) ([]models.Account, error) {
	rows, err := Pool.Query(ctx, AccountReadRequest)
	if err != nil {
		logrus.WithError(err).Error("Failed to query accounts table")
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var a models.Account
		err = rows.Scan(
			&a.Cilveks.Key,
			&a.Cilveks.FullName,
			&a.Cilveks.DateOfBirth,
			&a.Cilveks.Role,
			&a.Cilveks.ID,
			&a.Password,
			&a.Username,
			&a.Email,
			&a.PhoneNumber,
			&a.Salt,
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to scan account row")
			return nil, err
		}
		accounts = append(accounts, a)
	}

	if err = rows.Err(); err != nil {
		logrus.WithError(err).Error("Failed iterating account rows")
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, nil
	}

	return accounts, nil
}

func (DB *RealDB) UpdateAccount(ctx context.Context, updatedAccount models.Account) error {
	_, err := Pool.Exec(ctx, AccountUpdateRequest,
		updatedAccount.Cilveks.FullName,
		updatedAccount.Cilveks.DateOfBirth,
		updatedAccount.Cilveks.Role,
		updatedAccount.Cilveks.ID,
		updatedAccount.Password,
		updatedAccount.Username,
		updatedAccount.Email,
		updatedAccount.PhoneNumber,
		updatedAccount.Salt,
		updatedAccount.Cilveks.Key,
	)

	if err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	return nil
}

// Account applications for user creation

func (DB *RealDB) GetAccountApplications(ctx context.Context) ([]models.AccountApplication, error) {
	rows, err := Pool.Query(ctx, AccountApplicationReadRequest)
	if err != nil {
		logrus.WithError(err).Error("Failed to query account applications table")
		return nil, err
	}
	defer rows.Close()

	var accountApplications []models.AccountApplication
	for rows.Next() {
		var a models.AccountApplication
		err = rows.Scan(
			&a.FullName,
			&a.DateOfBirth,
			&a.Role,
			&a.Password,
			&a.Username,
			&a.Email,
			&a.PhoneNumber,
			&a.TeamName,
			&a.Key,
		)
		if err != nil {
			logrus.WithError(err).Error("Failed to scan account application row")
			return nil, err
		}
		accountApplications = append(accountApplications, a)
	}

	if err = rows.Err(); err != nil {
		logrus.WithError(err).Error("Failed iterating account application rows")
		return nil, err
	}

	if len(accountApplications) == 0 {
		return nil, nil
	}

	return accountApplications, nil
}

func (DB *RealDB) RegisterNewAccountApplication(ctx context.Context, newAccountApplication models.AccountApplication) error {
	_, err := Pool.Exec(ctx, AccountApplicationWriteRequest,
		newAccountApplication.FullName,
		newAccountApplication.DateOfBirth,
		newAccountApplication.Role,
		newAccountApplication.Password,
		newAccountApplication.Username,
		newAccountApplication.Email,
		newAccountApplication.PhoneNumber,
		newAccountApplication.TeamName,
		newAccountApplication.Key,
	)

	if err != nil {
		logrus.WithError(err).Error("Failed to write to accounts table")
		return err
	}

	return nil
}

func (DB *RealDB) GetAccountApplicationByKey(ctx context.Context, key string) (*models.AccountApplication, error) {
	a := &models.AccountApplication{}
	err := Pool.QueryRow(ctx, AccountApplicationReadRequestByKey, key).Scan(
		&a.FullName,
		&a.DateOfBirth,
		&a.Role,
		&a.Password,
		&a.Username,
		&a.Email,
		&a.PhoneNumber,
		&a.TeamName,
		&a.Key,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("Failed to scan account application row")
		return nil, err
	}

	return a, nil
}

func (DB *RealDB) DeleteAccountApplicationByKey(ctx context.Context, key string) error {
	_, err := Pool.Exec(ctx, AccountApplicationDeleteRequestByKey, key)
	if err != nil {
		logrus.WithError(err).Error("Failed to delete account application by key")
		return err
	}
	return nil
}

// Admin accounts

func (DB *RealDB) RegisterNewAdmin(ctx context.Context, newAccount models.AdminAccount) error {
	fmt.Printf("Registering new admin: %+v\n", newAccount)
	_, err := Pool.Exec(ctx, AdminAccountWriteRequest,
		newAccount.Username,
		newAccount.Password,
		newAccount.Salt,
		newAccount.Superadmin,
		newAccount.Key,
	)

	fmt.Printf("ierakstits")

	if err != nil {
		logrus.WithError(err).Error("Failed to write to admin accounts table")
		return err
	}

	return nil
}

func (DB *RealDB) GetAllAdmins(ctx context.Context) ([]models.AdminAccount, error) {
	rows, err := Pool.Query(ctx, AdminAccountReadRequest)
	if err != nil {
		logrus.WithError(err).Error("Failed to query admin accounts table")
		return nil, err
	}

	defer rows.Close()

	var admins []models.AdminAccount
	for rows.Next() {
		var a models.AdminAccount
		err = rows.Scan(
			&a.Username,
			&a.Password,
			&a.Salt,
			&a.Superadmin,
			&a.Key,
		)

		if err != nil {
			logrus.WithError(err).Error("Failed to scan admin account row")
			return nil, err
		}

		admins = append(admins, a)
	}

	if err = rows.Err(); err != nil {
		logrus.WithError(err).Error("Failed iterating admin account rows")
		return nil, err
	}

	if len(admins) == 0 {
		return nil, nil
	}

	return admins, nil
}

func (DB *RealDB) GetAdminAccountByUsername(ctx context.Context, username string) (*models.AdminAccount, error) {
	a := &models.AdminAccount{}
	err := Pool.QueryRow(ctx, AdminAccountReadRequestByUsername, username).Scan(
		&a.Username,
		&a.Password,
		&a.Salt,
		&a.Superadmin,
		&a.Key,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("Failed to scan admin account row")
		return nil, err
	}

	return a, nil
}

func (DB *RealDB) GetAdminAccountByKey(ctx context.Context, key string) (*models.AdminAccount, error) {
	a := &models.AdminAccount{}
	err := Pool.QueryRow(ctx, AdminAccountReadRequestByKey, key).Scan(
		&a.Username,
		&a.Password,
		&a.Salt,
		&a.Superadmin,
		&a.Key,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("Failed to scan admin account row")
		return nil, err
	}

	return a, nil
}

func (DB *RealDB) UpdateAdminAccount(ctx context.Context, updatedAccount models.AdminAccount) error {
	_, err := Pool.Exec(ctx, AdminAccountUpdateRequest,
		updatedAccount.Username,
		updatedAccount.Password,
		updatedAccount.Salt,
		updatedAccount.Superadmin,
		updatedAccount.Key,
	)
	if err != nil {
		return fmt.Errorf("failed to update admin account: %w", err)
	}
	return nil
}
