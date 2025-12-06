package models

// Language represents an ISO 639 language without region/script variants.
// Examples: en, nl, zh, sr, pt
type Language struct {
	ID string `gorm:"primaryKey;size:8"`
}
