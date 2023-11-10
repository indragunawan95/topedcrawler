package product

import (
	"context"

	"github.com/indragunawan95/topedcrawler/internal/entity"
)

type ProductRepoItf interface {
	CreateUser(ctx context.Context, input entity.Product) (entity.Product, error)
}

type Usecase struct {
	product ProductRepoItf
}

func New(product ProductRepoItf) *Usecase {
	return &Usecase{
		product: product,
	}
}

func (uc Usecase) CreateProduct(ctx context.Context, input entity.Product) (entity.Product, error) {
	newProduct, err := uc.product.CreateUser(ctx, input)

	if err != nil {
		return entity.Product{}, err
	}

	return newProduct, nil
}
