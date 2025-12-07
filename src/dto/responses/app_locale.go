package responses

import "api-i18n/main/src/models"

// AppLocale represents the response structure for application locale settings.
type AppLocale struct {
	AppName string   `json:"app_name"`
	Locales []string `json:"locales"`
}

// SetAppLocale sets the application name and locales in the response.
func (al *AppLocale) SetAppLocale(appName string, locales []models.Locale) {
	al.AppName = appName
	al.Locales = make([]string, len(locales))

	for i, locale := range locales {
		al.Locales[i] = locale.ID
	}
}

// SetAppLocaleSimple sets the application name and locales in the response using a slice of strings.
func (al *AppLocale) SetAppLocaleSimple(appName string, locales []string) {
	al.AppName = appName
	al.Locales = locales
}
