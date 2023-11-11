package main

import (
	"context"
	"fmt"
	"log"

	"github.com/indragunawan95/topedcrawler/files/config"
	"github.com/indragunawan95/topedcrawler/internal/entity"
	csvRepo "github.com/indragunawan95/topedcrawler/internal/repo/csv"
	productRepo "github.com/indragunawan95/topedcrawler/internal/repo/product"
	scrapperRepo "github.com/indragunawan95/topedcrawler/internal/repo/scrapper"
	urlRepo "github.com/indragunawan95/topedcrawler/internal/repo/url"
	scrapperUsecase "github.com/indragunawan95/topedcrawler/internal/usecase/scrappermanager"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/playwright-community/playwright-go"
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

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true), // Set to false to run in non-headless mode
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	numWorkers := cfg.App.NumWorkers
	numProducts := cfg.App.NumProducts

	productRepo := productRepo.New(db)
	urlRepo := urlRepo.New(db)
	scrapperRepo := scrapperRepo.New(browser)
	csvRepo := csvRepo.New("data.csv")

	scrapperUc := scrapperUsecase.New(productRepo, urlRepo, scrapperRepo, csvRepo, numWorkers)
	// Get Seed Url
	err = scrapperUc.GetAllProductLinks(context.Background(), numProducts)
	if err != nil {
		log.Fatalf("Error scraping product links: %v", err)
	}

	// Call the scrapper to process product details.
	err = scrapperUc.ProductDetailsScrapper()
	if err != nil {
		log.Fatalf("Error scraping product details: %v", err)
	}
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
