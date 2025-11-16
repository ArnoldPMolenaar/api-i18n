package models

type App struct {
	Name string `gorm:"primaryKey:true;autoIncrement:false"`

	// Relationships.
	Languages []Language `gorm:"many2many:app_languages;foreignKey:Name;joinForeignKey:AppName;references:ID;joinReferences:LanguageId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
