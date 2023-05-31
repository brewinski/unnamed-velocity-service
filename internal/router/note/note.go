package noteRoutes

import (
	noteHandlers "github.com/brewinski/unnamed-fiber/internal/handlers/note"
	"github.com/gofiber/fiber/v2"
)

func SetupNotesRoutes(router fiber.Router) {
	note := router.Group("note")

	note.Post("", noteHandlers.CreateNotes)
}
