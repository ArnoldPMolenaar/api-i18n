package requests

type InsertKeyTranslation struct {
	LanguageID string `json:"languageId" validate:"required"`
	ValueType  string `json:"valueType" validate:"required"`
	Value      string `json:"value" validate:"required"`
}
