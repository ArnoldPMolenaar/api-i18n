package models

// TerritoryName represents a localized name for a given language/locale.
type TerritoryName struct {
	TerritoryID string `gorm:"primaryKey;size:32"`
	LocaleID    string `gorm:"primaryKey;size:32"`
	Name        string `gorm:"not null"`

	// Relationships.
	Territory Territory `gorm:"foreignKey:TerritoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Locale    Locale    `gorm:"foreignKey:LocaleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
