package middleware

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = validator.New()

func validateStruct(body interface{}) []*ErrorResponse {
	var errors []*ErrorResponse

	err := validate.Struct(body)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}

	return errors
}

func ValidationHandlerFactory(validationType interface{}) fiber.Handler {
	typeOfStruct := reflect.TypeOf(validationType)

	return func(c *fiber.Ctx) error {
		validatoinStruct := reflect.New(typeOfStruct).Interface()

		// store the body of the note in the note variable
		err := c.BodyParser(&validatoinStruct)
		if err != nil {
			return fiber.ErrBadRequest
		}

		errors := validateStruct(validatoinStruct)

		if errors != nil {
			return c.Status(fiber.ErrBadRequest.Code).JSON(errors)
		}

		return c.Next()
	}
}
