package responses

import (
	"api-i18n/main/src/models"
	"time"
)

type PaginatedKey struct {
	ID                 uint       `json:"id"`
	CategoryID         *uint      `json:"categoryId"`
	Name               string     `json:"name"`
	AppName            string     `json:"appName"`
	DisabledAt         *time.Time `json:"disabledAt"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	CategoryName       *string    `json:"categoryName"`
	CategoryDisabledAt *time.Time `json:"categoryDisabledAt"`
}

// SetPaginatedKey method to set key data from models.Key{}.
func (k *PaginatedKey) SetPaginatedKey(key *models.Key) {
	k.ID = key.ID
	k.CategoryID = func() *uint {
		if key.CategoryID.Valid {
			return &key.CategoryID.V
		}
		return nil
	}()
	k.AppName = key.AppName
	k.Name = key.Name
	k.CreatedAt = key.CreatedAt
	k.UpdatedAt = key.UpdatedAt
	k.DisabledAt = func() *time.Time {
		if key.DisabledAt.Valid {
			return &key.DisabledAt.Time
		}
		return nil
	}()

	if key.Category != nil {
		k.CategoryName = &key.Category.Name
		k.CategoryDisabledAt = func() *time.Time {
			if key.Category.DisabledAt.Valid {
				return &key.Category.DisabledAt.Time
			}
			return nil
		}()
	}
}
