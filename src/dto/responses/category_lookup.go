package responses

import "api-i18n/main/src/models"

type CategoryLookup struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// SetCategoryLookup sets the category lookup fields from a Category model.
func (cl *CategoryLookup) SetCategoryLookup(c *models.Category) {
	cl.ID = c.ID
	cl.Name = c.Name
}
