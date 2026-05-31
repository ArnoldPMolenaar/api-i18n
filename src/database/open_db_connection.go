package database

import (
	"api-i18n/main/src/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/database"
	"gorm.io/gorm"
)

var Pg *gorm.DB

// OpenDBConnection Start a new database connection.
// Also tries to migrate the database schema.
func OpenDBConnection() error {
	// Open connection to database.
	db, err := database.PostgresSQLConnection()
	if err != nil {
		return err
	}

	// Migrate the database schema.
	err = Migrate(db)
	if err != nil {
		return err
	}

	// Set the global DB variable.
	Pg = db

	return nil
}

// ReadinessCheck verifies that the database connection is initialized and reachable.
func ReadinessCheck() error {
	if Pg == nil {
		return errors.New("database connection is not initialized")
	}

	sqlDB, err := Pg.DB()
	if err != nil {
		return fmt.Errorf("database sql handle unavailable: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// MigrationReadinessCheck verifies that required tables and enum types exist.
func MigrationReadinessCheck() error {
	if Pg == nil {
		return errors.New("database connection is not initialized")
	}

	requiredTables := []any{
		// Core CLDR entities.
		&models.Language{},
		&models.Script{},
		&models.Territory{},
		&models.Variant{},
		&models.Locale{},

		// Localized display name entities.
		&models.LocaleName{},
		&models.ScriptName{},
		&models.TerritoryName{},
		&models.VariantName{},

		// API domain entities.
		&models.App{},
		&models.Category{},
		&models.Key{},
		&models.KeyTranslation{},
	}
	for _, table := range requiredTables {
		if !Pg.Migrator().HasTable(table) {
			return fmt.Errorf("missing required table for %T", table)
		}
	}

	// Join table used by App <-> Locale many-to-many relation.
	if !Pg.Migrator().HasTable("app_locales") {
		return errors.New("missing required table app_locales")
	}

	valueTypeLabels, err := getEnumLabels("value_type")
	if err != nil {
		return fmt.Errorf("value_type enum check failed: %w", err)
	}

	expectedValueTypeLabels := []string{"text", "html", "json"}
	if len(valueTypeLabels) != len(expectedValueTypeLabels) {
		return fmt.Errorf("value_type enum labels mismatch: have %v, want %v", valueTypeLabels, expectedValueTypeLabels)
	}

	for i := range valueTypeLabels {
		if valueTypeLabels[i] != expectedValueTypeLabels[i] {
			return fmt.Errorf("value_type enum labels mismatch: have %v, want %v", valueTypeLabels, expectedValueTypeLabels)
		}
	}

	territoryTypeLabels, err := getEnumLabels("territory_type")
	if err != nil {
		return fmt.Errorf("territory_type enum check failed: %w", err)
	}

	expectedTerritoryTypeLabels := []string{"country", "numeric"}
	if len(territoryTypeLabels) != len(expectedTerritoryTypeLabels) {
		return fmt.Errorf("territory_type enum labels mismatch: have %v, want %v", territoryTypeLabels, expectedTerritoryTypeLabels)
	}

	for i := range territoryTypeLabels {
		if territoryTypeLabels[i] != expectedTerritoryTypeLabels[i] {
			return fmt.Errorf("territory_type enum labels mismatch: have %v, want %v", territoryTypeLabels, expectedTerritoryTypeLabels)
		}
	}

	return nil
}

func getEnumLabels(typeName string) ([]string, error) {
	rows, err := Pg.Raw(`
		SELECT e.enumlabel
		FROM pg_type t
		JOIN pg_enum e ON t.oid = e.enumtypid
		WHERE t.typname = ?
		ORDER BY e.enumsortorder
	`, typeName).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	labels := make([]string, 0)
	for rows.Next() {
		var label string
		if scanErr := rows.Scan(&label); scanErr != nil {
			return nil, scanErr
		}
		labels = append(labels, label)
	}

	if len(labels) == 0 {
		return nil, fmt.Errorf("enum type %q not found", typeName)
	}

	return labels, nil
}
