package responses

import "api-i18n/main/src/models"

type LocaleLookupList struct {
	Locales []LocaleLookup `json:"locales"`
}

// SetLocaleLookupList sets the list of locale lookups.
func (lll *LocaleLookupList) SetLocaleLookupList(locales *[]models.LocaleName) {
	lll.Locales = make([]LocaleLookup, len(*locales))
	for i, locale := range *locales {
		var ll LocaleLookup
		ll.SetLocaleLookup(&locale)
		lll.Locales[i] = ll
	}
}
