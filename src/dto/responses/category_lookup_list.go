package responses

import "api-i18n/main/src/models"

type CategoryLookupList struct {
	Categories []CategoryLookup `json:"categories"`
}

// SetCategoryLookupList sets the list of category lookups.
func (cll *CategoryLookupList) SetCategoryLookupList(categories *[]models.Category) {
	cll.Categories = make([]CategoryLookup, len(*categories))
	for i, category := range *categories {
		var cl CategoryLookup
		cl.SetCategoryLookup(&category)
		cll.Categories[i] = cl
	}
}
