package repository

import (
	"fiscalgo/model"

	"gorm.io/gorm"
)

// Criamos uma interface para que os objetos externos possam interagir com a interface por meio dela
type ImageRepositoryInterface interface {
	Create(image *model.Image) error
}

type ImageRepository struct {
	DB *gorm.DB
}

func NewImageRepo(db *gorm.DB) ImageRepositoryInterface {
	return &ImageRepository{DB: db}
}

func (r *ImageRepository) Create(image *model.Image) error {
	if err := r.DB.Create(image).Error; err != nil {
		return err
	}
	return nil
}
