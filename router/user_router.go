package router

import (
	"github.com/brewinski/unnamed-fiber/handlers"
	"github.com/gofiber/fiber/v2"
)

func setupUserRoutes(router fiber.Router) {
	user := router.Group("user")

	user.Get("",
		handlers.ListUsersHandler)
}
