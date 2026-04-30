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

<<<<<<< HEAD
type Account struct {
	Cilveks TeamMember `json:"dati"`
	Password         string `json:"parole"`
	Username  string `json:"lietotajvards"`
	Email         string `json:"epasts"`
	PhoneNumber string `json:"telefonanumurs"`
=======
type Konts struct {
	Cilveks        TeamMember `json:"dati"`
	Parole         string     `json:"parole"`
	Lietotajvards  string     `json:"lietotajvards"`
	Epasts         string     `json:"epasts"`
	Telefonanumurs string     `json:"telefonanumurs"`
>>>>>>> ed7934f4400024e9fab42d6d268359390e88e893
}

type AccountApplication struct {
	FullName       string `json:"fullName"`
<<<<<<< HEAD
	DateOfBirth    string `json:"dateOfBirth"` // Matches "YYYY-MM-DD" format
	Password         string `json:"parole"`
	Username  string `json:"lietotajvards"`
	Email         string `json:"epasts"`
	PhoneNumber string `json:"telefonanumurs"`
	TeamName string `json:"komandasNosaukums"`
=======
	DateOfBirth    string `json:"dateOfBirth"` // Sakrīt ar formātu "GGGG-MM-DD".
	Parole         string `json:"parole"`
	Lietotajvards  string `json:"lietotajvards"`
	Epasts         string `json:"epasts"`
	Telefonanumurs string `json:"telefonanumurs"`
>>>>>>> ed7934f4400024e9fab42d6d268359390e88e893
}
