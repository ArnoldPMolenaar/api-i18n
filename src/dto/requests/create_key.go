package requests

import "time"

type CreateKey struct {
	CategoryID   *uint                  `json:"categoryId"`
	AppName      string                 `json:"appName" validate:"required"`
	Name         string                 `json:"name" validate:"required"`
	Description  *string                `json:"description"`
	DisabledAt   *time.Time             `json:"disabledAt"`
	Translations []CreateKeyTranslation `json:"translations" validate:"required,min=1,dive"`
}
