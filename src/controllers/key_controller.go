package controllers

import (
	"api-i18n/main/src/dto/requests"
	"api-i18n/main/src/dto/responses"
	"api-i18n/main/src/errors"
	"api-i18n/main/src/services"

	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	util "github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

// GetKeys func for getting all keys paginated.
func GetKeys(c *fiber.Ctx) error {
	paginationModel, err := services.GetKeys(c)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(paginationModel)
}

// GetKeyByID func for getting a key by ID.
func GetKeyByID(c *fiber.Ctx) error {
	keyIDParam := c.Params("id")
	keyID, err := util.StringToUint(keyIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	key, err := services.GetKeyByID(keyID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if key.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.KeyExists, "Key does not exist.")
	}

	response := responses.Key{}
	response.SetKey(key)

	return c.Status(fiber.StatusOK).JSON(response)
}

// CreateKey func for creating a key.
func CreateKey(c *fiber.Ctx) error {
	// Create a new key struct for the request.
	keyRequest := &requests.CreateKey{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(keyRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate key fields.
	validate := util.NewValidator()
	if err := validate.Struct(keyRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	// Check if the app exists.
	appAvailable, err := services.IsAppAvailable(keyRequest.AppName)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !appAvailable {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AppNotFound, "App not found.")
	}

	// Check if key exists.
	if available, err := services.IsKeyAvailable(keyRequest.AppName, keyRequest.Name, keyRequest.CategoryID, nil); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.KeyAvailable, "Key name already exist.")
	}

	// Check if key is not a category.
	if available, err := services.IsCategoryAvailable(keyRequest.Name, nil); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.KeyIsCategory, "key name is equal to a category name.")
	}

	// Check if key has valid translations.
	localeIds := lo.Map(keyRequest.Translations, func(t requests.CreateKeyTranslation, _ int) string { return t.LocaleID })
	if valid, err := services.HasValidTranslations(keyRequest.AppName, localeIds); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !valid {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.InvalidTranslations, "One or more translations are invalid.")
	}

	// Create key.
	key, err := services.CreateKey(*keyRequest)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	key, err = services.GetKeyByID(key.ID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if key.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.KeyExists, "Key does not exist.")
	}

	// Return the key.
	response := responses.Key{}
	response.SetKey(key)

	return c.Status(fiber.StatusCreated).JSON(response)
}

// UpdateKey func for updating a key.
func UpdateKey(c *fiber.Ctx) error {
	// Get the keyID parameter from the URL.
	keyIDParam := c.Params("id")
	keyID, err := util.StringToUint(keyIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Create a new key struct for the request.
	keyRequest := &requests.UpdateKey{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(keyRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate key fields.
	validate := util.NewValidator()
	if err := validate.Struct(keyRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	// Get old key.
	oldKey, err := services.GetKeyByID(keyID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if oldKey.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.KeyExists, "Key does not exist.")
	}

	// Check if the key has been modified since it was last fetched.
	if keyRequest.UpdatedAt.Unix() < oldKey.UpdatedAt.Unix() {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.OutOfSync, "Data is out of sync.")
	}
	translationMap := make(map[string]requests.UpdateKeyTranslation)
	for _, translation := range keyRequest.Translations {
		translationMap[translation.LocaleID] = translation
	}
	for _, translation := range oldKey.Translations {
		if _, exists := translationMap[translation.LocaleID]; exists {
			if translationMap[translation.LocaleID].UpdatedAt.Unix() < translation.UpdatedAt.Unix() {
				return errorutil.Response(c, fiber.StatusBadRequest, errorutil.OutOfSync, "One or more translations are out of sync.")
			}
		}
	}

	// Check if key exists.
	if keyRequest.Name != oldKey.Name {
		var ignore *string
		if keyRequest.CategoryID == &oldKey.CategoryID.V {
			ignore = &oldKey.Name
		}
		if available, err := services.IsKeyAvailable(oldKey.AppName, keyRequest.Name, keyRequest.CategoryID, ignore); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
		} else if !available {
			return errorutil.Response(c, fiber.StatusBadRequest, errors.KeyAvailable, "Key name already exist.")
		}

		// Check if key is not a category.
		if available, err := services.IsCategoryAvailable(keyRequest.Name, nil); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
		} else if !available {
			return errorutil.Response(c, fiber.StatusBadRequest, errors.KeyIsCategory, "key name is equal to a category name.")
		}
	}

	// Check if key has valid translations.
	localeIds := lo.Map(keyRequest.Translations, func(t requests.UpdateKeyTranslation, _ int) string { return t.LocaleID })
	if valid, err := services.HasValidTranslations(oldKey.AppName, localeIds); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !valid {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.InvalidTranslations, "One or more translations are invalid.")
	}

	// Update key.
	updatedKey, err := services.UpdateKey(*oldKey, *keyRequest)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	key, err := services.GetKeyByID(updatedKey.ID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if key.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.KeyExists, "Key does not exist.")
	}

	// Return the key.
	response := responses.Key{}
	response.SetKey(key)

	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteKey func for deleting a key.
func DeleteKey(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := util.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the Key.
	key, err := services.GetKeyByID(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if key.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.KeyExists, "Key does not exist.")
	}

	// Delete the Key.
	if err := services.DeleteKey(key.ID); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RestoreKey func for restoring a deleted Key.
func RestoreKey(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := util.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Check if key is deleted.
	if isDeleted, err := services.IsKeyDeleted(id); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !isDeleted {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.KeyAvailable, "Key is not deleted.")
	}

	// Restore the Key.
	if err := services.RestoreKey(id); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
