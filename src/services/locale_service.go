package services

import (
	"api-i18n/main/src/cache"
	"api-i18n/main/src/database"
	"api-i18n/main/src/models"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/valkey-io/valkey-go"
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

	if inCache, err := isLocalesLookupInCache(localeId); err != nil {
		return nil, err
	} else if inCache {
		if cacheLocales, err := getLocalesLookupFromCache(localeId); err != nil {
			return nil, err
		} else if cacheLocales != nil && len(*cacheLocales) > 0 {
			locales = *cacheLocales
		}
	}

	if len(locales) == 0 {
		query := database.Pg.Model(&models.LocaleName{}).
			Select("locale_id_target", "name")

		if result := query.Find(&locales, "locale_id_viewer = ?", localeId); result.Error != nil {
			return nil, result.Error
		}

		_ = setLocalesLookupToCache(localeId, &locales)
	}

	// If a name filter is provided, perform case-insensitive substring match on the list
	if name != nil {
		target := strings.TrimSpace(*name)
		if target != "" {
			lowerTarget := strings.ToLower(target)
			filtered := make([]models.LocaleName, 0, len(locales))
			for i := range locales {
				if strings.Contains(strings.ToLower(locales[i].Name), lowerTarget) {
					filtered = append(filtered, locales[i])
				}
			}
			locales = filtered
		}
	}

	return &locales, nil
}

// isLocalesLookupInCache checks if the locales exists in the cache.
func isLocalesLookupInCache(localeId string) (bool, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Exists().Key(localeLookupCacheKey(localeId)).Build())
	if result.Error() != nil {
		return false, result.Error()
	}

	value, err := result.ToInt64()
	if err != nil {
		return false, err
	}

	return value == 1, nil
}

// getLocalesLookupFromCache gets the locales from the cache.
func getLocalesLookupFromCache(localeId string) (*[]models.LocaleName, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Get().Key(localeLookupCacheKey(localeId)).Build())
	if result.Error() != nil {
		return nil, result.Error()
	}

	value, err := result.ToString()
	if err != nil {
		return nil, err
	}

	var locales []models.LocaleName
	if err := json.Unmarshal([]byte(value), &locales); err != nil {
		return nil, err
	}

	return &locales, nil
}

// setLocalesLookupToCache sets the locales to the cache.
func setLocalesLookupToCache(localeId string, locales *[]models.LocaleName) error {
	value, err := json.Marshal(locales)
	if err != nil {
		return err
	}

	expiration := os.Getenv("VALKEY_EXPIRATION")
	duration, err := time.ParseDuration(expiration)
	if err != nil {
		return err
	}

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Set().Key(localeLookupCacheKey(localeId)).Value(valkey.BinaryString(value)).Ex(duration).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// localeLookupCacheKey returns the key for the locales cache.
func localeLookupCacheKey(localeId string) string {
	return fmt.Sprintf("locales:lookup:%s", localeId)
}
