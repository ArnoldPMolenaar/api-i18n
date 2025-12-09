package requests

import "time"

type CreateCategory struct {
	Name       string     `json:"name" validate:"required"`
	DisabledAt *time.Time `json:"disabledAt"`
}
