package responses

import (
	"api-i18n/main/src/models"
	"time"
)

type PaginatedCategory struct {
	ID         uint       `json:"id"`
	Name       string     `json:"name"`
	DisabledAt *time.Time `json:"disabledAt"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

// SetPaginatedCategory method to set category data from models.Category{}.
func (c *PaginatedCategory) SetPaginatedCategory(category *models.Category) {
	c.ID = category.ID
	c.Name = category.Name
	c.CreatedAt = category.CreatedAt
	c.UpdatedAt = category.UpdatedAt
	c.DisabledAt = func() *time.Time {
		if category.DisabledAt.Valid {
			return &category.DisabledAt.Time
		}
		return nil
	}()
}
