package services

import (
	"api-i18n/main/src/database"
	"api-i18n/main/src/dto/requests"
	"api-i18n/main/src/dto/responses"
	"api-i18n/main/src/enums"
	"api-i18n/main/src/models"
	"database/sql"

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

// IsKeyAvailableGlobal method to check if a key name is available globally (regardless of app or category).
func IsKeyAvailableGlobal(keyName string, ignore *string) (bool, error) {
	query := database.Pg.Limit(1)
	var result *gorm.DB

	if ignore != nil {
		result = query.Find(&models.Key{}, "name = ? AND name != ?", keyName, ignore)
	} else {
		result = query.Find(&models.Key{}, "name = ?", keyName)
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

// HasValidTranslations method to check if the provided locale IDs exactly match the app's locales.
func HasValidTranslations(appName string, localeIDs []string) (bool, error) {
	locales, err := GetAppLocales(appName)
	if err != nil {
		return false, err
	}

	// Build a set of IDs from locales.
	localeIDSet := make(map[string]struct{})
	for _, loc := range locales {
		localeIDSet[loc.ID] = struct{}{}
	}

	// Check that every localeID is in locales.
	for _, id := range localeIDs {
		if _, exists := localeIDSet[id]; !exists {
			return false, nil
		}
	}

	// Check that every locale in locales is in localeIDs.
	localeIDsSet := make(map[string]struct{})
	for _, id := range localeIDs {
		localeIDsSet[id] = struct{}{}
	}
	for _, loc := range locales {
		if _, exists := localeIDsSet[loc.ID]; !exists {
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
	dbResult := database.Pg.Scopes(queryFunc, sortFunc, scopeExcludeDeletedCategory).
		Preload("Category").
		Limit(limit).
		Offset(offset)

	total := int64(0)
	dbCount := database.Pg.Scopes(queryFunc, scopeExcludeDeletedCategory).
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

// GetKeyByID method to get a key by ID.
func GetKeyByID(keyID uint) (*models.Key, error) {
	key := &models.Key{}

	if result := database.Pg.
		Scopes(scopeExcludeDeletedCategory).
		Preload("Category").
		Preload("Translations").
		Find(key, "keys.id = ?", keyID); result.Error != nil {
		return nil, result.Error
	}

	return key, nil
}

// CreateKey method to create a key.
func CreateKey(keyDto requests.CreateKey) (*models.Key, error) {
	key := &models.Key{AppName: keyDto.AppName, Name: keyDto.Name}
	if keyDto.CategoryID != nil {
		key.CategoryID = sql.Null[uint]{V: *keyDto.CategoryID, Valid: true}
	}
	if keyDto.DisabledAt != nil {
		key.DisabledAt = sql.NullTime{Time: *keyDto.DisabledAt, Valid: true}
	}
	if keyDto.Description != nil {
		key.Description = sql.NullString{String: *keyDto.Description, Valid: true}
	}

	key.Translations = make([]models.KeyTranslation, len(keyDto.Translations))
	for i, translation := range keyDto.Translations {
		key.Translations[i] = models.KeyTranslation{
			LocaleID:  translation.LocaleID,
			ValueType: enums.ValueType(translation.ValueType),
			Value:     translation.Value,
		}
	}

	if err := database.Pg.Create(&key).Error; err != nil {
		return nil, err
	}

	for i := range key.Translations {
		_ = deleteTranslationFromCache(key.AppName, key.Translations[i].LocaleID)
	}

	return key, nil
}

// UpdateKey method to update a key.
func UpdateKey(oldKey models.Key, keyDto requests.UpdateKey) (*models.Key, error) {
	oldKey.Name = keyDto.Name
	if keyDto.CategoryID != nil {
		oldKey.CategoryID = sql.Null[uint]{V: *keyDto.CategoryID, Valid: true}
	} else {
		oldKey.CategoryID = sql.Null[uint]{Valid: false}
	}
	if keyDto.DisabledAt != nil {
		oldKey.DisabledAt = sql.NullTime{Time: *keyDto.DisabledAt, Valid: true}
	} else {
		oldKey.DisabledAt = sql.NullTime{Valid: false}
	}
	if keyDto.Description != nil {
		oldKey.Description = sql.NullString{String: *keyDto.Description, Valid: true}
	} else {
		oldKey.Description = sql.NullString{Valid: false}
	}

	// Update or add translations
	existingTranslations := make(map[string]*models.KeyTranslation)
	for i := range oldKey.Translations {
		existingTranslations[oldKey.Translations[i].LocaleID] = &oldKey.Translations[i]
	}

	for _, dtoTranslation := range keyDto.Translations {
		if existing, found := existingTranslations[dtoTranslation.LocaleID]; found {
			existing.Value = dtoTranslation.Value
			existing.ValueType = enums.ValueType(dtoTranslation.ValueType)
		} else {
			oldKey.Translations = append(oldKey.Translations, models.KeyTranslation{
				LocaleID:  dtoTranslation.LocaleID,
				ValueType: enums.ValueType(dtoTranslation.ValueType),
				Value:     dtoTranslation.Value,
			})
		}
	}

	if result := database.Pg.Session(&gorm.Session{FullSaveAssociations: true}).Save(&oldKey); result.Error != nil {
		return nil, result.Error
	}

	for i := range oldKey.Translations {
		_ = deleteTranslationFromCache(oldKey.AppName, oldKey.Translations[i].LocaleID)
	}

	return &oldKey, nil
}

// DeleteKey method to delete a key.
func DeleteKey(key *models.Key) error {
	if err := database.Pg.Delete(&models.Key{Model: gorm.Model{ID: key.ID}}).Error; err != nil {
		return err
	}

	for i := range key.Translations {
		_ = deleteTranslationFromCache(key.AppName, key.Translations[i].LocaleID)
	}

	return nil
}

// RestoreKey method to restore a deleted key.
func RestoreKey(keyID uint) error {
	err := database.Pg.Unscoped().Model(&models.Key{}).Where("id = ?", keyID).Update("deleted_at", nil).Error
	if err != nil {
		return err
	}

	if key, err := GetKeyByID(keyID); err != nil {
		return err
	} else {
		for i := range key.Translations {
			_ = deleteTranslationFromCache(key.AppName, key.Translations[i].LocaleID)
		}
	}

	return nil
}

// scopeExcludeDeletedCategory excludes keys whose Category was soft-deleted.
// Keeps keys with NULL category_id.
func scopeExcludeDeletedCategory(db *gorm.DB) *gorm.DB {
	return db.
		Joins("LEFT JOIN categories ON categories.id = keys.category_id").
		Where("(keys.category_id IS NULL OR categories.deleted_at IS NULL)")
}
