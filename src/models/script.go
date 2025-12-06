package models

// Script represents an ISO 15924 script subtag (e.g. Latn, Cyrl, Hans, Hant).
type Script struct {
	ID string `gorm:"primaryKey;size:32"` // Title-case code
}
