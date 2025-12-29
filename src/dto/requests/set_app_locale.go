package requests

// SetAppLocale struct for setting the locale of an app.
type SetAppLocale struct {
	Locales []string `json:"locales" validate:"required,min=1"`
}
