package models

type Country struct {
	ID     string `gorm:"primaryKey;size:2"`
	Alpha3 string `gorm:"uniqueIndex:idx_alpha3,sort:asc;size:3;not null"`
	Code   uint16 `gorm:"uniqueIndex:idx_code,sort:asc;not null"`
}
