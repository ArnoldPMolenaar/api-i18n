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

// GetAppLocales method to get the locales of an app.
func GetAppLocales(app string) ([]models.Locale, error) {
	a := models.App{}

	if result := database.Pg.Preload("Locales").Find(&a, "name = ?", app); result.Error != nil {
		return nil, result.Error
	}

	return a.Locales, nil
}

// CreateApp method to create an app.
func CreateApp(name string) (*models.App, error) {
	app := &models.App{Name: name}

	if err := database.Pg.FirstOrCreate(&models.App{}, app).Error; err != nil {
		return nil, err
	}

	return app, nil
}

// SetAppLocales method to set the locales of an app.
// It also restores existing translations for newly added locales
// and deletes translations for removed locales.
func SetAppLocales(app string, locales []string) error {
	a := models.App{Name: app}
	currentLocales, err := GetAppLocales(app)
	if err != nil {
		return err
	}

	currentLocaleIDs := make([]string, len(currentLocales))
	for i, lang := range currentLocales {
		currentLocaleIDs[i] = lang.ID
	}

	newLocales := make([]models.Locale, len(locales))
	for i, langID := range locales {
		newLocales[i] = models.Locale{ID: langID}
	}

	// Start a new transaction
	tx := database.Pg.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Clear existing locales
	if err := tx.Model(&a).Association("Locales").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// Set new locales
	for _, locale := range locales {
		if err := tx.Model(&a).Association("Locales").Append(&newLocales); err != nil {
			tx.Rollback()
			return err
		}

		// Check if the locale was not already associated
		if len(currentLocales) == 0 || slices.Contains(currentLocaleIDs, locale) {
			continue
		}

		// Try to restore existing translations for the new locale
		if result := tx.Model(&models.KeyTranslation{}).
			Unscoped().
			Joins("JOIN keys ON key_translations.key_id = keys.id").
			Where("key_translations.locale_id = ?", locale).
			Where("keys.app_name = ?", app).
			Update("deleted_at", nil); result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	// Delete translations for removed locales
	for _, localeID := range currentLocaleIDs {
		if slices.Contains(locales, localeID) {
			continue
		}

		if result := tx.Model(&models.KeyTranslation{}).
			Joins("JOIN keys ON key_translations.key_id = keys.id").
			Where("key_translations.locale_id = ?", localeID).
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
