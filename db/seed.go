package db

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/ksvaza/ees-link/data"
	"github.com/ksvaza/ees-link/models"
	"github.com/sirupsen/logrus"
)

func SeedTestData(ctx context.Context) error {
	var realDB data.Database
	realDB = &RealDB{}

	if err := seedAccounts(ctx, realDB); err != nil {
		return fmt.Errorf("seed accounts: %w", err)
	}
	if err := seedAdmins(ctx, realDB); err != nil {
		return fmt.Errorf("seed admins: %w", err)
	}
	if err := seedAccountApplications(ctx, realDB); err != nil { // neverificētie konti
		return fmt.Errorf("seed account applications: %w", err)
	}
	if err := seedApplications(ctx, realDB); err != nil {
		return fmt.Errorf("seed team applications: %w", err)
	}
	if err := seedTeams(ctx, realDB); err != nil {
		return fmt.Errorf("seed teams: %w", err)
	}

	return nil
}

func DeleteSeedData(ctx context.Context) error {
	if err := DeleteAllApplications(ctx); err != nil {
		return fmt.Errorf("delete applications: %w", err)
	}
	if err := DeleteAllAccounts(ctx); err != nil {
		return fmt.Errorf("delete accounts: %w", err)
	}
	if err := DeleteAllAdmins(ctx); err != nil {
		return fmt.Errorf("delete admins: %w", err)
	}
	if err := DeleteAllAccountApplications(ctx); err != nil {
		return fmt.Errorf("delete account applications: %w", err)
	}

	return nil
}

func DeleteAllApplications(ctx context.Context) error {
	_, err := Pool.Exec(ctx, "DELETE FROM applicants")
	return err
}

func DeleteAllAccounts(ctx context.Context) error {
	_, err := Pool.Exec(ctx, "DELETE FROM konti")
	return err
}

func DeleteAllAdmins(ctx context.Context) error {
	_, err := Pool.Exec(ctx, "DELETE FROM admini")
	return err
}

func DeleteAllAccountApplications(ctx context.Context) error {
	_, err := Pool.Exec(ctx, "DELETE FROM kontu_pieteikumi")
	return err
}

func seedHash(password, salt string) string {
	h := sha512.New()
	h.Write([]byte(password + salt))
	return hex.EncodeToString(h.Sum(nil))
}

func currentEETTime() time.Time {
	loc, err := time.LoadLocation("Europe/Riga")
	if err != nil {
		loc = time.FixedZone("EET", 2*60*60)
	}
	return time.Now().In(loc)
}

func seedApplications(ctx context.Context, realDB data.Database) error {
	jtv := "JTV"
	rtu := "RTU"
	howHeard := "internet"

	app := []models.RegistrationFormData{
		{
			ID:           "seed-team-002",
			TeamName:     "Seed Team Beta",
			AgeGroup:     "U16",
			Institution:  &jtv,
			CityOrRegion: "Jelgava",
			Members: []models.TeamMember{
				{
					Key:                    "seed-key-leader-002",
					FullName:               "Andris Paraudziņš",
					DateOfBirth:            "2010-01-10",
					EducationalInstitution: "JTV",
					Role:                   "team_leader",
					ClassOrYear:            "9",
					ID:                     "seed-team-002",
					Email:                  "andris.paraudzins@example.com",
				},
				{
					Key:                    "seed-key-member-002",
					FullName:               "Anna Paraudziņa",
					DateOfBirth:            "2011-04-11",
					EducationalInstitution: "JTV",
					Role:                   "member",
					ClassOrYear:            "8",
					ID:                     "seed-team-002",
					Email:                  "anna.paraudzina@example.com",
				},
			},
			ResponsiblePerson: models.ResponsiblePerson{
				FullName: "Skolotājs Paraudziņš",
				Phone:    "+37120000001",
				Email:    "skol.otajs@paraugs.lv",
				Status:   "teacher",
			},
			HowHeardAbout:   &howHeard,
			Comments:        nil,
			ConfirmTruthful: true,
			ConfirmRules:    true,
			ConfirmMedia:    true,
			AppliedAt:       currentEETTime(),
			Status:          "pending",
		},
		{
			ID:           "seed-team-001",
			TeamName:     "Seed Team Alpha",
			AgeGroup:     "U25",
			Institution:  &rtu,
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
					Email:                  "janis.seed@example.com",
				},
				{
					Key:                    "seed-key-member-001",
					FullName:               "Anna Seeder",
					DateOfBirth:            "2005-04-11",
					EducationalInstitution: "RTU",
					Role:                   "member",
					ClassOrYear:            "12",
					ID:                     "seed-team-001",
					Email:                  "anna.seed@example.com",
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
			AppliedAt:       currentEETTime(),
			Status:          "accepted",
		},
	}

	pievC := 0

	for _, a := range app {
		existing, err := realDB.GetApplicationByID(ctx, a.ID)
		if err != nil {
			return err
		}
		if existing != nil {
			continue
		}

		if err := realDB.RegisterNewApplication(ctx, a); err != nil {
			logrus.Info("Te 2")
			return err
		}
		pievC++
	}

	if pievC <= 0 {
		return errors.New("Neizdevās pievienot nevienu komandu pieteikumu")
	}

	return nil
}

func seedTeams(ctx context.Context, realDB data.Database) error {
	teamData := models.TeamData{
		Key:      "seed-team-001",
		CarID:    "seed-car-001",
		TeamName: "Seed Team Alpha",
		Accounts: []models.Account{
			{
				Cilveks: models.TeamMember{
					Key:                    "seed-key-leader-001",
					FullName:               "Janis Seeder",
					DateOfBirth:            "2005-01-10",
					EducationalInstitution: "RTU",
					Role:                   "team_leader",
					ClassOrYear:            "12",
					ID:                     "seed-team-001",
					Email:                  "janis.seed@example.com",
				},
				Password:               seedHash("teamA", "seed-salt-001"),
				Username:               "janis",
				Email:                  "janis.seed@example.com",
				PhoneNumber:            "+37126000001",
				Salt:                   "seed-salt-001",
				PendingTeamID:          "seed-team-001",
				Verified:               true,
				EducationalInstitution: "RTU",
				ClassOrYear:            "12",
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
					Email:                  "anna.seed@example.com",
				},
				Password:               seedHash("teamA", "seed-salt-001"),
				Username:               "annaS",
				Email:                  "anna.seed@example.com",
				PhoneNumber:            "+37126000002",
				Salt:                   "seed-salt-001",
				PendingTeamID:          "seed-team-001",
				Verified:               true,
				EducationalInstitution: "RTU",
				ClassOrYear:            "12",
			},
		},
		AgeGroup:     "U25",
		Institution:  "RTU",
		CityOrRegion: "Riga",
		ResponsiblePerson: models.ResponsiblePerson{
			FullName: "Teacher Seeder",
			Phone:    "+37120000001",
			Email:    "teacher.seed@example.com",
			Status:   "teacher",
		},
		Avatar: "", // base64 iekodēta bilde
		CarData: models.Car{
			SetVoltage: 12,
			MaxCurrent: 100,
			Mass:       250.5,
		},
	}
	existing, err := realDB.GetTeamDataByKey(ctx, teamData.Key)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}
	return realDB.RegisterTeamData(ctx, teamData)
}

func seedAccounts(ctx context.Context, realDB data.Database) error {
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
				Email:                  "janis.seed@example.com",
			},
			Password:               seedHash("teamA", "seed-salt-001"),
			Username:               "janis",
			Email:                  "janis.seed@example.com",
			PhoneNumber:            "+37126000001",
			Salt:                   "seed-salt-001",
			PendingTeamID:          "seed-team-001",
			Verified:               true,
			EducationalInstitution: "RTU",
			ClassOrYear:            "12",
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
				Email:                  "anna.seed@example.com",
			},
			Password:               seedHash("teamA", "seed-salt-001"),
			Username:               "annaS", // vai būs kļūda, ka sakrīt lietotājvārdi?
			Email:                  "anna.seed@example.com",
			PhoneNumber:            "+37126000002",
			Salt:                   "seed-salt-001",
			PendingTeamID:          "seed-team-001",
			Verified:               true,
			EducationalInstitution: "RTU",
			ClassOrYear:            "12",
		},
	}

	pievC := 0

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
		pievC++
	}

	if pievC <= 0 {
		return errors.New("Neizdevās pievienot nevienu kontu")
	}

	return nil
}

func seedAdmins(ctx context.Context, realDB data.Database) error {
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
		Key:        "seed-key-admin-001",
	}

	return realDB.RegisterNewAdmin(ctx, admin)
}

func seedAccountApplications(ctx context.Context, realDB data.Database) error {
	// pending := models.AccountApplication{
	// 	FullName:    "Liene Pending",
	// 	DateOfBirth: "2006-01-01",
	// 	Role:        "member",
	// 	Password:    "temp-pass-change-me",
	// 	Username:    "seed_pending_member",
	// 	Email:       "liene.pending@example.com",
	// 	PhoneNumber: "+37127000001",
	// 	TeamName:    "Seed Team Alpha",
	// }
	// pending.Key = "seed-key-pending-001"

	pending := []models.Account{
		{
			Cilveks: models.TeamMember{
				Key:                    "seed-key-pending-001",
				FullName:               "Liene Pending",
				DateOfBirth:            "2006-01-01",
				EducationalInstitution: "RTU",
				Role:                   "member",
				ClassOrYear:            "12",
				ID:                     "seed-team-001",
				Email:                  "liene.pending@example.com",
			},
			Password:               seedHash("teamA", "seed-pending-salt"),
			Username:               "liene",
			Email:                  "liene.pending@example.com",
			PhoneNumber:            "+37127000001",
			Salt:                   "seed-pending-salt",
			PendingTeamID:          "seed-team-001",
			Verified:               false,
			EducationalInstitution: "RTU",
			ClassOrYear:            "12",
		},
		{
			Cilveks: models.TeamMember{
				Key:                    "seed-key-pending-002",
				FullName:               "Marta Paraudziņa",
				DateOfBirth:            "2011-02-02",
				EducationalInstitution: "",
				Role:                   "member",
				ClassOrYear:            "",
				ID:                     "",
				Email:                  "marta.paraudzina@example.com",
			},
			Password:               seedHash("teamB", "seed-pending-salt"),
			Username:               "marta",
			Email:                  "marta.paraudzina@example.com",
			PhoneNumber:            "+37127000002",
			Salt:                   "seed-pending-salt",
			PendingTeamID:          "seed-team-002",
			Verified:               false,
			EducationalInstitution: "",
			ClassOrYear:            "",
		},
		{
			Cilveks: models.TeamMember{
				Key:                    "seed-key-pending-leader-003",
				FullName:               "Andris Paraudziņš",
				DateOfBirth:            "2010-01-10",
				EducationalInstitution: "JTV",
				Role:                   "team_leader",
				ClassOrYear:            "9",
				ID:                     "",
				Email:                  "andris.paraudzins@example.com",
			},
			Password:               seedHash("teamB", "seed-pending-salt"),
			Username:               "andris",
			Email:                  "andris.paraudzins@example.com",
			PhoneNumber:            "+37127000002",
			Salt:                   "seed-pending-salt",
			PendingTeamID:          "seed-team-002",
			Verified:               false,
			EducationalInstitution: "JTV",
			ClassOrYear:            "9",
		},
		{
			Cilveks: models.TeamMember{
				Key:                    "seed-key-pending-004",
				FullName:               "Anna Paraudziņa",
				DateOfBirth:            "2011-04-11",
				EducationalInstitution: "JTV",
				Role:                   "member",
				ClassOrYear:            "8",
				ID:                     "",
				Email:                  "anna.paraudzina@example.com",
			},
			Password:               seedHash("teamB", "seed-pending-salt"),
			Username:               "annaP", // vai būs kļūda, ka sakrīt lietotājvārdi?
			Email:                  "anna.paraudzina@example.com",
			PhoneNumber:            "+37127000002",
			Salt:                   "seed-pending-salt",
			PendingTeamID:          "seed-team-002",
			Verified:               false,
			EducationalInstitution: "JTV",
			ClassOrYear:            "8",
		},
	}

	// all, err := realDB.GetAccountApplications(ctx)
	// if err != nil {
	// 	return err
	// }
	// for _, existing := range all {
	// 	if existing.Username == pending.Username {
	// 		return nil
	// 	}
	// }

	pievC := 0

	for _, account := range pending {
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
		pievC++
	}

	if pievC <= 0 {
		return errors.New("Neizdevās pievienot nevienu kontu pieteikumu")
	}

	return nil
}
