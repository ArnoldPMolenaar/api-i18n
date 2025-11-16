package database

import (
	"api-i18n/main/src/models"

	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// Migrate the database schema.
// See: https://gorm.io/docs/migration.html#Auto-Migration
func Migrate(db *gorm.DB) error {
	// Adds the size enum type to the database.
	if tx := db.Exec(`DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'value_type') THEN 
			CREATE TYPE value_type AS ENUM ('text', 'html'); 
		END IF; 
	END $$;`); tx.Error != nil {
		return tx.Error
	}

	err := db.AutoMigrate(&models.Language{}, &models.Country{}, &models.CountryName{}, &models.App{}, &models.Category{}, &models.Key{}, &models.KeyTranslation{})
	if err != nil {
		return err
	}
	// --- Begin seed script ---
	if err := seedLanguages(db); err != nil {
		return err
	}
	if err := seedCountries(db); err != nil {
		return err
	}
	if err := seedCountryNames(db); err != nil {
		return err
	}
	// --- End seed script ---

	return nil
}

// seedLanguages imports Languages.csv if the languages table is empty.
func seedLanguages(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Language{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil // already seeded or manually populated
	}

	path := filepath.Join("src", "database", "fixtures", "languages.csv")
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Printf("seedLanguages: error closing file %s: %v", path, cerr)
		}
	}()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	var languages []models.Language
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(rec) < 2 {
			continue
		}
		id := strings.TrimSpace(rec[0])
		name := strings.TrimSpace(rec[1])
		if id == "" || name == "" {
			continue
		}
		// Current schema restricts Language.ID size to 4. Skip longer codes to avoid errors.
		if len(id) > 4 {
			log.Printf("seedLanguages: skipping language code '%s' (>4 chars). Consider widening Language.ID size if needed.", id)
			continue
		}
		languages = append(languages, models.Language{ID: id, Name: name})
	}
	if len(languages) == 0 {
		return nil
	}
	if err := db.Create(&languages).Error; err != nil {
		return err
	}
	return nil
}

// seedCountries imports Countries.csv if the countries table is empty.
func seedCountries(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Country{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	path := filepath.Join("src", "database", "fixtures", "countries.csv")
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Printf("seedCountries: error closing file %s: %v", path, cerr)
		}
	}()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	var countries []models.Country
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(rec) < 3 {
			continue
		}
		id := strings.TrimSpace(rec[0]) // alpha-2
		alpha3 := strings.TrimSpace(rec[1])
		codeStr := strings.TrimSpace(rec[2])
		if id == "" || alpha3 == "" || codeStr == "" {
			continue
		}
		// Current schema restricts Country.ID size to 2. Skip longer codes to avoid errors.
		if len(id) > 2 {
			log.Printf("seedCountries: skipping country code '%s' (>2 chars). Consider widening Country.ID size if needed.", id)
			continue
		}
		// Current schema restricts Country.Alpha3 size to 3. Skip longer codes to avoid errors.
		if len(alpha3) > 3 {
			log.Printf("seedCountries: skipping country alpha3 '%s' (>3 chars). Consider widening Country.Alpha3 size if needed.", id)
			continue
		}
		// Parse numeric code (may exceed uint16 range present in model; validation).
		codeInt, err := strconv.Atoi(codeStr)
		if err != nil {
			return err
		}
		if codeInt < 0 || codeInt > 65535 {
			// Warn and skip to avoid silent overflow.
			log.Printf("seedCountries: skipping country %s numeric code %d outside uint8 range; consider changing Country.Code to a larger integer type.", id, codeInt)
			continue
		}
		countries = append(countries, models.Country{ID: id, Alpha3: alpha3, Code: uint16(codeInt)})
	}
	if len(countries) == 0 {
		return nil
	}
	if err := db.Create(&countries).Error; err != nil {
		return err
	}
	return nil
}

// seedCountryNames imports CountryNames.csv if the country_names table is empty.
func seedCountryNames(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.CountryName{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// Ensure prerequisite tables have data (languages & countries). If not, skip to avoid FK errors.
	var langCount, countryCount int64
	if err := db.Model(&models.Language{}).Count(&langCount).Error; err != nil {
		return err
	}
	if err := db.Model(&models.Country{}).Count(&countryCount).Error; err != nil {
		return err
	}
	if langCount == 0 || countryCount == 0 {
		return errors.New("seedCountryNames: prerequisite tables empty (languages or countries); cannot seed country names")
	}

	path := filepath.Join("src", "database", "fixtures", "country_names.csv")
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Printf("seedCountryNames: error closing file %s: %v", path, cerr)
		}
	}()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	var names []models.CountryName
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(rec) < 3 {
			continue
		}
		languageID := strings.TrimSpace(rec[0])
		countryID := strings.TrimSpace(rec[1])
		name := strings.TrimSpace(rec[2])
		if languageID == "" || countryID == "" || name == "" {
			continue
		}
		if len(languageID) > 4 { // schema restriction
			log.Printf("seedCountryNames: skipping entry language '%s' (>4 chars) for country '%s'", languageID, countryID)
			continue
		}
		if len(countryID) > 2 {
			log.Printf("seedCountryNames: skipping entry invalid country code '%s'", countryID)
			continue
		}
		names = append(names, models.CountryName{CountryID: countryID, LanguageID: languageID, Name: name})
	}
	if len(names) == 0 {
		return nil
	}
	// Batch create with controlled chunking to avoid huge insert statements (optional, here simple).
	if err := db.Create(&names).Error; err != nil {
		return err
	}
	return nil
}
