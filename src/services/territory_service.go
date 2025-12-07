package services

import (
	"api-i18n/main/src/database"
	"api-i18n/main/src/enums"
	"api-i18n/main/src/models"
)

// GetTerritoriesLookup method to get territories lookup by locale ID, type and optional name filter.
func GetTerritoriesLookup(localeID string, t *enums.TerritoryType, name *string) ([]models.TerritoryName, error) {
	territories := make([]models.TerritoryName, 0)

	query := database.Pg.Model(&models.TerritoryName{}).
		Joins("JOIN territories ON territory_names.territory_id = territories.id")

	if t != nil {
		query = query.Where("territories.type = ?", *t)
	}

	if name != nil {
		query = query.Where("name ILIKE ?", "%"+*name+"%")
	}

	if result := query.Find(&territories, "locale_id = ?", localeID); result.Error != nil {
		return nil, result.Error
	}

	return territories, nil
}
