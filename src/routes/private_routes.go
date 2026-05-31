package routes

import (
	"api-i18n/main/src/controllers"

	"github.com/ArnoldPMolenaar/api-utils/middleware"
	"github.com/gofiber/fiber/v3"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create private routes group.
	route := a.Group("/v1")

	// Register route group for /v1/apps.
	apps := route.Group("/apps")
	apps.Post("/", middleware.MachineProtected(), controllers.CreateApp)
	apps.Patch("/:name/locales", middleware.MachineProtected(), controllers.SetAppLocales)

	// Register route group for /v1/categories.
	categories := route.Group("/categories", middleware.MachineProtected())
	categories.Get("/", controllers.GetCategories)
	categories.Post("/", controllers.CreateCategory)
	categories.Get("/lookup", controllers.GetCategoryLookup)
	categories.Get("/:id", controllers.GetCategoryByID)
	categories.Patch("/:id", controllers.UpdateCategory)
	categories.Delete("/:id", controllers.DeleteCategory)
	categories.Post("/:id/restore", controllers.RestoreCategory)

	// Register route group for /v1/keys.
	keys := route.Group("/keys", middleware.MachineProtected())
	keys.Get("/", controllers.GetKeys)
	keys.Post("/", controllers.CreateKey)
	keys.Get("/:id", controllers.GetKeyByID)
	keys.Patch("/:id", controllers.UpdateKey)
	keys.Delete("/:id", controllers.DeleteKey)
	keys.Post("/:id/restore", controllers.RestoreKey)
}
