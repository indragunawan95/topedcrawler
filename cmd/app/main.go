package main

import (
	"fmt"
	"log"

	"github.com/indragunawan95/topedcrawler/files/config"    // Adjust to your project's structure
	"github.com/indragunawan95/topedcrawler/internal/entity" // Adjust to your project's structure
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.NewConfig()
	// Format the DSN using the configuration values
	datasource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Username,
		cfg.DB.Password, // Ensure you are using cfg.DB.Password instead of cfg.Password
		cfg.DB.Name,
	)

	// Open the database connection using the formatted DSN
	db, err := gorm.Open(postgres.Open(datasource), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	// Automigrate your models
	err = db.AutoMigrate(&entity.ProductModel{})
	if err != nil {
		log.Fatal("failed to migrate:", err)
	}
}
