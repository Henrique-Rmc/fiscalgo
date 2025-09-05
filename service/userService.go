package service

import (
	"context"
	"errors"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, data *model.UserDto) (*model.User, error)
	GetUserById(ctx context.Context, userId uuid.UUID) (*model.User, error)
}

type userService struct {
	UserRepo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface, redis *redis.Client) UserServiceInterface {
	return &userService{UserRepo: repo}
}

func (service *userService) CreateUser(ctx context.Context, data *model.UserDto) (*model.User, error) {
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

func (service *userService) GetUserById(ctx context.Context, userId uuid.UUID) (*model.User, error) {
	user, err := service.UserRepo.FindUserById(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return user, nil
}
