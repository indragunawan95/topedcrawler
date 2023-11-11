package csv

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"

	"github.com/indragunawan95/topedcrawler/internal/entity"
)

type CSVRepository struct {
	filePath string
}

func New(filePath string) *CSVRepository {
	return &CSVRepository{filePath: filePath}
}

func (r *CSVRepository) SaveProductsToCSV(ctx context.Context, products []entity.Product) error {
	// Open the file with append mode and write permissions
	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Check if the file is new to write the header
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	// If the file is new or empty, write the header
	if fileInfo.Size() == 0 {
		header := []string{"Name", "Description", "StoreName", "Price", "Rating", "ImageLink"}
		if err := writer.Write(header); err != nil {
			return err
		}
	}

	// Writing the product data
	for _, product := range products {
		record := []string{
			product.Name,
			product.Description,
			product.StoreName,
			product.Price,
			strconv.FormatFloat(float64(product.Rating), 'f', 2, 32),
			product.ImageLink,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
