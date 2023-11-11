package url

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/indragunawan95/topedcrawler/internal/entity"
	"gorm.io/gorm"
)

type UrlRepo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *UrlRepo {
	return &UrlRepo{
		db: db,
	}
}

func (ur UrlRepo) CreateUrl(ctx context.Context, input entity.Url) (entity.Url, error) {
	input.ID = uuid.New().String()
	model := input.ToModel()

	err := ur.db.WithContext(ctx).Create(&model).Error

	if err != nil {
		return entity.Url{}, err
	}
	output := model.ToEntity()
	return output, nil
}

func (ur *UrlRepo) CreateUrls(ctx context.Context, inputs []entity.Url) ([]entity.Url, error) {
	// Convert the slice of Url entities to a slice of Url models
	var models []entity.UrlModel
	for _, input := range inputs {
		input.ID = uuid.New().String() // Assign a new UUID for each Url
		models = append(models, input.ToModel())
	}

	// Perform bulk insert
	err := ur.db.WithContext(ctx).Create(&models).Error
	if err != nil {
		return nil, err
	}

	// Convert back to slice of Url entities
	var output []entity.Url
	for _, model := range models {
		output = append(output, model.ToEntity())
	}

	return output, nil
}

func (ur UrlRepo) GetUrls(ctx context.Context) ([]entity.Url, error) {
	var models []entity.UrlModel

	err := ur.db.WithContext(ctx).Where("is_scrapped = ?", false).Find(&models).Error
	if err != nil {
		return nil, err
	}

	var output []entity.Url
	for _, model := range models {
		output = append(output, model.ToEntity()) // Convert each model back to entity
	}

	return output, nil
}

func (ur UrlRepo) MarkUrlAsScrapped(ctx context.Context, urlID string) error {
	result := ur.db.WithContext(ctx).Model(&entity.UrlModel{}).Where("id = ?", urlID).Update("is_scrapped", true)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no rows affected, check if the URL ID exists")
	}

	return nil
}
