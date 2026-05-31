package models

type TeamMember struct {
	Key                    string `json:"key"`
	FullName               string `json:"fullName"`
	DateOfBirth            string `json:"dateOfBirth"` // Matches "YYYY-MM-DD" format
	EducationalInstitution string `json:"educationalInstitution"`
	Role                   string `json:"role"`
	ClassOrYear            string `json:"classOrYear"`
	ID                     string `json:"id"` // teamID
	Email                  string `json:"email"`
}

type Account struct {
	Cilveks                TeamMember `json:"dati"`
	Password               string     `json:"parole"`
	Username               string     `json:"lietotajvards"`
	Email                  string     `json:"epasts"`
	PhoneNumber            string     `json:"telefonanumurs"`
	Salt                   string     `json:"salt"`
	PendingTeamID          string     `json:"pendingTeamId,omitempty"` // ID of the team the user has applied to join, if any
	Verified               bool       `json:"registered"`
	EducationalInstitution string     `json:"educationalInstitution"`
	ClassOrYear            string     `json:"classOrYear"`
	Avatar                 []byte     `json:"avatar"` // tehniski base64 enkodēts binārs fails - bilde
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
	Salt        string `json:"salt"`
	Username    string `json:"lietotajvards"`
	Email       string `json:"epasts"`
	PhoneNumber string `json:"telefonanumurs"`
	TeamName    string `json:"komandasNosaukums"`
	Role        string `json:"role"`
}

type AccountVerificationCriteria struct {
	CanRegister   bool `json:"canRegister"`
	CanLinkToTeam bool `json:"canLinkToTeam"`

	// General requirements (visible to everyone)
	RequiredFieldsPresent bool `json:"requiredFieldsPresent"`
	UsernameAvailable     bool `json:"usernameAvailable"`

	// Standalone account requirements (visible to everyone)
	NoTeamNameProvided *bool `json:"noTeamNameProvided,omitempty"`
	NotTeamLeaderRole  *bool `json:"notTeamLeaderRole,omitempty"`

	// Team member linking requirements (visible to team leaders and superadmin)
	TeamNameProvided      *bool `json:"teamNameProvided,omitempty"`
	TeamApplicationExists *bool `json:"teamApplicationExists,omitempty"`
	TeamMemberFound       *bool `json:"teamMemberFound,omitempty"`
	TeamMemberRoleMatches *bool `json:"teamMemberRoleMatches,omitempty"`
	NoMatchingTeamMember  *bool `json:"noMatchingTeamMember,omitempty"`

	// Team leader requirements (visible to team leaders and superadmin)
	TeamLeaderRole             *bool `json:"teamLeaderRole,omitempty"`
	TeamLeaderMemberFound      *bool `json:"teamLeaderMemberFound,omitempty"`
	TeamLeaderDateOfBirthMatch *bool `json:"teamLeaderDateOfBirthMatch,omitempty"`
}
