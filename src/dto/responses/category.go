package responses

import (
	"api-i18n/main/src/models"
	"time"
)

type Category struct {
	ID         uint       `json:"id"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DisabledAt *time.Time `json:"disabledAt"`
}

// SetCategory sets the category fields from a Category model.
func (c *Category) SetCategory(cat *models.Category) {
	c.ID = cat.ID
	c.Name = cat.Name
	c.CreatedAt = cat.CreatedAt
	c.UpdatedAt = cat.UpdatedAt

	if cat.DisabledAt.Valid {
		c.DisabledAt = &cat.DisabledAt.Time
	}
}
