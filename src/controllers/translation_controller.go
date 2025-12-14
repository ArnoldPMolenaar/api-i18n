package controllers

import (
	"api-i18n/main/src/errors"
	"api-i18n/main/src/services"
	"api-i18n/main/src/utils"

	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/gofiber/fiber/v2"
)

// GetTranslationsByLocaleId func for getting translations by locale id.
func GetTranslationsByLocaleId(c *fiber.Ctx) error {
	localeId := c.Params("localeId")
	appName := c.Query("app")

	// Check if the app exists.
	appAvailable, err := services.IsAppAvailable(appName)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !appAvailable {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AppNotFound, "App not found.")
	}

	// Resolve the locale id for backwards compatibility.
	resolvedLocaleId := utils.ResolveLocaleId(localeId)
	if resolvedLocaleId == nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.LocaleNotFound, "Locale not found.")
	}

	// Check if locales are set in the app.
	hasLocales, err := HasAppLocales(appName, *resolvedLocaleId)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !hasLocales {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.LocaleNotFound, "Locale not found in app.")
	}

	translations, err := services.GetTranslationsByLocaleId(appName, *resolvedLocaleId)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(translations)
}
