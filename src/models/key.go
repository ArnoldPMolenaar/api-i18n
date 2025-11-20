package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type Key struct {
	gorm.Model
	DisabledAt  sql.NullTime
	AppName     string         `gorm:"not null;index:idx_app_category_name,unique,priority:1"`
	CategoryID  sql.Null[uint] `gorm:"index:idx_app_category_name,unique,priority:2"`
	Name        string         `gorm:"not null;index:idx_app_category_name,unique,priority:3"`
	Description sql.NullString

	// Relationships.
	App          App              `gorm:"foreignKey:AppName;references:Name;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Category     *Category        `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Translations []KeyTranslation `gorm:"foreignKey:KeyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
