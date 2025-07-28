package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"gorm.io/gorm"
)

// Criamos uma interface para definir que quem implementa essa ineterface, deve possuir os métodos da interface
// Nesse caso ,
type ImageRepositoryInterface interface {
	CreateImage(ctx context.Context, image *model.Image) error
	FindByUniqueFileName(ctx context.Context, uniqueName string) error
}

// Na strutct ImageRepository estamos fazendo uma injeção de dependencia no objeto indicando que
// ImageRepository deve possuir uma instancia de um gorm.DB
type ImageRepository struct {
	DB *gorm.DB
}

// *
// Esse é o construtor do ImageRepository e ele indica que deve receber um db Gorm e retornar um Image
// Repository interface
// Como isso conecta os métodos do ImageRepositoryInterface com o ImageRepository? *//
func NewImageRepo(db *gorm.DB) ImageRepositoryInterface {
	return &ImageRepository{DB: db}
}

// Essa função define por meio do reciever que todo ImageRepository deve impllementar um Create
// Dessa forma, é estabelecida uma conexão implicita entre o struct e a Inteface
func (repo *ImageRepository) CreateImage(ctx context.Context, image *model.Image) error {
	if err := repo.DB.Create(image).Error; err != nil {
		return err
	}
	return nil
}
func (repo *ImageRepository) FindByUniqueFileName(ctx context.Context, uniqueName string) error {
	var image model.Image

	_, err := gorm.G[model.Image](repo.DB).Where("unique_file_name = ?", uniqueName).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	fmt.Println(image.UniqueFileName)

	return nil
}
