package service

import (
	"context"
	"errors"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/repository"
	"gorm.io/gorm"
)

type AuthServiceInterface interface {
	Login(ctx context.Context, email string, plainPassword string) (*model.User, error)
}

type AuthService struct{
	userRepo repository.UserRepositoryInterface
}

func (authService *AuthService) Login(ctx context.Context, email string, plainPassword string) (*model.User, error) {
	user,err := authService.userRepo.FindUserByEmail(ctx,email)
	if err != nil{
		if errors.Is(err, gorm.ErrRecordNotFound){
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return user, nil

}