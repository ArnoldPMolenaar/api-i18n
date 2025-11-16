package models

type Language struct {
	ID   string `gorm:"primaryKey;size:4"`
	Name string `gorm:"uniqueIndex:idx_name,sort:asc;not null"`
}
