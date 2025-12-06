package models

// LocaleName stores the localized display name of a target locale in the context of a viewer locale.
// Example: LocaleIDViewer = "nl, LocaleIDTarget = "en", Name = "Engels (Verenigde Staten)".
// Composite primary key ensures single translation per (viewer,target) pair.
type LocaleName struct {
	LocaleIDViewer string `gorm:"primaryKey;size:32"`
	LocaleIDTarget string `gorm:"primaryKey;size:32"`
	Name           string `gorm:"not null"`

	LocaleViewer Locale `gorm:"foreignKey:LocaleIDViewer;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	LocaleTarget Locale `gorm:"foreignKey:LocaleIDTarget;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
