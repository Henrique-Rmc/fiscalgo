package repository

import (
	"fiscalgo/model"

	"gorm.io/gorm"
)

// Criamos uma interface para definir que quem implementa essa ineterface, deve possuir os métodos da interface
// Nesse caso ,
type ImageRepositoryInterface interface {
	Create(image *model.Image) error
	Get(ID uint) (*model.Image, error)
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
func (r *ImageRepository) Create(image *model.Image) error {
	if err := r.DB.Create(image).Error; err != nil {
		return err
	}
	return nil
}
func (r *ImageRepository) Get(ID uint) (*model.Image, error) {
	var image model.Image

	result := r.DB.Find(&image, ID)

	if result.Error != nil {
		// Retorna nil para o ponteiro da imagem e o erro encontrado.
		return nil, result.Error
	}
	// Se não houve erro, retorna o endereço da struct 'image' (ponteiro) e nil para erro.
	return &image, nil
}
