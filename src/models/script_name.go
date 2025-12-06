package models

// ScriptName represents a script name for a given language/locale.
type ScriptName struct {
	ScriptID string `gorm:"primaryKey;size:32"`
	LocaleID string `gorm:"primaryKey;size:32"`
	Name     string `gorm:"not null"`

	// Relationships.
	Script Script `gorm:"foreignKey:ScriptID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Locale Locale `gorm:"foreignKey:LocaleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
