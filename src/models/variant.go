package models

// Variant represents a code subtag like "1901" (Traditional German orthography).
type Variant struct {
	ID string `gorm:"primaryKey;size:32"` // Title-case code
}
