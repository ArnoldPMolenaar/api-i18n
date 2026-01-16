package routes

import (
	"api-i18n/main/src/controllers"

	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create public routes group.
	route := a.Group("/v1")

	// Register route group for /v1/apps.
	apps := route.Group("/apps")
	apps.Get("/:name/locales", controllers.GetAppLocales)

	// Register route group for /v1/territories.
	territories := route.Group("/territories")
	territories.Get("/lookup", controllers.GetTerritoryLookup)

	// Register route group for /v1/locales.
	locales := route.Group("/locales")
	locales.Get("/lookup", controllers.GetLocaleLookup)

	// Register route group for /v1/translations.
	translations := route.Group("/translations")
	translations.Get("/:localeId", controllers.GetTranslationsByLocaleId)

	// Register route group for /v1/phones.
	phones := route.Group("/phones")
	phones.Get("/lookup", controllers.GetPhoneLookup)
	phones.Get("/validate", controllers.GetPhoneNumberValidation)
	phones.Get("/format", controllers.GetPhoneNumberFormat)
}
