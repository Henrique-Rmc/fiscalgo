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
	if err != nil{
		return nil, err
	}
	userId := uuid.New()
	user := &model.User{
		ID: userId,
		Name: data.Name,
		Email: data.Email,
		PasswordHash: string(hashedPassword),
		Occupation: data.Occupation,
	}
	createdUser, err := service.UserRepo.CreateUser(ctx, user)
	if err != nil{
		return nil,err
	}
	return createdUser, nil
}
