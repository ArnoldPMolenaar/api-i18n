package services

import (
	"api-i18n/main/src/database"
	"api-i18n/main/src/models"
)

// IsLocaleAvailable method to check if a locale is available.
func IsLocaleAvailable(localeID string) (bool, error) {
	if result := database.Pg.Limit(1).Find(&models.Locale{}, "id = ?", localeID); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}
