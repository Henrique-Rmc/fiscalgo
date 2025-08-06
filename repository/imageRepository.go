package repository

import (
	"context"
	"errors"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"gorm.io/gorm"
)

type ImageRepositoryInterface interface {
	CreateImage(ctx context.Context, image *model.Image) error
	FindByUniqueFileName(ctx context.Context, uniqueName string) (*model.Image, error)
}

type ImageRepository struct {
	DB *gorm.DB
}

func NewImageRepo(db *gorm.DB) ImageRepositoryInterface {
	return &ImageRepository{DB: db}
}

func (imageRepo *ImageRepository) CreateImage(ctx context.Context, image *model.Image) error {
	if err := imageRepo.DB.Create(image).Error; err != nil {
		return err
	}
	return nil
}
func (imageRepo *ImageRepository) FindByUniqueFileName(ctx context.Context, uniqueName string) (*model.Image, error) {
	var image model.Image

	err := imageRepo.DB.WithContext(ctx).Where("uniqueName = ?", uniqueName).First(&image).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &image, nil
}
