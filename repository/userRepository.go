package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	FindUserById(ctx context.Context, id uuid.UUID) (*model.User, error)
}
type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{DB: db}
}

func (userRepo *UserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := userRepo.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepo *UserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := userRepo.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	fmt.Println(user)
	return &user, err
}

func (userRepo *UserRepository) FindUserById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := userRepo.DB.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	fmt.Println(user)
	return &user, err
}
