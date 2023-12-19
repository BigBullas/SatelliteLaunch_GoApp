package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"RIP_lab1/internal/api/pkg"
	"RIP_lab1/internal/models"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(pkg.GetConnectionString()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Product{})
	if err != nil {
		panic("cant migrate db")
	}
}