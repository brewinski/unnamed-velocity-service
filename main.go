package main

import (
	router "github.com/brewinski/unnamed-fiber/pkg/router"
	"github.com/brewinski/unnamed-fiber/platform/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// create a new Fiber instance
	app := fiber.New()
	// connect to the database
	// database.ConnectDB()
	database.ConnectSqliteDB()
	// setup routes
	router.SetupRoutes(app)
	// Listen on PORT 3000
	app.Listen(":5001")
}
