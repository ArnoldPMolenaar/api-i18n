package requests

import "time"

type UpdateCategory struct {
	Name       string     `json:"name" validate:"required"`
	UpdatedAt  time.Time  `json:"updatedAt" validate:"required"`
	DisabledAt *time.Time `json:"disabledAt"`
}
