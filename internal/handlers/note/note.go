package noteHandlers

import (
	"github.com/brewinski/unnamed-fiber/internal/model"
	"github.com/gofiber/fiber/v2"
)

func CreateNotes(c *fiber.Ctx) error {
	// db := database.DB
	note := model.Note{}

	// store the body of the note in the note variable
	err := c.BodyParser(&note)

	if err != nil {
		return fiber.ErrBadRequest
	}

	return c.JSON(note)
}
