package services

import (
	"api-i18n/main/src/database"
	"api-i18n/main/src/models"

	"github.com/samber/lo"
)

// GetTranslationsByLocaleId func to get translations by locale ID.
func GetTranslationsByLocaleId(appName, localeId string) (*map[string]interface{}, error) {
	var keys []models.Key

	tx := database.Pg.Model(&models.Key{}).
		Preload("Category").
		Preload("Translations", "locale_id = ?", localeId).
		Joins("LEFT JOIN categories ON categories.id = category_id").
		Where("app_name = ? AND keys.disabled_at IS NULL AND categories.disabled_at IS NULL", appName).
		Find(&keys)
	if tx.Error != nil {
		return nil, tx.Error
	}

	translations := make(map[string]interface{})
	for _, key := range keys {
		keyName := lo.CamelCase(key.Name)

		if len(key.Translations) == 0 {
			translations[keyName] = nil
			continue
		}

		if !key.CategoryID.Valid {
			translations[key.Name] = key.Translations[0].Value
			continue
		}

		categoryName := lo.CamelCase(key.Category.Name)

		if _, exists := translations[categoryName]; !exists {
			translations[categoryName] = make(map[string]string)
		}

		categoryMap := translations[categoryName].(map[string]string)
		categoryMap[keyName] = key.Translations[0].Value
	}

	return &translations, nil
}
