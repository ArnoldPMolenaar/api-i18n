package requests

type InsertKeyTranslation struct {
	LocaleID  string `json:"localeId" validate:"required"`
	ValueType string `json:"valueType" validate:"required"`
	Value     string `json:"value" validate:"required"`
}
