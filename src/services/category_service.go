package services

import (
	"api-i18n/main/src/cache"
	"api-i18n/main/src/database"
	"api-i18n/main/src/dto/responses"
	"api-i18n/main/src/models"
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/pagination"
	"github.com/gofiber/fiber/v2"
	"github.com/valkey-io/valkey-go"
	"gorm.io/gorm"
)

// categoriesLookupCacheKey returns the key for the categories cache.
const categoriesLookupCacheKey = "categories:lookup"

// IsCategoryAvailable method to check if a category is available.
func IsCategoryAvailable(categoryName string, ignore *string) (bool, error) {
	query := database.Pg.Limit(1)
	var result *gorm.DB
	if ignore != nil {
		result = query.Find(&models.Category{}, "name = ? AND name != ?", categoryName, ignore)
	} else {
		result = query.Find(&models.Category{}, "name = ?", categoryName)
	}

	if result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 0, nil
	}
}

// IsCategoryDeleted method to check if a category is deleted.
func IsCategoryDeleted(categoryID uint) (bool, error) {
	if result := database.Pg.Unscoped().Limit(1).Find(&models.Category{}, "id = ? AND deleted_at IS NOT NULL", categoryID); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// GetCategories method to get paginated categories.
func GetCategories(c *fiber.Ctx) (*pagination.Model, error) {
	categories := make([]models.Category, 0)
	values := c.Request().URI().QueryArgs()
	allowedColumns := map[string]bool{
		"id":          true,
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
		Limit(limit).
		Offset(offset)

	total := int64(0)
	dbCount := database.Pg.Scopes(queryFunc).
		Model(&models.Category{})

	if result := dbResult.Find(&categories); result.Error != nil {
		return nil, result.Error
	}

	dbCount.Count(&total)
	pageCount := pagination.Count(int(total), limit)

	paginatedCategories := make([]responses.PaginatedCategory, 0)
	for i := range categories {
		paginatedCategory := responses.PaginatedCategory{}
		paginatedCategory.SetPaginatedCategory(&categories[i])
		paginatedCategories = append(paginatedCategories, paginatedCategory)
	}

	paginationModel := pagination.CreatePaginationModel(limit, page, pageCount, int(total), paginatedCategories)

	return &paginationModel, nil
}

// GetCategoryLookup method to get a lookup of categories.
func GetCategoryLookup(name *string) (*[]models.Category, error) {
	categories := make([]models.Category, 0)

	if inCache, err := isCategoriesLookupInCache(); err != nil {
		return nil, err
	} else if inCache {
		if cacheCategories, err := getCategoriesLookupFromCache(); err != nil {
			return nil, err
		} else if cacheCategories != nil && len(*cacheCategories) > 0 {
			categories = *cacheCategories
		}
	}

	if len(categories) == 0 {
		query := database.Pg.Model(&models.Category{}).
			Select("id", "name")

		if name != nil {
			query = query.Where("name ILIKE ?", "%"+*name+"%")
		}

		if result := query.Find(&categories, "disabled_at IS NULL"); result.Error != nil {
			return nil, result.Error
		}

		_ = setCategoriesLookupToCache(&categories)
	}

	return &categories, nil
}

// GetCategoryByID method to get a category by ID.
func GetCategoryByID(categoryID uint) (*models.Category, error) {
	category := &models.Category{}

	if result := database.Pg.Find(category, "id = ?", categoryID); result.Error != nil {
		return nil, result.Error
	}

	return category, nil
}

// CreateCategory method to create a category.
func CreateCategory(name string, disabledAt *time.Time) (*models.Category, error) {
	category := &models.Category{Name: name}
	if disabledAt != nil {
		category.DisabledAt = sql.NullTime{Time: *disabledAt, Valid: true}
	}

	result := &models.Category{}
	if err := database.Pg.FirstOrCreate(&result, category).Error; err != nil {
		return nil, err
	}

	_ = deleteCategoriesLookupFromCache()

	return result, nil
}

// UpdateCategory method to update a category.
func UpdateCategory(oldCategory models.Category, name string, disabledAt *time.Time) (*models.Category, error) {
	oldCategory.Name = name
	if disabledAt != nil {
		oldCategory.DisabledAt = sql.NullTime{Time: *disabledAt, Valid: true}
	} else {
		oldCategory.DisabledAt = sql.NullTime{Valid: false}
	}

	if result := database.Pg.Save(&oldCategory); result.Error != nil {
		return nil, result.Error
	}

	_ = deleteCategoriesLookupFromCache()

	return &oldCategory, nil
}

// DeleteCategory method to delete a category.
func DeleteCategory(categoryID uint) error {
	err := database.Pg.Delete(&models.Category{Model: gorm.Model{ID: categoryID}}).Error
	if err == nil {
		_ = deleteCategoriesLookupFromCache()
	}

	return err
}

// RestoreCategory method to restore a deleted category.
func RestoreCategory(categoryID uint) error {
	err := database.Pg.Unscoped().Model(&models.Category{}).Where("id = ?", categoryID).Update("deleted_at", nil).Error
	if err == nil {
		_ = deleteCategoriesLookupFromCache()
	}

	return err
}

// isCategoriesLookupInCache checks if the categories exists in the cache.
func isCategoriesLookupInCache() (bool, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Exists().Key(categoriesLookupCacheKey).Build())
	if result.Error() != nil {
		return false, result.Error()
	}

	value, err := result.ToInt64()
	if err != nil {
		return false, err
	}

	return value == 1, nil
}

// getCategoriesLookupFromCache gets the categories from the cache.
func getCategoriesLookupFromCache() (*[]models.Category, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Get().Key(categoriesLookupCacheKey).Build())
	if result.Error() != nil {
		return nil, result.Error()
	}

	value, err := result.ToString()
	if err != nil {
		return nil, err
	}

	var categories []models.Category
	if err := json.Unmarshal([]byte(value), &categories); err != nil {
		return nil, err
	}

	return &categories, nil
}

// setCategoriesLookupToCache sets the categories to the cache.
func setCategoriesLookupToCache(categories *[]models.Category) error {
	value, err := json.Marshal(categories)
	if err != nil {
		return err
	}

	expiration := os.Getenv("VALKEY_EXPIRATION")
	duration, err := time.ParseDuration(expiration)
	if err != nil {
		return err
	}

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Set().Key(categoriesLookupCacheKey).Value(valkey.BinaryString(value)).Ex(duration).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// deleteCategoriesLookupFromCache deletes an existing categories from the cache.
func deleteCategoriesLookupFromCache() error {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Del().Key(categoriesLookupCacheKey).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}
