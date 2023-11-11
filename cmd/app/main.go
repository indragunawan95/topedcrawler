package main

import (
	"context"
	"fmt"
	"log"

	"github.com/indragunawan95/topedcrawler/files/config"
	"github.com/indragunawan95/topedcrawler/internal/entity"
	urlRepo "github.com/indragunawan95/topedcrawler/internal/repo/url"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error getting information varibale: %v", err)
	}

	db, err := dbSetup(cfg)
	if err != nil {
		log.Fatal("Error:", err)
	}

	// Try add URL
	// prodRepo := productRepo.New(db)
	// productUsecase(prodRepo)
	repo := urlRepo.New(db)
	url, err := repo.CreateUrl(context.Background(), entity.Url{
		Url: "https://www.tokopedia.com/wahanaacesories/hp-redmi-note-9-smartphone-xiaomi-mi-note-9-ram-4-64gb-garansi-resmi-hijau",
	})
	if err != nil {
		log.Fatal("Error:", err)
	}
	fmt.Println(url)

}

func dbSetup(cfg *config.Config) (*gorm.DB, error) {
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
		return nil, err
	}

	// Ensure uuid-ossp extension is available
	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		log.Fatal("failed to create uuid extension:", err)
		return nil, err
	}

	// Automigrate your models
	err = db.AutoMigrate(&entity.ProductModel{}, &entity.UrlModel{})
	if err != nil {
		log.Fatal("failed to migrate:", err)
		return nil, err
	}

	return db, nil
}
