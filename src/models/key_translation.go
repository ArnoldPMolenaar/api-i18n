package models

import (
	"api-i18n/main/src/enums"
	"time"

	"gorm.io/gorm"
)

type KeyTranslation struct {
	KeyID     uint            `gorm:"primaryKey"`
	LocaleID  string          `gorm:"primaryKey;size:32"`
	ValueType enums.ValueType `gorm:"not null;type:value_type;default:text"`
	Value     string          `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relationships.
	Key    Key    `gorm:"foreignKey:KeyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Locale Locale `gorm:"foreignKey:LocaleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
