package models

import "time"

// ResponsiblePerson represents the contact person details
type ResponsiblePerson struct {
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

type RegistrationFormData struct {
	ID                string            `json:"id"`
	TeamName          string            `json:"teamName"`
	AgeGroup          string            `json:"ageGroup"` // Assumes AgeGroup is a string enum
	Institution       *string           `json:"institution,omitempty"`
	CityOrRegion      string            `json:"cityOrRegion"`
	Members           []TeamMember      `json:"members"`
	ResponsiblePerson ResponsiblePerson `json:"responsiblePerson"`
	HowHeardAbout     *string           `json:"howHeardAbout,omitempty"`
	Comments          *string           `json:"comments,omitempty"`
	ConfirmTruthful   bool              `json:"confirmTruthful"`
	ConfirmRules      bool              `json:"confirmRules"`
	ConfirmMedia      bool              `json:"confirmMedia"`
	AppliedAt         time.Time         `json:"appliedAt"`
	Status            string            `json:"status"`
}

func (r RegistrationFormData) FindTeamMemberByFullName(fullName string) *TeamMember {
	for _, member := range r.Members {
		if member.FullName == fullName {
			return &member
		}
	}
	return nil
}

func (r RegistrationFormData) FindTeamLeader() *TeamMember {
	for _, member := range r.Members {
		if member.Role == "team_leader" {
			return &member
		}
	}
	return nil
}

type RegistrationFormDataRestricted struct {
	TeamName    string    `json:"teamName"`
	Institution *string   `json:"institution,omitempty"`
	MemberCount int       `json:"memberCount"`
	AppliedAt   time.Time `json:"appliedAt"`
	Status      string    `json:"status"`
}

type TeamData struct {
	Key               string            `json:"key"`
	CarID             string            `json:"carId"`
	TeamName          string            `json:"teamName"`
	Accounts          []Account         `json:"members"`
	AgeGroup          string            `json:"ageGroup"`
	Institution       string            `json:"institution,omitempty"`
	CityOrRegion      string            `json:"cityOrRegion"`
	ResponsiblePerson ResponsiblePerson `json:"responsiblePerson"`
	Avatar            string            `json:"avatar"`
	CarData           Car               `json:"carData,omitempty"`
}

type Car struct {
	SetVoltage int     `json:"setVoltage"`
	MaxCurrent int     `json:"maxCurrent"`
	Mass       float32 `json:"mass"`
}
