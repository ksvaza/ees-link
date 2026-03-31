package fakedb

import (
	"encoding/json"
	"os"
	"time"
)

// Read write to file in json
// --------------------------

// RegistrationFormData represents the top-level registration structure
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
	AplliedAt         time.Time         `json:"appliedAt"`
	Status            string            `json:"status"`
}

// TeamMember represents individual members in the members array
type TeamMember struct {
	Key                    string `json:"key"`
	FullName               string `json:"fullName"`
	DateOfBirth            string `json:"dateOfBirth"` // Matches "YYYY-MM-DD" format
	EducationalInstitution string `json:"educationalInstitution"`
	Role                   string `json:"role"`
	ClassOrYear            string `json:"classOrYear"`
}

// ResponsiblePerson represents the contact person details
type ResponsiblePerson struct {
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

var path string = "fakedbdata/data.json"

// Read and return all applications (registration forms) from the JSON file
func GetAllApplications() ([]RegistrationFormData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []RegistrationFormData{}, nil
	}

	var apps []RegistrationFormData
	if err := json.Unmarshal(data, &apps); err != nil {
		return nil, err
	}

	return apps, nil
}

// Read and return a single application by its ID
func GetApplicationByID(id string) (*RegistrationFormData, error) {
	apps, err := GetAllApplications()
	if err != nil {
		return nil, err
	}

	for _, app := range apps {
		if app.ID == id {
			return &app, nil
		}
	}

	return nil, nil
}

// Add a new application to the JSON file
func AddApplication(app RegistrationFormData) error {
	apps, err := GetAllApplications()
	if err != nil {
		return err
	}

	// Check if application with same ID already exists
	found := false
	for i, existing := range apps {
		if existing.ID == app.ID {
			apps[i] = app
			found = true
			break
		}
	}

	// If not found, append new application
	if !found {
		apps = append(apps, app)
	}

	// Write back to file
	data, err := json.MarshalIndent(apps, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
