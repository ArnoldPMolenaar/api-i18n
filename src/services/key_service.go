package services

import (
	"api-i18n/main/src/database"
	"api-i18n/main/src/dto/responses"
	"api-i18n/main/src/models"

	"github.com/ArnoldPMolenaar/api-utils/pagination"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// IsKeyAvailable method to check if a key is available.
//
// WARNING: If you are using the ignore parameter it is assumed you are updating a key inside the same category.
// When the categoryId is different from old to new, the ignore parameter should also stay null so it checks for
// uniqueness in the new category.
func IsKeyAvailable(appName, keyName string, categoryID *uint, ignore *string) (bool, error) {
	query := database.Pg.Limit(1)
	var result *gorm.DB
	if ignore != nil {
		result = query.Find(&models.Key{}, "app_name = ? AND name = ? AND category_id = ? AND name != ?", appName, keyName, categoryID, ignore)
	} else {
		result = query.Find(&models.Key{}, "app_name = ? AND name = ? AND category_id = ?", appName, keyName, categoryID)
	}

	if result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 0, nil
	}
}

// IsKeyDeleted method to check if a key is deleted.
func IsKeyDeleted(keyID uint) (bool, error) {
	if result := database.Pg.Unscoped().Limit(1).Find(&models.Key{}, "id = ? AND deleted_at IS NOT NULL", keyID); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// HasValidTranslations method to check if the provided language IDs exactly match the app's languages.
func HasValidTranslations(appName string, languageIDs []string) (bool, error) {
	languages, err := GetAppLanguages(appName)
	if err != nil {
		return false, err
	}

	// Build a set of IDs from languages.
	languageIDSet := make(map[string]struct{})
	for _, lang := range languages {
		languageIDSet[lang.ID] = struct{}{}
	}

	// Check that every languageID is in languages.
	for _, id := range languageIDs {
		if _, exists := languageIDSet[id]; !exists {
			return false, nil
		}
	}

	// Check that every language in languages is in languageIDs.
	languageIDsSet := make(map[string]struct{})
	for _, id := range languageIDs {
		languageIDsSet[id] = struct{}{}
	}
	for _, lang := range languages {
		if _, exists := languageIDsSet[lang.ID]; !exists {
			return false, nil
		}
	}

	return true, nil
}

// GetKeys method to get paginated keys.
func GetKeys(c *fiber.Ctx) (*pagination.Model, error) {
	keys := make([]models.Key, 0)
	values := c.Request().URI().QueryArgs()
	allowedColumns := map[string]bool{
		"id":          true,
		"category_id": true,
		"app_name":    true,
		"name":        true,
		"disabled_at": true,
		"created_at":  true,
		"updated_at":  true,
	}

	queryFunc := pagination.Query(values, allowedColumns)
	sortFunc := pagination.Sort(values, allowedColumns)
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	limit := c.QueryInt("limit", 10)
	if limit < 1 {
		limit = 10
	}
	offset := pagination.Offset(page, limit)
	dbResult := database.Pg.Scopes(queryFunc, sortFunc).
		Preload("Category").
		Limit(limit).
		Offset(offset)

	total := int64(0)
	dbCount := database.Pg.Scopes(queryFunc).
		Model(&models.Key{})

	if result := dbResult.Find(&keys); result.Error != nil {
		return nil, result.Error
	}

	dbCount.Count(&total)
	pageCount := pagination.Count(int(total), limit)

	paginatedKeys := make([]responses.PaginatedKey, 0)
	for i := range keys {
		paginatedKey := responses.PaginatedKey{}
		paginatedKey.SetPaginatedKey(&keys[i])
		paginatedKeys = append(paginatedKeys, paginatedKey)
	}

	paginationModel := pagination.CreatePaginationModel(limit, page, pageCount, int(total), paginatedKeys)

	return &paginationModel, nil
}
