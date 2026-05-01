package db

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ksvaza/ees-link/models"
)


func SeedTestData(ctx context.Context) error {
	realDB := &RealDB{}

	if err := seedApplications(ctx, realDB); err != nil {
		return fmt.Errorf("seed applications: %w", err)
	}
	if err := seedAccounts(ctx, realDB); err != nil {
		return fmt.Errorf("seed accounts: %w", err)
	}
	if err := seedAdmins(ctx, realDB); err != nil {
		return fmt.Errorf("seed admins: %w", err)
	}
	if err := seedAccountApplications(ctx, realDB); err != nil {
		return fmt.Errorf("seed account applications: %w", err)
	}

	return nil
}

func seedHash(password, salt string) string {
	h := sha512.New()
	h.Write([]byte(password + salt))
	return hex.EncodeToString(h.Sum(nil))
}

func seedApplications(ctx context.Context, realDB *RealDB) error {
	institution := "RTU"
	howHeard := "internet"

	app := models.RegistrationFormData{
		ID:           "seed-team-001",
		TeamName:     "Seed Team Alpha",
		AgeGroup:     "senior",
		Institution:  &institution,
		CityOrRegion: "Riga",
		Members: []models.TeamMember{
			{
				Key:                    "seed-key-leader-001",
				FullName:               "Janis Seeder",
				DateOfBirth:            "2005-01-10",
				EducationalInstitution: "RTU",
				Role:                   "team_leader",
				ClassOrYear:            "12",
				ID:                     "seed-team-001",
			},
			{
				Key:                    "seed-key-member-001",
				FullName:               "Anna Seeder",
				DateOfBirth:            "2005-04-11",
				EducationalInstitution: "RTU",
				Role:                   "member",
				ClassOrYear:            "12",
				ID:                     "seed-team-001",
			},
		},
		ResponsiblePerson: models.ResponsiblePerson{
			FullName: "Teacher Seeder",
			Phone:    "+37120000001",
			Email:    "teacher.seed@example.com",
			Status:   "teacher",
		},
		HowHeardAbout:   &howHeard,
		Comments:        nil,
		ConfirmTruthful: true,
		ConfirmRules:    true,
		ConfirmMedia:    true,
		AppliedAt:       time.Now().UTC(),
		Status:          "pending",
	}

	existing, err := realDB.GetApplicationByID(ctx, app.ID)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}

	return realDB.RegisterNewApplication(ctx, app)
}

func seedAccounts(ctx context.Context, realDB *RealDB) error {
	const salt = "seed-salt-001"

	accounts := []models.Account{
		{
			Cilveks: models.TeamMember{
				Key:                    "seed-key-leader-001",
				FullName:               "Janis Seeder",
				DateOfBirth:            "2005-01-10",
				EducationalInstitution: "RTU",
				Role:                   "team_leader",
				ClassOrYear:            "12",
				ID:                     "seed-team-001",
			},
			Password:    seedHash("seedPass123", salt),
			Username:    "seed_leader",
			Email:       "janis.seed@example.com",
			PhoneNumber: "+37126000001",
			Salt:        salt,
		},
		{
			Cilveks: models.TeamMember{
				Key:                    "seed-key-member-001",
				FullName:               "Anna Seeder",
				DateOfBirth:            "2005-04-11",
				EducationalInstitution: "RTU",
				Role:                   "member",
				ClassOrYear:            "12",
				ID:                     "seed-team-001",
			},
			Password:    seedHash("seedPass123", salt),
			Username:    "seed_member",
			Email:       "anna.seed@example.com",
			PhoneNumber: "+37126000002",
			Salt:        salt,
		},
	}

	for _, account := range accounts {
		existing, err := realDB.GetAccountByUsername(ctx, account.Username)
		if err != nil {
			return err
		}
		if existing != nil {
			continue
		}

		if err := realDB.RegisterNewAccount(ctx, account); err != nil {
			return err
		}
	}

	return nil
}

func seedAdmins(ctx context.Context, realDB *RealDB) error {
	const username = "seed_admin"
	const salt = "seed-admin-salt"

	existing, err := realDB.GetAdminAccountByUsername(ctx, username)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}

	admin := models.AdminAccount{
		Username:   username,
		Password:   seedHash("seedAdminPass123", salt),
		Salt:       salt,
		Superadmin: false,
	}

	return realDB.RegisterNewAdmin(ctx, admin)
}

func seedAccountApplications(ctx context.Context, realDB *RealDB) error {
	pending := models.AccountApplication{
		FullName:    "Liene Pending",
		DateOfBirth: "2006-01-01",
		Role:        "member",
		Password:    "temp-pass-change-me",
		Username:    "seed_pending_member",
		Email:       "liene.pending@example.com",
		PhoneNumber: "+37127000001",
		TeamName:    "Seed Team Alpha",
	}

	all, err := realDB.GetAccountApplications(ctx)
	if err != nil {
		return err
	}
	for _, existing := range all {
		if existing.Username == pending.Username {
			return nil
		}
	}

	return realDB.RegisterNewAccountApplication(ctx, pending)
}
