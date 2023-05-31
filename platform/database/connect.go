package database

import (
	"fmt"
	"log"
	"strconv"

	"github.com/brewinski/unnamed-fiber/internal/model"
	"github.com/brewinski/unnamed-fiber/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Declare the variable for the database
var DB *gorm.DB

// ConnectDB connect to db
func ConnectDB() {
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)

	if err != nil {
		log.Println("Idiot")
	}

	// Connection URL to connect to Postgres Database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Config("DB_HOST"), port, config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"))
	// Connect to the DB and initialize the DB variable
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")

	// Migrate the schema
	err = DB.AutoMigrate(&model.Note{})

	if err != nil {
		panic("failed to migrate database")
	}

	fmt.Println("Database Migrated")
}

func ConnectSqliteDB() {
	// Connect to the DB and initialize the DB variable
	DB, err := gorm.Open(sqlite.Open("fiber.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")

	// Migrate the schema
	err = DB.AutoMigrate(&model.Note{})

	if err != nil {
		panic("failed to migrate database")
	}

	fmt.Println("Database Migrated")
}