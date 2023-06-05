package model

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model // Adds some metadata fields to the table
	// ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4();unique;not null;"` // Explicitly specify the type to be uuid
	ID       int `gorm:"primaryKey;unique;not null;"` // Explicitly specify the type to be uuid
	Title    string
	SubTitle string
	Text     string
}

// transport layer
type CreateNoteBody struct {
	Title    string `validate:"required,min=1,max=50"`
	SubTitle string `validate:"required"`
	Text     string `validate:"required"`
}
