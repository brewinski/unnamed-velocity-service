package router

import (
	noteRoutes "github.com/brewinski/unnamed-fiber/internal/router/note"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// SetupRoutes func to setup routes
func SetupRoutes(app *fiber.App) {

	api := app.Group("api", logger.New(), recover.New(recover.Config{
		EnableStackTrace: true,
	})) // Group endpoints with param 'api' and log whenever this endpoint is hit.

	noteRoutes.SetupNotesRoutes(api)

	user := api.Group("user")

	user.Get("/", func(c *fiber.Ctx) error {
		if true {
			return fiber.ErrBadRequest
		}
		return c.SendString("All users")
	})
}
