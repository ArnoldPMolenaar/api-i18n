package controllers

import (
	"api-i18n/main/src/dto/responses"
	"api-i18n/main/src/enums"
	"api-i18n/main/src/errors"
	"api-i18n/main/src/services"
	"api-i18n/main/src/utils"

	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/gofiber/fiber/v2"
)

// GetTerritoryLookup func for getting territory lookup by locale ID, type and optional name filter.
func GetTerritoryLookup(c *fiber.Ctx) error {
	localeIDParam := c.Query("localeId")
	if localeIDParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "localeId query parameter is required.")
	}

	nameParam := c.Query("name")
	var name *string
	if nameParam != "" {
		name = &nameParam
	}

	typeParam := c.Query("type")
	var t *enums.TerritoryType
	if typeParam != "" {
		var tt enums.TerritoryType
		tt.Convert(typeParam)
		t = &tt
	}

	// Resolve the locale id for backwards compatibility.
	resolvedLocaleId := utils.ResolveLocaleId(localeIDParam)
	if resolvedLocaleId == nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.LocaleNotFound, "Locale not found.")
	}

	territories, err := services.GetTerritoriesLookup(*resolvedLocaleId, t, name)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	response := responses.TerritoryLookupList{}
	response.SetTerritoryLookupList(territories)

	return c.Status(fiber.StatusOK).JSON(response)
}
