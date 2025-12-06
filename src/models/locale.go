package models

import "database/sql"

// Locale represents a fully composed BCP47 tag broken into components.
// Examples: en, en-US, zh-Hant-TW, sr-Cyrl, es-419, en-polyton.
type Locale struct {
	ID          string         `gorm:"primaryKey;size:32"` // canonical BCP47 tag
	LanguageID  string         `gorm:"size:8;not null"`
	ScriptID    sql.NullString `gorm:"size:32"`
	TerritoryID sql.NullString `gorm:"size:32"`
	VariantID   sql.NullString `gorm:"size:32"`
	Canonical   bool           `gorm:"not null;default:true"` // whether tag is canonical form
	// Unique composite to prevent duplicate component combinations forming multiple locale IDs.
	// (GORM automatically builds index; variant/extensions add uniqueness scope.)
	// NOTE: GORM name for index is uid_locale; adjust if needed.

	// Relationships.
	Language  Language   `gorm:"foreignKey:LanguageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Script    *Script    `gorm:"foreignKey:ScriptID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Territory *Territory `gorm:"foreignKey:TerritoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Variant   *Variant   `gorm:"foreignKey:VariantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
