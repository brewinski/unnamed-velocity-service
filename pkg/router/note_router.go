package router

import (
	"github.com/brewinski/unnamed-fiber/data"
	"github.com/brewinski/unnamed-fiber/internal/handler"
	"github.com/brewinski/unnamed-fiber/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func setupNotesRoutes(router fiber.Router) {
	note := router.Group("note")

	note.Post("",
		middleware.ValidationHandlerFactory(data.CreateNoteBody{}),
		handler.CreateNoteHandler)
}
