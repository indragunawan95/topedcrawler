package scrapper

import (
	"context"
	"errors"

	"github.com/indragunawan95/topedcrawler/internal/entity"
	"github.com/playwright-community/playwright-go"
)

type UrlRepoItf interface {
	GetUrls(ctx context.Context) ([]entity.Url, error)
}

type ScrapperRepo struct {
	browser playwright.Browser
	page    playwright.Page
}

func New(browser playwright.Browser) *ScrapperRepo {
	return &ScrapperRepo{
		browser: browser,
	}
}

func (s *ScrapperRepo) LaunchTab() error {
	page, err := s.browser.NewPage()
	if err != nil {
		return err
	}
	s.page = page
	return nil
}

func (s *ScrapperRepo) OpenPage(url string) error {
	_, err := s.page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *ScrapperRepo) ScrollPage() error {
	// Scroll to the bottom of the page to trigger lazy loading.
	// This JavaScript snippet will scroll to the bottom.
	_, err := s.page.Evaluate("window.scrollTo(0, document.body.scrollHeight)")
	if err != nil {
		return err
	}

	// Wait for the network to be idle after scrolling.
	loadStateOptions := playwright.PageWaitForLoadStateOptions{
		State: (*playwright.LoadState)(playwright.WaitUntilStateLoad),
	}
	err = s.page.WaitForLoadState(loadStateOptions)
	if err != nil {
		return err
	}
	return nil
}

func (s *ScrapperRepo) ClosePage() error {
	return s.page.Close()
}

func (s *ScrapperRepo) GetProductTitle() (string, error) {
	selector := "[data-testid='lblPDPDetailProductName']"
	locator := s.page.Locator(selector)
	if locator == nil {
		return "", errors.New("locator not found")
	}
	title, err := locator.TextContent()
	if err != nil {
		return "", err
	}
	return title, nil
}

func (s *ScrapperRepo) GetProductDescription() (string, error) {
	selector := "[data-testid='lblPDPDescriptionProduk']"
	locator := s.page.Locator(selector)
	if locator == nil {
		return "", errors.New("locator not found")
	}
	description, err := locator.TextContent()
	if err != nil {
		return "", err
	}
	return description, nil
}

func (s *ScrapperRepo) GetProductStoreName() (string, error) {
	selector := "a[data-testid='llbPDPFooterShopName'] h2"
	locator := s.page.Locator(selector)
	if locator == nil {
		return "", errors.New("locator not found")
	}
	storeName, err := locator.TextContent()
	if err != nil {
		return "", err
	}
	return storeName, nil
}

func (s *ScrapperRepo) GetProductPrice() (string, error) {
	selector := "[data-testid='lblPDPDetailProductPrice']"
	locator := s.page.Locator(selector)
	if locator == nil {
		return "", errors.New("locator not found")
	}
	price, err := locator.TextContent()
	if err != nil {
		return "", err
	}
	return price, nil
}

func (s *ScrapperRepo) GetProductRating() (string, error) {
	selector := "[data-testid='lblPDPDetailProductRatingNumber']"
	locator := s.page.Locator(selector)
	// Attempt to get the text content of the locator
	rating, err := locator.TextContent()
	// If there is an error, return "0" as the default rating
	if err != nil {
		// You might want to log the error or handle different types of errors appropriately
		return "0", nil
	}
	// No error, return the rating
	return rating, nil
}

func (s *ScrapperRepo) GetProductImageLink() (string, error) {
	// Using data-testid to select the image element
	selector := "[data-testid='PDPMainImage']"
	locator := s.page.Locator(selector)
	if locator == nil {
		return "", errors.New("locator not found")
	}
	// Get the "src" attribute of the image element
	imageLink, err := locator.GetAttribute("src")
	if err != nil {
		return "", err
	}
	return imageLink, nil
}

func (s *ScrapperRepo) GetAllProductLinks() ([]string, error) {
	// Selector for all elements with the specific data-testid
	selector := "a[data-testid='lnkProductContainer']"
	// Create a locator for all elements matching the selector
	locator := s.page.Locator(selector)

	// Count the number of elements matched by the locator
	count, err := locator.Count()
	if err != nil {
		return nil, err
	}

	var links []string
	for i := 0; i < count; i++ {
		// Get the "href" attribute of each element by using the nth locator
		href, err := locator.Nth(i).GetAttribute("href")
		if err != nil {
			// Decide how to handle the error; skip to the next iteration or return
			continue
		}
		links = append(links, href)
	}
	return links, nil
}
