package handlers

import (
	"github.com/brewinski/unnamed-fiber/velocity"
	"github.com/gofiber/fiber/v2"
)

func PointsEarnHandler(c *fiber.Ctx) error {
	bodyParams := velocity.VelocityEarPointsParams{}
	err := c.BodyParser(&bodyParams)

	if err != nil {
		return fiber.ErrBadRequest
	}
	res, err := velocity.AllocatePoint(bodyParams)

	if err != nil {
		return err
	}
	return c.JSON(res)
}
