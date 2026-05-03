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

type Account struct {
	Cilveks     TeamMember `json:"dati"`
	Password    string     `json:"parole"`
	Username    string     `json:"lietotajvards"`
	Email       string     `json:"epasts"`
	PhoneNumber string     `json:"telefonanumurs"`
	Salt        string     `json:"salt"`
	// varbūt kaut kas trūkst tīri moderēšans pēc
}

type AdminAccount struct {
	Key        string `json:"key"`
	Username   string `json:"lietotajvards"`
	Password   string `json:"parole"`
	Salt       string `json:"salt"`
	Superadmin bool   `json:"superadmin"`
}

type AccountApplication struct {
	Key         string `json:"key"`
	FullName    string `json:"fullName"`
	DateOfBirth string `json:"dateOfBirth"` // Matches "YYYY-MM-DD" format
	Password    string `json:"parole"`
	Username    string `json:"lietotajvards"`
	Email       string `json:"epasts"`
	PhoneNumber string `json:"telefonanumurs"`
	TeamName    string `json:"komandasNosaukums"`
	Role        string `json:"role"`
}

type AccountVerificationCriteria struct {
	// General requirements (visible to everyone)
	RequiredFieldsPresent bool `json:"requiredFieldsPresent"`
	UsernameAvailable     bool `json:"usernameAvailable"`

	// Standalone account requirements (visible to everyone)
	NoTeamNameProvided *bool `json:"noTeamNameProvided,omitempty"`
	NotTeamLeaderRole  *bool `json:"notTeamLeaderRole,omitempty"`

	// Team member linking requirements (visible to team leaders and superadmin)
	TeamNameProvided         *bool `json:"teamNameProvided,omitempty"`
	TeamApplicationExists    *bool `json:"teamApplicationExists,omitempty"`
	TeamMemberFound          *bool `json:"teamMemberFound,omitempty"`
	TeamMemberRoleMatches    *bool `json:"teamMemberRoleMatches,omitempty"`
	NoMatchingTeamMember     *bool `json:"noMatchingTeamMember,omitempty"`

	// Team leader requirements (visible to team leaders and superadmin)
	TeamLeaderRole             *bool `json:"teamLeaderRole,omitempty"`
	TeamLeaderMemberFound      *bool `json:"teamLeaderMemberFound,omitempty"`
	TeamLeaderDateOfBirthMatch *bool `json:"teamLeaderDateOfBirthMatch,omitempty"`
}
