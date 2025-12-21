package routes

import (
	"api-i18n/main/src/controllers"

	"github.com/ArnoldPMolenaar/api-utils/middleware"
	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create private routes group.
	route := a.Group("/v1")

	// Register route group for /v1/apps.
	apps := route.Group("/apps", middleware.MachineProtected())
	apps.Post("/", controllers.CreateApp)
	apps.Get("/:name/locales", controllers.GetAppLocales)
	apps.Put("/:name/locales", controllers.SetAppLocales)

	// Register route group for /v1/categories.
	categories := route.Group("/categories", middleware.MachineProtected())
	categories.Get("/", controllers.GetCategories)
	categories.Post("/", controllers.CreateCategory)
	categories.Get("/lookup", controllers.GetCategoryLookup)
	categories.Get("/:id", controllers.GetCategoryByID)
	categories.Put("/:id", controllers.UpdateCategory)
	categories.Delete("/:id", controllers.DeleteCategory)
	categories.Put("/:id/restore", controllers.RestoreCategory)

	// Register route group for /v1/keys.
	keys := route.Group("/keys", middleware.MachineProtected())
	keys.Get("/", controllers.GetKeys)
	keys.Post("/", controllers.CreateKey)
	keys.Get("/:id", controllers.GetKeyByID)
	keys.Put("/:id", controllers.UpdateKey)
	keys.Delete("/:id", controllers.DeleteKey)
	keys.Put("/:id/restore", controllers.RestoreKey)
}
