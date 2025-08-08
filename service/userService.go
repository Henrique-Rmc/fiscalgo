package service

import (
	"context"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, data *model.UserData) (*model.User, error)
}

type UserService struct {
	UserRepo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) UserServiceInterface {
	return &UserService{UserRepo: repo}
}

func (service *UserService) CreateUser(ctx context.Context, data *model.UserData) (*model.User, error) {
	/*Receber o UserModel, aplicar validação de dados e criptografar senha
	Enviar para o repositorio salvar o usuario
	*/
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	userId := uuid.New()
	userToSave := &model.User{
		ID:                   userId,
		Name:                 data.Name,
		Email:                data.Email,
		CPF:                  data.CPF,
		PasswordHash:         string(hashedPassword),
		Occupation:           data.Occupation,
		ProfessionalRegistry: data.ProfessionalRegistry,
	}
	_, err = service.UserRepo.CreateUser(ctx, userToSave)
	if err != nil {
		return nil, err
	}
	return userToSave, nil
}
