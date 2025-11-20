package services

import (
	"api-i18n/main/src/database"
	"api-i18n/main/src/models"
)

// IsLanguageAvailable method to check if a language is available.
func IsLanguageAvailable(languageID string) (bool, error) {
	if result := database.Pg.Limit(1).Find(&models.Language{}, "id = ?", languageID); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}
