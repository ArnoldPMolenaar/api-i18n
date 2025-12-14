package responses

import (
	"api-i18n/main/src/models"
	"time"
)

type KeyTranslation struct {
	KeyID     uint      `json:"keyId"`
	LocaleID  string    `json:"localeId"`
	ValueType string    `json:"valueType"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SetKeyTranslation func to set key translation response model from key translation model.
func (kt *KeyTranslation) SetKeyTranslation(keyTranslation *models.KeyTranslation) {
	kt.KeyID = keyTranslation.KeyID
	kt.LocaleID = keyTranslation.LocaleID
	kt.ValueType = keyTranslation.ValueType.String()
	kt.Value = keyTranslation.Value
	kt.CreatedAt = keyTranslation.CreatedAt
	kt.UpdatedAt = keyTranslation.UpdatedAt
}
