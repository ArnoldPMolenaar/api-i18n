package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	DisabledAt sql.NullTime
	Name       string `gorm:"uniqueIndex:idx_categories_name,sort:asc;not null"`

	// Relationships.
	Keys []Key
}
