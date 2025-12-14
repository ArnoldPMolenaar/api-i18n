package controllers

import (
	"api-i18n/main/src/database"
	"api-i18n/main/src/dto/requests"
	"api-i18n/main/src/dto/responses"
	"api-i18n/main/src/models"
	"api-i18n/main/src/services"
	"api-i18n/main/src/utils"

	util "github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"

	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
)

// HasAppLocales to check if an app has all the given locales.
func HasAppLocales(appName string, locales ...string) (bool, error) {
	var count int64

	tx := database.Pg.Model(&models.App{}).
		Joins("JOIN app_locales ON app_locales.app_name = apps.name").
		Where("name = ? AND locale_id IN ?", appName, locales).
		Count(&count)

	if tx.Error != nil {
		return false, tx.Error
	}

	return int(count) == len(locales), nil
}

// GetAppLocales to get the locales of an app.
func GetAppLocales(c *fiber.Ctx) error {
	// Get the appName parameter from the URL.
	appNameParam := c.Params("name")
	if appNameParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.MissingRequiredParam, "App Name is required.")
	}

	// Check if the app exists.
	appAvailable, err := services.IsAppAvailable(appNameParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !appAvailable {
		return errorutil.Response(c, fiber.StatusNotFound, errorutil.NotFound, "App not found.")
	}

	// Get the app locales.
	locales, err := services.GetAppLocales(appNameParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the app.
	response := responses.AppLocale{}
	response.SetAppLocale(appNameParam, locales)

	return c.JSON(response)
}

// CreateApp method to create an app.
func CreateApp(c *fiber.Ctx) error {
	// Parse the request.
	request := requests.CreateApp{}
	if err := c.BodyParser(&request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate document fields.
	validate := util.NewValidator()
	if err := validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	// Create the app.
	app, err := services.CreateApp(request.Name)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	}

	// Return the document.
	response := responses.App{}
	response.SetApp(app)

	return c.JSON(response)
}

// SetAppLocales to set the locale of an app.
func SetAppLocales(c *fiber.Ctx) error {
	// Get the appName parameter from the URL.
	appNameParam := c.Params("name")
	if appNameParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.MissingRequiredParam, "App Name is required.")
	}

	// Parse the request.
	request := requests.SetAppLocale{}
	if err := c.BodyParser(&request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate document fields.
	validate := util.NewValidator()
	if err := validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	// Check if the app exists.
	appAvailable, err := services.IsAppAvailable(appNameParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !appAvailable {
		return errorutil.Response(c, fiber.StatusNotFound, errorutil.NotFound, "App not found.")
	}

	for _, locale := range request.Locales {
		id := utils.ResolveLocaleId(locale)
		if id == nil {
			return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "Invalid locale: "+locale)
		}
	}

	// Set the app locale.
	if err := services.SetAppLocales(appNameParam, request.Locales); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the app.
	response := responses.AppLocale{}
	response.SetAppLocaleSimple(appNameParam, request.Locales)

	return c.JSON(response)
}
