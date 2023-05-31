package main

import (
	"github.com/brewinski/unnamed-fiber/database"
	"github.com/brewinski/unnamed-fiber/internal/router"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// create a new Fiber instance
	app := fiber.New()
	// connect to the database
	database.ConnectDB()
	// setup routes
	router.SetupRoutes(app)
	// Listen on PORT 3000
	app.Listen(":5001")
}
