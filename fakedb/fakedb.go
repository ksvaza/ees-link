package fakedb

import (
	"encoding/json"
	"os"

	"github.com/ksvaza/ees-link/models"
)

// Read write to file in json
// --------------------------

var path string = "fakedbdata/data.json"

func testExistance() error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return err
	}
	return nil
}

// Read and return all applications (registration forms) from the JSON file
func getAllApplications() ([]models.RegistrationFormData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []models.RegistrationFormData{}, nil
	}

	var apps []models.RegistrationFormData
	if err := json.Unmarshal(data, &apps); err != nil {
		return nil, err
	}

	return apps, nil
}

// Read and return a single application by its ID
func getApplicationByID(id string) (*models.RegistrationFormData, error) {
	apps, err := getAllApplications()
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
func addApplication(app models.RegistrationFormData) error {
	apps, err := getAllApplications()
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
