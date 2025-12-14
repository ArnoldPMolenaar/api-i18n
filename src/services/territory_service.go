package services

import (
	"api-i18n/main/src/cache"
	"api-i18n/main/src/database"
	"api-i18n/main/src/enums"
	"api-i18n/main/src/models"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/valkey-io/valkey-go"
)

// GetTerritoriesLookup method to get territories lookup by locale ID, type and optional name filter.
func GetTerritoriesLookup(localeID string, t *enums.TerritoryType, name *string) (*[]models.TerritoryName, error) {
	territories := make([]models.TerritoryName, 0)

	if inCache, err := isTerritoriesLookupInCache(localeID, t); err != nil {
		return nil, err
	} else if inCache {
		if cacheTerritories, err := getTerritoriesLookupFromCache(localeID, t); err != nil {
			return nil, err
		} else if cacheTerritories != nil && len(*cacheTerritories) > 0 {
			territories = *cacheTerritories
		}
	}

	if len(territories) == 0 {
		query := database.Pg.Model(&models.TerritoryName{}).
			Preload("Territory").
			Joins("JOIN territories ON territory_names.territory_id = territories.id")

		if t != nil {
			query = query.Where("territories.type = ?", *t)
		}

		if result := query.Find(&territories, "locale_id = ?", localeID); result.Error != nil {
			return nil, result.Error
		}

		_ = setTerritoriesLookupToCache(localeID, t, &territories)
	}

	// If a name filter is provided, perform case-insensitive substring match on the list
	if name != nil {
		target := strings.TrimSpace(*name)
		if target != "" {
			lowerTarget := strings.ToLower(target)
			filtered := make([]models.TerritoryName, 0, len(territories))
			for i := range territories {
				if strings.Contains(strings.ToLower(territories[i].Name), lowerTarget) {
					filtered = append(filtered, territories[i])
				}
			}
			territories = filtered
		}
	}

	return &territories, nil
}

// isTerritoriesLookupInCache checks if the territories exists in the cache.
func isTerritoriesLookupInCache(localeID string, t *enums.TerritoryType) (bool, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Exists().Key(territoryLookupCacheKey(localeID, t)).Build())
	if result.Error() != nil {
		return false, result.Error()
	}

	value, err := result.ToInt64()
	if err != nil {
		return false, err
	}

	return value == 1, nil
}

// getTerritoriesLookupFromCache gets the territories from the cache.
func getTerritoriesLookupFromCache(localeID string, t *enums.TerritoryType) (*[]models.TerritoryName, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Get().Key(territoryLookupCacheKey(localeID, t)).Build())
	if result.Error() != nil {
		return nil, result.Error()
	}

	value, err := result.ToString()
	if err != nil {
		return nil, err
	}

	var territories []models.TerritoryName
	if err := json.Unmarshal([]byte(value), &territories); err != nil {
		return nil, err
	}

	return &territories, nil
}

// setTerritoriesLookupToCache sets the territories to the cache.
func setTerritoriesLookupToCache(localeID string, t *enums.TerritoryType, territories *[]models.TerritoryName) error {
	value, err := json.Marshal(territories)
	if err != nil {
		return err
	}

	expiration := os.Getenv("VALKEY_EXPIRATION")
	duration, err := time.ParseDuration(expiration)
	if err != nil {
		return err
	}

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Set().Key(territoryLookupCacheKey(localeID, t)).Value(valkey.BinaryString(value)).Ex(duration).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// territoryLookupCacheKey returns the key for the territories cache.
func territoryLookupCacheKey(localeID string, t *enums.TerritoryType) string {
	tt := "all"
	if t != nil {
		tt = t.String()
	}

	return fmt.Sprintf("territories:lookup:%s:%s", localeID, tt)
}
