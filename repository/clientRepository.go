package repository

import (
	"context"
	"errors"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"gorm.io/gorm"
)

type ClientRepositoryInterface interface {
	CreateClient(ctx context.Context, client *model.Client) error
	FindClientById(ctx context.Context, id string) (*model.Client, error)
}

type ClientRepository struct {
	DB *gorm.DB
}

func NewClientRepository(db *gorm.DB) ClientRepositoryInterface {
	return &ClientRepository{DB: db}
}

func (clientRepo *ClientRepository) CreateClient(ctx context.Context, clientData *model.Client) error {
	err := clientRepo.DB.WithContext(ctx).Create(clientData).Error
	if err != nil {
		return err
	}
	return nil
}

func (clientRepo *ClientRepository) FindClientById(ctx context.Context, id string) (*model.Client, error) {
	var client model.Client
	err := clientRepo.DB.WithContext(ctx).Where("id = ?", id).First(&client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &client, nil
}
