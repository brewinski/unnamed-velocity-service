package noteRoutes

import (
	noteHandlers "github.com/brewinski/unnamed-fiber/internal/handlers/note"
	"github.com/brewinski/unnamed-fiber/internal/model"
	utilValidation "github.com/brewinski/unnamed-fiber/internal/util/validation"
	"github.com/gofiber/fiber/v2"
)

func SetupNotesRoutes(router fiber.Router) {
	note := router.Group("note")

	note.Post("", ValidateCreateNotes, noteHandlers.CreateNotes)
}

func ValidateCreateNotes(c *fiber.Ctx) error {
	note := model.NoteBody{}

	// store the body of the note in the note variable
	err := c.BodyParser(&note)
	if err != nil {
		return fiber.ErrBadRequest
	}

	errors := utilValidation.ValidateStruct(note)

	if errors != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(errors)
	}

	return c.Next()
}
