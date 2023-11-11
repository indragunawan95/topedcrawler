package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Used across layer except Persistence
type Product struct {
	ID          string
	Name        string
	Description string
	ImageLink   string
	Price       string
	Rating      float32
	StoreName   string
}

func (p Product) ToModel() ProductModel {
	return ProductModel{
		ID:          uuid.MustParse(p.ID),
		Name:        p.Name,
		Description: p.Description,
		ImageLink:   p.ImageLink,
		Price:       p.Price,
		Rating:      p.Rating,
		StoreName:   p.StoreName,
	}
}

// Used in by Gorm
type ProductModel struct {
	gorm.Model            // Embeds fields like ID, CreatedAt, UpdatedAt, DeletedAt
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:text;not null"`
	ImageLink   string    `gorm:"type:text;not null"`
	Price       string    `gorm:"type:varchar(100);not null"`
	Rating      float32   `gorm:"type:decimal(10,2)"`
	StoreName   string    `gorm:"type:varchar(100);not null"`
}

// TableName overrides the table name used by ProductModel to `products`
func (ProductModel) TableName() string {
	return "products"
}

// ToDomain converts the persistence model to the domain entity
func (p ProductModel) ToEntity() Product {
	return Product{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		ImageLink:   p.ImageLink,
		Price:       p.Price,
		Rating:      p.Rating,
		StoreName:   p.StoreName,
	}
}
