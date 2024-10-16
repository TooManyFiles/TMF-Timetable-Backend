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
		On("CONFLICT (userid, settingtype, settingname) DO UPDATE").
		Set("settingsvariable = EXCLUDED.settingsvariable").
		Model(&setting).
		Exec(ctx)
	return err
}

// GetUserSetting retrieves a specific user setting.
func (database *Database) GetUserSetting(userID int, settingType, settingName string, ctx context.Context) (string, error) {
	var setting dbModels.UserSetting
	err := database.DB.NewSelect().
		Model(&setting).
		Where("userid = ? AND settingtype = ? AND settingname = ?", userID, settingType, settingName).
		Scan(ctx)
	return setting.SettingsVariable, err
}
