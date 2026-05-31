package db

import (
	"context"

	"github.com/ksvaza/ees-link/models"
)

var adminSettingsLocal models.AdminSettings

func (DB *RealDB) GetAdminSettings(ctx context.Context) (models.AdminSettings, error) {
	return adminSettingsLocal, nil
}

func (DB *RealDB) UpdateAdminSettings(ctx context.Context, newSettings models.AdminSettings) error {
	adminSettingsLocal = newSettings
	return nil
}
