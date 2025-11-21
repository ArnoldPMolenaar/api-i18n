package requests

import "time"

type UpdateKeyTranslation struct {
	LanguageID string    `json:"languageId" validate:"required"`
	ValueType  string    `json:"valueType" validate:"required"`
	Value      string    `json:"value" validate:"required"`
	UpdatedAt  time.Time `json:"updatedAt" validate:"required"`
}
