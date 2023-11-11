package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Url struct {
	ID         string
	Url        string
	IsScrapped bool
}

func (url Url) ToModel() UrlModel {
	return UrlModel{
		ID:         uuid.MustParse(url.ID),
		Url:        url.Url,
		IsScrapped: url.IsScrapped,
	}
}

type UrlModel struct {
	gorm.Model           // Embeds fields like ID, CreatedAt, UpdatedAt, DeletedAt
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Url        string    `gorm:"type:text;not null"`
	IsScrapped bool      `gorm:"type:boolean;not null"`
}

func (UrlModel) TableName() string {
	return "urls"
}

func (url UrlModel) ToEntity() Url {
	return Url{
		ID:         url.ID.String(),
		Url:        url.Url,
		IsScrapped: url.IsScrapped,
	}
}
