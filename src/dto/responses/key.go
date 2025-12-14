package responses

import (
	"api-i18n/main/src/models"
	"time"
)

type Key struct {
	ID           uint             `json:"id"`
	AppName      string           `json:"appName"`
	CategoryID   *uint            `json:"categoryId"`
	Name         string           `json:"name"`
	Description  *string          `json:"description"`
	DisabledAt   *time.Time       `json:"disabledAt"`
	CreatedAt    time.Time        `json:"createdAt"`
	UpdatedAt    time.Time        `json:"updatedAt"`
	Category     *Category        `json:"category"`
	Translations []KeyTranslation `json:"translations"`
}

// SetKey func to set key response from key model.
func (k *Key) SetKey(key *models.Key) {
	k.ID = key.ID
	k.AppName = key.AppName

	if key.CategoryID.Valid {
		k.CategoryID = &key.CategoryID.V
	}

	k.Name = key.Name

	if key.Description.Valid {
		k.Description = &key.Description.String
	}

	if key.DisabledAt.Valid {
		k.DisabledAt = &key.DisabledAt.Time
	}

	k.CreatedAt = key.CreatedAt
	k.UpdatedAt = key.UpdatedAt

	if key.Category != nil {
		k.Category = &Category{}
		k.Category.SetCategory(key.Category)
	}

	if len(key.Translations) > 0 {
		k.Translations = make([]KeyTranslation, len(key.Translations))
		for i, translation := range key.Translations {
			k.Translations[i] = KeyTranslation{}
			k.Translations[i].SetKeyTranslation(&translation)
		}
	}
}
