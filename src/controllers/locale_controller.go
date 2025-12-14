package controllers

import (
	"api-i18n/main/src/dto/responses"
	"api-i18n/main/src/errors"
	"api-i18n/main/src/services"
	"api-i18n/main/src/utils"

	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/gofiber/fiber/v2"
)

// GetLocaleLookup func for getting locale lookup by locale ID, optional name filter.
func GetLocaleLookup(c *fiber.Ctx) error {
	localeIDParam := c.Query("localeId")
	if localeIDParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "localeId query parameter is required.")
	}

	nameParam := c.Query("name")
	var name *string
	if nameParam != "" {
		name = &nameParam
	}

	// Resolve the locale id for backwards compatibility.
	resolvedLocaleId := utils.ResolveLocaleId(localeIDParam)
	if resolvedLocaleId == nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.LocaleNotFound, "Locale not found.")
	}

	locales, err := services.GetLocaleLookup(*resolvedLocaleId, name)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	response := responses.LocaleLookupList{}
	response.SetLocaleLookupList(locales)

	return c.Status(fiber.StatusOK).JSON(response)
}
