package scrappermanager

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/indragunawan95/topedcrawler/internal/entity"
)

const (
	BaseURL = "https://www.tokopedia.com/p/handphone-tablet/handphone?ob=23&page=%d"
)

type ProductRepoItf interface {
	CreateProduct(ctx context.Context, input entity.Product) (entity.Product, error)
}

type CSVRepoItf interface {
	SaveProductsToCSV(ctx context.Context, products []entity.Product) error
}

type UrlRepoItf interface {
	CreateUrls(ctx context.Context, inputs []entity.Url) ([]entity.Url, error)
	GetUrls(ctx context.Context) ([]entity.Url, error)
	MarkUrlAsScrapped(ctx context.Context, urlID string) error
}

type ScrapperRepoItf interface {
	LaunchTab() error
	OpenPage(url string) error
	ScrollPage() error
	ClosePage() error
	GetProductTitle() (string, error)
	GetProductDescription() (string, error)
	GetProductStoreName() (string, error)
	GetProductPrice() (string, error)
	GetProductRating() (string, error)
	GetProductImageLink() (string, error)
	GetAllProductLinks() ([]string, error)
}

type Usecase struct {
	scrapperRepo ScrapperRepoItf
	productRepo  ProductRepoItf
	urlRepo      UrlRepoItf
	csvRepo      CSVRepoItf
	NumWorkers   int
}

func New(productRepo ProductRepoItf, urlRepo UrlRepoItf, scrapperRepo ScrapperRepoItf, csvRepo CSVRepoItf, numWorkers int) *Usecase {

	return &Usecase{
		productRepo:  productRepo,
		urlRepo:      urlRepo,
		scrapperRepo: scrapperRepo,
		csvRepo:      csvRepo,
		NumWorkers:   numWorkers,
	}
}

// Get all seed product link first
func (uc *Usecase) GetAllProductLinks(ctx context.Context, maxLinks int) error {
	if err := uc.scrapperRepo.LaunchTab(); err != nil {
		return fmt.Errorf("failed to launch tab: %w", err)
	}

	urls := make([]entity.Url, 0, maxLinks)
	pageIndex := 1

	for len(urls) < maxLinks {
		pageURL := fmt.Sprintf(BaseURL, pageIndex)
		if err := uc.scrapperRepo.OpenPage(pageURL); err != nil {
			return fmt.Errorf("failed to open page: %w", err)
		}

		if err := uc.scrapperRepo.ScrollPage(); err != nil {
			return fmt.Errorf("failed to scroll page: %w", err)
		}

		links, err := uc.scrapperRepo.GetAllProductLinks()
		if err != nil {
			return fmt.Errorf("failed to scrape product links: %w", err)
		}

		for _, link := range links {
			if len(urls) < maxLinks {
				urls = append(urls, entity.Url{Url: link})
			} else {
				break // We have reached the maxLinks limit
			}
		}

		pageIndex++ // Move to the next page
	}

	if _, err := uc.urlRepo.CreateUrls(ctx, urls); err != nil {
		return fmt.Errorf("failed to save product links: %w", err)
	}

	return nil
}

// Scrap product detail from seed product link
func (uc *Usecase) ProductDetailsScrapper() error {
	urls, err := uc.urlRepo.GetUrls(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get URLs: %w", err)
	}

	// Create a channel to send URLs to be processed.
	urlsChan := make(chan entity.Url)
	// Create a channel to communicate errors from goroutines.
	errChan := make(chan error, uc.NumWorkers)
	// WaitGroup to wait for all goroutines to finish.
	var wg sync.WaitGroup

	// Start the specified number of worker goroutines.
	for i := 0; i < uc.NumWorkers; i++ {
		wg.Add(1)
		go worker(&wg, urlsChan, errChan, uc)
	}

	// Send URLs to the channel for the workers to process.
	go func() {
		for _, urlEntity := range urls {
			urlsChan <- urlEntity
		}
		close(urlsChan) // Close the channel to signal workers to stop.
	}()

	// Wait for all goroutines to complete and close the error channel.
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collect errors, if any.
	var scrappingErrors []error
	for e := range errChan {
		if e != nil {
			scrappingErrors = append(scrappingErrors, e)
		}
	}

	if len(scrappingErrors) > 0 {
		// Handle errors accordingly. For example, you could log them or retry failed operations.
		// For now, just returning the first error.
		return scrappingErrors[0]
	}

	return nil
}

// Worker function that processes URLs from the urlsChan and sends errors to errChan.
func worker(wg *sync.WaitGroup, urlsChan <-chan entity.Url, errChan chan<- error, uc *Usecase) {
	defer wg.Done()
	for url := range urlsChan {
		if err := uc.processUrl(url); err != nil {
			errChan <- err
			// For now, let's just log the error and move on to the next URL.
			log.Printf("Error processing URL %s: %v", url.Url, err)
		}
	}
}

func (uc *Usecase) processUrl(url entity.Url) error {
	if err := uc.scrapperRepo.LaunchTab(); err != nil {
		return fmt.Errorf("failed to launch tab: %w", err)
	}
	defer uc.scrapperRepo.ClosePage()

	if err := uc.scrapperRepo.OpenPage(url.Url); err != nil {
		return fmt.Errorf("failed to open product detail page: %w", err)
	}

	if err := uc.scrapperRepo.ScrollPage(); err != nil {
		return fmt.Errorf("failed to scroll page: %w", err)
	}

	product, err := uc.scrapeProductDetails()
	if err != nil {
		return fmt.Errorf("failed to scrape product details: %w", err)
	}

	_, err = uc.productRepo.CreateProduct(context.Background(), product)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	err = uc.urlRepo.MarkUrlAsScrapped(context.Background(), url.ID)
	if err != nil {
		return fmt.Errorf("failed to update scrapped: %w", err)
	}

	// append the product details to the CSV file.
	err = uc.csvRepo.SaveProductsToCSV(context.Background(), []entity.Product{product})
	if err != nil {
		return fmt.Errorf("failed to save product to CSV: %w", err)
	}

	log.Printf("Processed product: %s\n", product.Name)
	return nil
}

func (uc *Usecase) scrapeProductDetails() (entity.Product, error) {
	var product entity.Product

	title, err := uc.scrapperRepo.GetProductTitle()
	if err != nil {
		return product, fmt.Errorf("failed to get product title: %w", err)
	}

	description, err := uc.scrapperRepo.GetProductDescription()
	if err != nil {
		return product, fmt.Errorf("failed to get product description: %w", err)
	}

	storeName, err := uc.scrapperRepo.GetProductStoreName()
	if err != nil {
		return product, fmt.Errorf("failed to get product storename: %w", err)
	}

	price, err := uc.scrapperRepo.GetProductPrice()
	if err != nil {
		return product, fmt.Errorf("failed to get product price: %w", err)
	}

	rating, err := uc.scrapperRepo.GetProductRating()
	if err != nil {
		return product, fmt.Errorf("failed to get product rating: %w", err)
	}

	ratingFloat, err := strconv.ParseFloat(rating, 32) // 32 specifies the precision
	if err != nil {
		return product, fmt.Errorf("error converting string to float: %w", err)
	}

	imageLink, err := uc.scrapperRepo.GetProductImageLink()
	if err != nil {
		return product, fmt.Errorf("failed to get product main image url: %w", err)
	}

	product = entity.Product{
		Name:        title,
		Description: description,
		StoreName:   storeName,
		Price:       extractPrice(price),
		Rating:      float32(ratingFloat),
		ImageLink:   imageLink,
	}

	return product, nil
}

func extractPrice(priceWithCurrency string) string {
	// Remove all non-digit characters
	return strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, priceWithCurrency)
}
