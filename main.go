package main

import (
	"github.com/brewinski/unnamed-fiber/db"
	router "github.com/brewinski/unnamed-fiber/router"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// create a new Fiber instance
	app := fiber.New()
	// connect to the database
	db.ConnectDB()
	// setup routes
	router.SetupRoutes(app)
	// Listen on PORT 3000
	app.Listen(":5001")
}
