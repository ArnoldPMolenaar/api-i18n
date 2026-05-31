package responses

import (
	"api-i18n/main/src/models"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/utils"
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
	k.CategoryID = utils.PtrFromNull[uint](key.CategoryID)
	k.AppName = key.AppName
	k.Name = key.Name
	k.DisabledAt = utils.PtrFromNullTime(key.DisabledAt)
	k.CreatedAt = key.CreatedAt
	k.UpdatedAt = key.UpdatedAt

	if key.Category != nil {
		k.CategoryName = &key.Category.Name
		k.CategoryDisabledAt = utils.PtrFromNullTime(key.Category.DisabledAt)
	}
}
