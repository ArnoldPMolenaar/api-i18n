package controllers

import (
	"api-i18n/main/src/dto/requests"
	"api-i18n/main/src/dto/responses"
	"api-i18n/main/src/errors"
	"api-i18n/main/src/services"

	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	util "github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
)

// GetCategories func for getting all categories paginated.
func GetCategories(c *fiber.Ctx) error {
	paginationModel, err := services.GetCategories(c)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(paginationModel)
}

// GetCategoryLookup func for getting category lookup.
func GetCategoryLookup(c *fiber.Ctx) error {
	nameParam := c.Query("name")
	var name *string
	if nameParam != "" {
		name = &nameParam
	}

	categories, err := services.GetCategoryLookup(name)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	response := responses.CategoryLookupList{}
	response.SetCategoryLookupList(categories)

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetCategoryByID func for getting a category by ID.
func GetCategoryByID(c *fiber.Ctx) error {
	categoryIDParam := c.Params("id")
	categoryID, err := util.StringToUint(categoryIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	category, err := services.GetCategoryByID(categoryID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if category.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.CategoryExists, "Category does not exist.")
	}

	response := responses.Category{}
	response.SetCategory(category)

	return c.Status(fiber.StatusOK).JSON(response)
}

// CreateCategory func for creating a category.
func CreateCategory(c *fiber.Ctx) error {
	// Create a new category struct for the request.
	categoryRequest := &requests.CreateCategory{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(categoryRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate category fields.
	validate := util.NewValidator()
	if err := validate.Struct(categoryRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	// Check if category exists.
	if available, err := services.IsCategoryAvailable(categoryRequest.Name, nil); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.CategoryAvailable, "Category name already exist.")
	}

	// Check if category name exists as key name.
	if available, err := services.IsKeyAvailableGlobal(categoryRequest.Name, nil); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.CategoryIsKey, "Category name is a key name.")
	}

	// Create category.
	category, err := services.CreateCategory(categoryRequest.Name, categoryRequest.DisabledAt)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the category.
	response := responses.Category{}
	response.SetCategory(category)

	return c.Status(fiber.StatusCreated).JSON(response)
}

// UpdateCategory func for updating a category.
func UpdateCategory(c *fiber.Ctx) error {
	// Get the categoryID parameter from the URL.
	categoryIDParam := c.Params("id")
	categoryID, err := util.StringToUint(categoryIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Create a new category struct for the request.
	categoryRequest := &requests.UpdateCategory{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(categoryRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate category fields.
	validate := util.NewValidator()
	if err := validate.Struct(categoryRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	// Get old category.
	oldCategory, err := services.GetCategoryByID(categoryID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if oldCategory.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.CategoryExists, "Category does not exist.")
	}

	// Check if the category has been modified since it was last fetched.
	if categoryRequest.UpdatedAt.Unix() < oldCategory.UpdatedAt.Unix() {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.OutOfSync, "Data is out of sync.")
	}

	// Check if category exists.
	if categoryRequest.Name != oldCategory.Name {
		if available, err := services.IsCategoryAvailable(categoryRequest.Name, &oldCategory.Name); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
		} else if !available {
			return errorutil.Response(c, fiber.StatusBadRequest, errors.CategoryAvailable, "Category name already exist.")
		}

		// Check if category name exists as key name.
		if available, err := services.IsKeyAvailableGlobal(categoryRequest.Name, nil); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
		} else if !available {
			return errorutil.Response(c, fiber.StatusBadRequest, errors.CategoryIsKey, "Category name is a key name.")
		}
	}

	// Update category.
	updatedCategory, err := services.UpdateCategory(*oldCategory, categoryRequest.Name, categoryRequest.DisabledAt)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the category.
	response := responses.Category{}
	response.SetCategory(updatedCategory)

	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteCategory func for deleting a category.
func DeleteCategory(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := util.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the Category.
	category, err := services.GetCategoryByID(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if category.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.CategoryExists, "Category does not exist.")
	}

	// Delete the Category.
	if err := services.DeleteCategory(category.ID); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RestoreCategory func for restoring a deleted Category.
func RestoreCategory(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := util.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Check if category is deleted.
	if isDeleted, err := services.IsCategoryDeleted(id); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !isDeleted {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.CategoryAvailable, "Category is not deleted.")
	}

	// Restore the Category.
	if err := services.RestoreCategory(id); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
