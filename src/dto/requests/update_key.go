package requests

import "time"

type UpdateKey struct {
	CategoryID   *uint                  `json:"categoryId"`
	Name         string                 `json:"name" validate:"required"`
	Description  *string                `json:"description"`
	DisabledAt   *time.Time             `json:"disabledAt"`
	UpdatedAt    time.Time              `json:"updatedAt" validate:"required"`
	Translations []UpdateKeyTranslation `json:"translations" validate:"required,dive"`
}
