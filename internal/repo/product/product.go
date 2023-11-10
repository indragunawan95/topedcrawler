package product

import (
	"context"

	"github.com/indragunawan95/topedcrawler/internal/entity"
	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

func (pr ProductRepo) CreateProduct(ctx context.Context, input entity.Product) (entity.Product, error) {
	model := input.ToModel()

	err := pr.db.WithContext(ctx).Create(&model).Error

	if err != nil {
		return entity.Product{}, err
	}
	output := model.ToEntity()
	return output, nil
}
