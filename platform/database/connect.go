package database

import (
	"fmt"
	"log"
	"strconv"

	"github.com/brewinski/unnamed-fiber/data"
	"github.com/brewinski/unnamed-fiber/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Declare the variable for the database
var (
	DB *gorm.DB
)

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
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")

	// Migrate the schema
	err = DB.AutoMigrate(&data.Note{}, &data.User{}, &data.Credit{})

	if err != nil {
		panic("failed to migrate database")
	}

	fmt.Println("Database Migrated")
}

func ConnectSqliteDB() {
	var err error
	// Connect to the DB and initialize the DB variable
	DB, err = gorm.Open(sqlite.Open("fiber.sqlite"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")

	// Migrate the schema
	err = DB.Debug().AutoMigrate(&data.Note{}, &data.User{}, &data.Credit{})

	if err != nil {
		panic("failed to migrate database")
	}

	fmt.Println("Database Migrated")
}
