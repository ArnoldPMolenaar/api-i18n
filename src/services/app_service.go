package services

import (
	"api-i18n/main/src/database"
	"api-i18n/main/src/models"
	"slices"
)

// IsAppAvailable method to check if an app is available.
func IsAppAvailable(app string) (bool, error) {
	if result := database.Pg.Limit(1).Find(&models.App{}, "name = ?", app); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// GetAppLanguages method to get the languages of an app.
func GetAppLanguages(app string) ([]models.Language, error) {
	a := models.App{}

	if result := database.Pg.Preload("Languages").Find(&a, "name = ?", app); result.Error != nil {
		return nil, result.Error
	}

	return a.Languages, nil
}

// CreateApp method to create an app.
func CreateApp(name string) (*models.App, error) {
	app := &models.App{Name: name}

	if err := database.Pg.FirstOrCreate(&models.App{}, app).Error; err != nil {
		return nil, err
	}

	return app, nil
}

// SetAppLanguages method to set the languages of an app.
// It also restores existing translations for newly added languages
// and deletes translations for removed languages.
func SetAppLanguages(app string, languages []string) error {
	a := models.App{Name: app}
	currentLanguages, err := GetAppLanguages(app)
	if err != nil {
		return err
	}

	currentLanguageIDs := make([]string, len(currentLanguages))
	for i, lang := range currentLanguages {
		currentLanguageIDs[i] = lang.ID
	}

	newLanguages := make([]models.Language, len(languages))
	for i, langID := range languages {
		newLanguages[i] = models.Language{ID: langID}
	}

	// Start a new transaction
	tx := database.Pg.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Clear existing languages
	if err := tx.Model(&a).Association("Languages").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// Set new languages
	for _, language := range languages {
		if err := tx.Model(&a).Association("Languages").Append(&newLanguages); err != nil {
			tx.Rollback()
			return err
		}

		// Check if the language was not already associated
		if len(currentLanguages) == 0 || slices.Contains(currentLanguageIDs, language) {
			continue
		}

		// Try to restore existing translations for the new language
		if result := tx.Model(&models.KeyTranslation{}).
			Unscoped().
			Joins("JOIN keys ON key_translations.key_id = keys.id").
			Where("key_translations.language_id = ?", language).
			Where("keys.app_name = ?", app).
			Update("deleted_at", nil); result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	// Delete translations for removed languages
	for _, languageID := range currentLanguageIDs {
		if slices.Contains(languages, languageID) {
			continue
		}

		if result := tx.Model(&models.KeyTranslation{}).
			Joins("JOIN keys ON key_translations.key_id = keys.id").
			Where("key_translations.language_id = ?", languageID).
			Where("keys.app_name = ?", app).
			Delete(&models.KeyTranslation{}); result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
