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

// GetLocaleLookup method to get locale lookup names.
func GetLocaleLookup(localeId string, name *string) (*[]models.LocaleName, error) {
	locales := make([]models.LocaleName, 0)

	query := database.Pg.Model(&models.LocaleName{}).
		Select("locale_id_target", "name")

	if name != nil {
		query = query.Where("name ILIKE ?", "%"+*name+"%")
	}

	if result := query.Find(&locales, "locale_id_viewer = ?", localeId); result.Error != nil {
		return nil, result.Error
	}

	return &locales, nil
}
