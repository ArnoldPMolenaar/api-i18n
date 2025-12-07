package controllers

import (
	"api-i18n/main/src/dto/requests"
	"api-i18n/main/src/dto/responses"
	"api-i18n/main/src/services"

	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"

	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
)

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
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
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
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}

	// Check if the app exists.
	appAvailable, err := services.IsAppAvailable(appNameParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !appAvailable {
		return errorutil.Response(c, fiber.StatusNotFound, errorutil.NotFound, "App not found.")
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
