package models

import "api-i18n/main/src/enums"

// Territory represents a geographic territory subtag (ISO 3166-1 alpha-2 or UN M.49 numeric).
// Examples: US, NL, 419.
// Type allows distinguishing country vs numeric codes.
type Territory struct {
	ID   string              `gorm:"primaryKey;size:32"`
	Type enums.TerritoryType `gorm:"type:territory_type;default:country;not null"` // e.g. country, numeric
}
