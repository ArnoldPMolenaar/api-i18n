package responses

import (
	"api-i18n/main/src/models"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/utils"
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
	c.DisabledAt = utils.PtrFromNullTime(cat.DisabledAt)
	c.CreatedAt = cat.CreatedAt
	c.UpdatedAt = cat.UpdatedAt
}
