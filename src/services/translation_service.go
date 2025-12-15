package services

import (
	"api-i18n/main/src/cache"
	"api-i18n/main/src/database"
	"api-i18n/main/src/models"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/samber/lo"
	"github.com/valkey-io/valkey-go"
)

// GetTranslationsByLocaleId func to get translations by locale ID.
func GetTranslationsByLocaleId(appName, localeID string) (*map[string]interface{}, error) {
	var keys []models.Key

	tx := database.Pg.Model(&models.Key{}).
		Preload("Category").
		Preload("Translations", "locale_id = ?", localeID).
		Joins("LEFT JOIN categories ON categories.id = category_id").
		Where("app_name = ? AND keys.disabled_at IS NULL AND categories.disabled_at IS NULL", appName).
		Find(&keys)
	if tx.Error != nil {
		return nil, tx.Error
	}

	translations := make(map[string]interface{})
	if inCache, err := isTranslationInCache(appName, localeID); err != nil {
		return nil, err
	} else if inCache {
		if cacheTranslation, err := getTranslationFromCache(appName, localeID); err != nil {
			return nil, err
		} else if cacheTranslation != nil && len(*cacheTranslation) > 0 {
			translations = *cacheTranslation
		}
	}

	if len(translations) == 0 {
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

		_ = setTranslationToCache(appName, localeID, &translations)
	}

	return &translations, nil
}

// isTranslationInCache checks if the translation exists in the cache.
func isTranslationInCache(appName, localeID string) (bool, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Exists().Key(translationCacheKey(appName, localeID)).Build())
	if result.Error() != nil {
		return false, result.Error()
	}

	value, err := result.ToInt64()
	if err != nil {
		return false, err
	}

	return value == 1, nil
}

// getTranslationFromCache gets the translation from the cache.
func getTranslationFromCache(appName, localeID string) (*map[string]interface{}, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Get().Key(translationCacheKey(appName, localeID)).Build())
	if result.Error() != nil {
		return nil, result.Error()
	}

	value, err := result.ToString()
	if err != nil {
		return nil, err
	}

	var translation map[string]interface{}
	if err := json.Unmarshal([]byte(value), &translation); err != nil {
		return nil, err
	}

	return &translation, nil
}

// setTranslationToCache sets the translation to the cache.
func setTranslationToCache(appName, localeID string, translations *map[string]interface{}) error {
	value, err := json.Marshal(translations)
	if err != nil {
		return err
	}

	expiration := os.Getenv("VALKEY_EXPIRATION")
	duration, err := time.ParseDuration(expiration)
	if err != nil {
		return err
	}

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Set().Key(translationCacheKey(appName, localeID)).Value(valkey.BinaryString(value)).Ex(duration).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// deleteTranslationFromCache deletes an existing translation from the cache.
func deleteTranslationFromCache(appName, localeID string) error {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Del().Key(translationCacheKey(appName, localeID)).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// deleteAllTranslationsFromCache deletes all translations from the cache.
func deleteAllTranslationsFromCache() error {
	apps, err := GetApps()
	if err != nil {
		return err
	}

	for _, app := range *apps {
		for _, locale := range app.Locales {
			result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Del().Key(translationCacheKey(app.Name, locale.ID)).Build())
			if result.Error() != nil {
				return result.Error()
			}
		}
	}

	return nil
}

// translationCacheKey returns the key for the locales cache.
func translationCacheKey(appName, localeID string) string {
	return fmt.Sprintf("translations:%s:%s", appName, localeID)
}
