package models

type App struct {
	Name string `gorm:"primaryKey:true;autoIncrement:false"`

	// Relationships.
	Locales []Locale `gorm:"many2many:app_locales;foreignKey:Name;joinForeignKey:AppName;references:ID;joinReferences:LocaleId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
