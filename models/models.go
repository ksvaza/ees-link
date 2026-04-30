package models

import "time"

type TeamMember struct {
	Key                    string `json:"key"`
	FullName               string `json:"fullName"`
	DateOfBirth            string `json:"dateOfBirth"` // Matches "YYYY-MM-DD" format
	EducationalInstitution string `json:"educationalInstitution"`
	Role                   string `json:"role"`
	ClassOrYear            string `json:"classOrYear"`
	ID                     string `json:"id"`
}

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

type RegistrationFormDataRestricted struct {
	TeamName    string    `json:"teamName"`
	Institution *string   `json:"institution,omitempty"`
	MemberCount int       `json:"memberCount"`
	AppliedAt   time.Time `json:"appliedAt"`
	Status      string    `json:"status"`
}

type Konts struct {
	Cilveks TeamMember `json:"dati"`
	Parole         string `json:"parole"`
	Lietotajvards  string `json:"lietotajvards"`
	Epasts         string `json:"epasts"`
	Telefonanumurs string `json:"telefonanumurs"`
}

type KontaPieteikums struct {
	FullName       string `json:"fullName"`
	DateOfBirth    string `json:"dateOfBirth"` // Matches "YYYY-MM-DD" format
	Parole         string `json:"parole"`
	Lietotajvards  string `json:"lietotajvards"`
	Epasts         string `json:"epasts"`
	Telefonanumurs string `json:"telefonanumurs"`
}
