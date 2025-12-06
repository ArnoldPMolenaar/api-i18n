package models

// VariantName represents a variant name for a given language/locale.
type VariantName struct {
	VariantID string `gorm:"primaryKey;size:32"`
	LocaleID  string `gorm:"primaryKey;size:32"`
	Name      string `gorm:"not null"`

	// Relationships.
	Variant Variant `gorm:"foreignKey:VariantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Locale  Locale  `gorm:"foreignKey:LocaleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
