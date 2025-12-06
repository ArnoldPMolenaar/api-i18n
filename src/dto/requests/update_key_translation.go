package requests

import "time"

type UpdateKeyTranslation struct {
	LocaleID  string    `json:"localeId" validate:"required"`
	ValueType string    `json:"valueType" validate:"required"`
	Value     string    `json:"value" validate:"required"`
	UpdatedAt time.Time `json:"updatedAt" validate:"required"`
}
