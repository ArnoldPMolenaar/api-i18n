package models

type CountryName struct {
	CountryID  string `gorm:"primaryKey;size:2"`
	LanguageID string `gorm:"primaryKey;size:4"`
	Name       string `gorm:"not null"`

	// Relationships.
	Country  Country  `gorm:"foreignKey:CountryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Language Language `gorm:"foreignKey:LanguageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
