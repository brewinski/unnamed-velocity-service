package handler

import (
	"github.com/brewinski/unnamed-fiber/internal/model"
	"github.com/gofiber/fiber/v2"
)

// This handler should abstract saving a note from the fiber fameowrk implementation details.
// Keep logic portable by extracting the request values we need and passing them to the worker functions.

func CreateNoteHandler(c *fiber.Ctx) error {
	// db := database.DB
	note := model.CreateNoteBody{}

	// store the body of the note in the note variable
	err := c.BodyParser(&note)
	if err != nil {
		return fiber.ErrBadRequest
	}

	CreateNote(&note)

	return c.JSON(note)
}

func CreateNote(note *model.CreateNoteBody) error {

	return nil
}
