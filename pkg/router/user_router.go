package router

import (
	"github.com/brewinski/unnamed-fiber/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func setupUserRoutes(router fiber.Router) {
	user := router.Group("user")

	user.Get("",
		handler.ListUsersHandler)
}
