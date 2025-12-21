package controllers

import (
	"api-i18n/main/src/dto/responses"
	"api-i18n/main/src/errors"
	"api-i18n/main/src/services"
	"api-i18n/main/src/utils"

	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/gofiber/fiber/v2"
)

// GetPhoneLookup handles the phone code lookup request.
func GetPhoneLookup(c *fiber.Ctx) error {
	localeIDParam := c.Query("localeId")
	if localeIDParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "localeId query parameter is required.")
	}

	// Resolve the locale id for backwards compatibility.
	resolvedLocaleId := utils.ResolveLocaleId(localeIDParam)
	if resolvedLocaleId == nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.LocaleNotFound, "Locale not found.")
	}

	phoneCodes, err := services.GetTerritoryPhoneCodes(*resolvedLocaleId)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.InternalServerError, err.Error())
	}

	response := responses.PhoneCodeLookupList{}
	response.SetPhoneCodeLookupList(phoneCodes)

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetPhoneNumberValidation handles the phone number validation request.
func GetPhoneNumberValidation(c *fiber.Ctx) error {
	localeIDParam := c.Query("localeId")
	phoneNumberParam := c.Query("phoneNumber")

	if localeIDParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "localeId query parameter is required.")
	}
	if phoneNumberParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "phoneNumber query parameter is required.")
	}

	// Resolve the locale id for backwards compatibility.
	resolvedLocaleId := utils.ResolveLocaleId(localeIDParam)
	if resolvedLocaleId == nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.LocaleNotFound, "Locale not found.")
	}

	isValid, err := services.ValidatePhoneNumber(phoneNumberParam, resolvedLocaleId)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.InternalServerError, err.Error())
	}

	response := &responses.PhoneNumberValid{}
	response.IsValid = isValid

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetPhoneNumberFormat handles the phone number format request.
// To format a phone number according to the specified locale.
func GetPhoneNumberFormat(c *fiber.Ctx) error {
	localeIDParam := c.Query("localeId")
	phoneNumberParam := c.Query("phoneNumber")

	if localeIDParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "localeId query parameter is required.")
	}
	if phoneNumberParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "phoneNumber query parameter is required.")
	}

	// Resolve the locale id for backwards compatibility.
	resolvedLocaleId := utils.ResolveLocaleId(localeIDParam)
	if resolvedLocaleId == nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.LocaleNotFound, "Locale not found.")
	}

	phoneNumberFormat, err := services.FormatPhoneNumber(phoneNumberParam, resolvedLocaleId)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.InternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(phoneNumberFormat)
}
