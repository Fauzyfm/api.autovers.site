package config

import (
	"log"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"fmt"
)

var DB *gorm.DB

func InitDatabase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	fmt.Println("Connecting to DB:", os.Getenv("DB_HOST"), os.Getenv("DB_NAME")) // debug

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db
}