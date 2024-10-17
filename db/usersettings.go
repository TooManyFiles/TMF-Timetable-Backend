package db

import (
	"context"

	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
)

// UpdateUserSetting stores or updates a user setting.
func (database *Database) UpdateUserSetting(userID int, settingType, settingName, settingValue string, ctx context.Context) error {
	setting := dbModels.UserSetting{
		UserID:           userID,
		SettingType:      settingType,
		SettingName:      settingName,
		SettingsVariable: settingValue,
	}
	_, err := database.DB.NewInsert().
		Model(&setting).
		On("CONFLICT (userid, setting_type, setting_name) DO UPDATE").
		Exec(ctx)
	return err
}

// GetUserSetting retrieves a specific user setting.
func (database *Database) GetUserSetting(userID int, settingType, settingName string, ctx context.Context) (string, error) {
	setting := dbModels.UserSetting{
		UserID:      userID,
		SettingType: settingType,
		SettingName: settingName,
	}
	err := database.DB.NewSelect().
		Model(&setting).
		WherePK().
		Scan(ctx)
	return setting.SettingsVariable, err
}
