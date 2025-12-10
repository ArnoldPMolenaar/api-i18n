package responses

import "api-i18n/main/src/models"

type LocaleLookup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SetLocaleLookup sets the locale lookup fields from a LocaleName model.
func (ll *LocaleLookup) SetLocaleLookup(ln *models.LocaleName) {
	ll.ID = ln.LocaleIDTarget
	ll.Name = ln.Name
}
