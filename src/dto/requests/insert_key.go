package requests

import "time"

type InsertKey struct {
	CategoryID   *uint                  `json:"categoryId"`
	AppName      string                 `json:"appName" validate:"required"`
	Name         string                 `json:"name" validate:"required"`
	Description  *string                `json:"description"`
	DisabledAt   *time.Time             `json:"disabledAt"`
	Translations []InsertKeyTranslation `json:"translations" validate:"required,dive"`
}
