package fakedb

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ksvaza/ees-link/models"
)

func createApplicationFile(a models.RegistrationFormData) error {

	data, err := json.Marshal(a)
	if err != nil {
		return err
	}

	return os.WriteFile(fmt.Sprintf("fakedbdata/applications/a_%s_%s.json", a.TeamName, a.ID), data, 0644)
}
