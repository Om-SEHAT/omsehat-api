package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Om-SEHAT/omsehat-api/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	// get database connection string from .env file
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// construct DSN
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	// open connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to database successfully")

	// migrate models to databbase
	err = db.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.Queue{},
		&models.Doctor{},
		&models.Message{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migrated successfully")

	// set db to global variable
	DB = db
}
