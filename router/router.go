package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// SetupRoutes func to setup routes
func SetupRoutes(app *fiber.App) {
	health := app.Group("", logger.New(), recover.New(recover.Config{
		EnableStackTrace: true,
	})) // Group endpoints with param 'health'

	health.Get("/", func(c *fiber.Ctx) error {
		println("Health check")
		return c.SendString("OK")
	})

	api := app.Group("api", logger.New(), recover.New(recover.Config{
		EnableStackTrace: true,
	})) // Group endpoints with param 'api' and log whenever this endpoint is hit.

	api_v1 := api.Group("v1") // Group endpoints with param 'v1'
	setupNotesRoutes(api_v1)
	setupUserRoutes(api_v1)
	setupVelocityRoutes(api_v1)

	api_v2 := api.Group("v2") // Group endpoints with param 'v2'
	setupNotesRoutes(api_v2)
}
