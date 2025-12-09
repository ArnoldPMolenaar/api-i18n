package responses

import "api-i18n/main/src/models"

type TerritoryLookupList struct {
	Territories []TerritoryLookup `json:"territories"`
}

// SetTerritoryLookupList sets the list of territory lookups.
func (tll *TerritoryLookupList) SetTerritoryLookupList(territories *[]models.TerritoryName) {
	tll.Territories = make([]TerritoryLookup, len(*territories))
	for i, territory := range *territories {
		var tl TerritoryLookup
		tl.SetTerritoryLookup(&territory)
		tll.Territories[i] = tl
	}
}
