package responses

import "api-i18n/main/src/models"

type TerritoryLookup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// SetTerritoryLookup sets the territory lookup fields from a TerritoryName model.
func (tl *TerritoryLookup) SetTerritoryLookup(tn *models.TerritoryName) {
	tl.ID = tn.TerritoryID
	tl.Name = tn.Name
	tl.Type = tn.Territory.Type.String()
}
