package database

import (
	"api-i18n/main/src/enums"
	"api-i18n/main/src/models"
	"database/sql"
	"encoding/json"
	"os"
	"slices"
	"unicode"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

// Migrate the database schema.
// See: https://gorm.io/docs/migration.html#Auto-Migration
func Migrate(db *gorm.DB) error {
	// Adds the value type enum type to the database.
	if tx := db.Exec(`DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'value_type') THEN 
			CREATE TYPE value_type AS ENUM ('text', 'html'); 
		END IF; 
	END $$;`); tx.Error != nil {
		return tx.Error
	}

	// Adds the region type enum type to the database.
	if tx := db.Exec(`DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'territory_type') THEN 
			CREATE TYPE territory_type AS ENUM ('country', 'numeric'); 
		END IF; 
	END $$;`); tx.Error != nil {
		return tx.Error
	}

	// Updated migration set: normalized models + existing domain models.
	err := db.AutoMigrate(&models.Language{}, &models.Script{}, &models.Territory{}, &models.Variant{}, &models.Locale{}, &models.LocaleName{}, &models.ScriptName{}, &models.TerritoryName{}, &models.VariantName{}, &models.App{}, &models.Category{}, &models.Key{}, &models.KeyTranslation{})
	if err != nil {
		return err
	}

	// -- Start CLDR script migration --
	if err := seedCLDRData(db); err != nil {
		return err
	}
	// -- End CLDR script migration --

	return nil
}

// readJSONFile reads a JSON file from the given path and unmarshals it into a map.
func readJSONFile(path string) (map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var doc map[string]interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&doc); err != nil {
		return nil, err
	}
	return doc, nil
}

// isNumeric checks if a string consists only of numeric characters.
func isNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return len(s) > 0
}

// seedCLDRData seeds the database with CLDR data for languages, scripts, territories, variants, locales, and their names.
func seedCLDRData(db *gorm.DB) error {
	const cldrBasePath = "src/database/fixtures/cldr-json/cldr-json/"

	var languageCount, scriptCount, territoryCount, variantCount, localeCount, scriptNameCount, territoryNameCount, variantNameCount, localeNameCount int64
	_ = db.Model(&models.Language{}).Count(&languageCount)
	_ = db.Model(&models.Script{}).Count(&scriptCount)
	_ = db.Model(&models.Territory{}).Count(&territoryCount)
	_ = db.Model(&models.Variant{}).Count(&variantCount)
	_ = db.Model(&models.Locale{}).Count(&localeCount)
	_ = db.Model(&models.ScriptName{}).Count(&scriptNameCount)
	_ = db.Model(&models.TerritoryName{}).Count(&territoryNameCount)
	_ = db.Model(&models.VariantName{}).Count(&variantNameCount)
	_ = db.Model(&models.LocaleName{}).Count(&localeNameCount)

	if languageCount > 0 && scriptCount > 0 && territoryCount > 0 && variantCount > 0 && localeCount > 0 && scriptNameCount > 0 && territoryNameCount > 0 && variantNameCount > 0 && localeNameCount > 0 {
		return nil // Data already seeded; skip.
	}

	log.Info("Seeding CLDR data into the database...")

	languages := make([]models.Language, 0)
	scripts := make([]models.Script, 0)
	scriptNames := make([]models.ScriptName, 0)
	territories := make([]models.Territory, 0)
	territoryNames := make([]models.TerritoryName, 0)
	variants := make([]models.Variant, 0)
	variantNames := make([]models.VariantName, 0)
	locales := make([]models.Locale, 0)
	localeNames := make([]models.LocaleName, 0)

	// Get all optional locales.
	localesDoc, err := readJSONFile(cldrBasePath + "cldr-core/availableLocales.json")
	if err != nil {
		return err
	}
	optionalLocales := localesDoc["availableLocales"].(map[string]interface{})["full"].([]interface{})

	for _, loc := range optionalLocales {
		locale, ok := loc.(string)
		if !ok || locale == "" {
			continue
		}

		localeNamesDoc, err := readJSONFile(cldrBasePath + "cldr-localenames-full/main/" + locale + "/localeDisplayNames.json")
		if err != nil {
			continue
		}

		// Set language if not exists.
		languageID, ok := localeNamesDoc["main"].(map[string]interface{})[locale].(map[string]interface{})["identity"].(map[string]interface{})["language"].(string)
		if !ok || languageID == "" {
			continue
		}
		language := models.Language{ID: languageID}
		if !slices.Contains(languages, language) {
			languages = append(languages, language)
		}

		// Set script if not exists.
		scriptID, ok := localeNamesDoc["main"].(map[string]interface{})[locale].(map[string]interface{})["identity"].(map[string]interface{})["script"].(string)
		if ok && scriptID != "" {
			script := models.Script{ID: scriptID}
			if !slices.Contains(scripts, script) {
				scripts = append(scripts, script)
			}
		}

		// Set territory if not exists.
		territoryID, ok := localeNamesDoc["main"].(map[string]interface{})[locale].(map[string]interface{})["identity"].(map[string]interface{})["territory"].(string)
		if ok && territoryID != "" {
			territoryType := enums.COUNTRY
			if isNumeric(territoryID) {
				territoryType = enums.NUMERIC
			}
			territory := models.Territory{ID: territoryID, Type: territoryType}
			if !slices.Contains(territories, territory) {
				territories = append(territories, territory)
			}
		}

		// Set variant if not exists.
		variantID, ok := localeNamesDoc["main"].(map[string]interface{})[locale].(map[string]interface{})["identity"].(map[string]interface{})["variant"].(string)
		if ok && variantID != "" {
			variant := models.Variant{ID: variantID}
			if !slices.Contains(variants, variant) {
				variants = append(variants, variant)
			}
		}

		// Set locale.
		localeModel := models.Locale{ID: locale, LanguageID: languageID}
		if scriptID != "" {
			localeModel.ScriptID = sql.NullString{String: scriptID, Valid: true}
		}
		if territoryID != "" {
			localeModel.TerritoryID = sql.NullString{String: territoryID, Valid: true}
		}
		if variantID != "" {
			localeModel.VariantID = sql.NullString{String: variantID, Valid: true}
		}
		locales = append(locales, localeModel)

		// Set script names.
		scriptNamesDoc, err := readJSONFile(cldrBasePath + "cldr-localenames-full/main/" + locale + "/scripts.json")
		if err == nil {
			scriptNamesMap := scriptNamesDoc["main"].(map[string]interface{})[locale].(map[string]interface{})["localeDisplayNames"].(map[string]interface{})["scripts"].(map[string]interface{})
			for id, name := range scriptNamesMap {
				n, ok := name.(string)
				if !ok || n == "" {
					continue
				}
				script := models.Script{ID: id}
				if !slices.Contains(scripts, script) {
					scripts = append(scripts, script)
				}
				scriptName := models.ScriptName{
					ScriptID: id,
					LocaleID: locale,
					Name:     n,
				}
				scriptNames = append(scriptNames, scriptName)
			}
		}

		// Set territory names.
		territoryNamesDoc, err := readJSONFile(cldrBasePath + "cldr-localenames-full/main/" + locale + "/territories.json")
		if err == nil {
			territoryNamesMap := territoryNamesDoc["main"].(map[string]interface{})[locale].(map[string]interface{})["localeDisplayNames"].(map[string]interface{})["territories"].(map[string]interface{})
			for id, name := range territoryNamesMap {
				n, ok := name.(string)
				if !ok || n == "" {
					continue
				}

				territoryType := enums.COUNTRY
				if isNumeric(id) {
					territoryType = enums.NUMERIC
				}
				territory := models.Territory{ID: id, Type: territoryType}
				if !slices.Contains(territories, territory) {
					territories = append(territories, territory)
				}

				territoryName := models.TerritoryName{
					TerritoryID: id,
					LocaleID:    locale,
					Name:        n,
				}
				territoryNames = append(territoryNames, territoryName)
			}
		}

		// Set variant names.
		variantNamesDoc, err := readJSONFile(cldrBasePath + "cldr-localenames-full/main/" + locale + "/variants.json")
		if err == nil {
			variantNamesMap := variantNamesDoc["main"].(map[string]interface{})[locale].(map[string]interface{})["localeDisplayNames"].(map[string]interface{})["variants"].(map[string]interface{})
			for id, name := range variantNamesMap {
				n, ok := name.(string)
				if !ok || n == "" {
					continue
				}

				variant := models.Variant{ID: id}
				if !slices.Contains(variants, variant) {
					variants = append(variants, variant)
				}

				variantName := models.VariantName{
					VariantID: id,
					LocaleID:  locale,
					Name:      n,
				}
				variantNames = append(variantNames, variantName)
			}
		}

		// Set locale names.
		languageNamesDoc, err := readJSONFile(cldrBasePath + "cldr-localenames-full/main/" + locale + "/languages.json")
		if err == nil {
			languageNamesMap := languageNamesDoc["main"].(map[string]interface{})[locale].(map[string]interface{})["localeDisplayNames"].(map[string]interface{})["languages"].(map[string]interface{})
			for id, name := range languageNamesMap {
				n, ok := name.(string)
				if !ok || n == "" {
					continue
				}
				localeName := models.LocaleName{
					LocaleIDViewer: locale,
					LocaleIDTarget: id,
					Name:           n,
				}
				localeNames = append(localeNames, localeName)
			}
		}
	}

	// Bulk insert collected data.

	if languageCount == 0 && len(languages) > 0 {
		log.Info("Inserting languages...")
		if tx := db.Create(&languages); tx.Error != nil {
			return tx.Error
		}
	}

	if scriptCount == 0 && len(scripts) > 0 {
		log.Info("Inserting scripts...")
		if tx := db.Create(&scripts); tx.Error != nil {
			return tx.Error
		}
	}

	if territoryCount == 0 && len(territories) > 0 {
		log.Info("Inserting territories...")
		if tx := db.Create(&territories); tx.Error != nil {
			return tx.Error
		}
	}

	if variantCount == 0 && len(variants) > 0 {
		log.Info("Inserting variants...")
		if tx := db.Create(&variants); tx.Error != nil {
			return tx.Error
		}
	}

	if localeCount == 0 && len(locales) > 0 {
		log.Info("Inserting locales...")
		if tx := db.Create(&locales); tx.Error != nil {
			return tx.Error
		}
	}

	if scriptNameCount == 0 && len(scriptNames) > 0 {
		log.Info("Inserting script names...")
		for i := range scriptNames {
			if tx := db.Create(&scriptNames[i]); tx.Error != nil {
				return tx.Error
			}
		}
	}

	if territoryNameCount == 0 && len(territoryNames) > 0 {
		log.Info("Inserting territory names...")
		for i := range territoryNames {
			if tx := db.Create(&territoryNames[i]); tx.Error != nil {
				return tx.Error
			}
		}
	}

	if variantNameCount == 0 && len(variantNames) > 0 {
		log.Info("Inserting variant names...")
		for i := range variantNames {
			if tx := db.Create(&variantNames[i]); tx.Error != nil {
				return tx.Error
			}
		}
	}

	if localeNameCount == 0 && len(localeNames) > 0 {
		log.Info("Inserting locale names...")
		for i := range localeNames {
			var foundViewerId, foundTargetId bool
			for _, loc := range optionalLocales {
				locale, ok := loc.(string)
				if !ok || locale == "" {
					continue
				}
				if !foundViewerId && localeNames[i].LocaleIDViewer == locale {
					foundViewerId = true
				}
				if !foundTargetId && localeNames[i].LocaleIDTarget == locale {
					foundTargetId = true
				}
				if foundViewerId && foundTargetId {
					break
				}
			}

			if foundViewerId && foundTargetId {
				if tx := db.Create(&localeNames[i]); tx.Error != nil {
					return tx.Error
				}
			}
		}
	}

	return nil
}
