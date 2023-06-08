package handler

import (
	"github.com/brewinski/unnamed-fiber/data"
	"github.com/brewinski/unnamed-fiber/platform/database"
	"github.com/gofiber/fiber/v2"
)

// This handler should abstract saving a note from the fiber fameowrk implementation details.
// Keep logic portable by extracting the request values we need and passing them to the worker functions.

func CreateNoteHandler(c *fiber.Ctx) error {
	// db := database.DB
	note := data.Note{}

	// store the body of the note in the note variable
	err := c.BodyParser(&note)

	if err != nil {
		return fiber.ErrBadRequest
	}

	createdNote, err := CreateNote(note)

	if err != nil {
		return err
	}

	readNote, err := readNote(createdNote.ID)

	if err != nil {
		return err
	}

	return c.JSON(readNote)
}

func CreateNote(note data.Note) (data.Note, error) {
	db := database.DB

	err := db.Create(&note).Error

	if err != nil {
		return note, err
	}

	return note, nil
}

func readNote(id int) (data.Note, error) {
	db := database.DB
	note := data.Note{}

	err := db.First(&note, id).Error

	if err != nil {
		return note, err
	}

	return note, nil
}
