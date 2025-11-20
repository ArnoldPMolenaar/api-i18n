package services

import (
	"api-i18n/main/src/database"
	"api-i18n/main/src/models"
)

// GetCountriesLookup method to get countries lookup by language ID and optional name filter.
func GetCountriesLookup(languageID string, name *string) ([]models.Country, error) {
	var countries []models.Country

	query := database.Pg.Model(&models.Country{}).
		Where("language_id = ?", languageID)

	if name != nil {
		query = query.Where("name ILIKE ?", "%"+*name+"%")
	}

	if result := query.Find(&countries); result.Error != nil {
		return nil, result.Error
	}

	return countries, nil
}
