package repository

import (
	"context"
	"errors"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientRepositoryInterface interface {
	CreateClient(ctx context.Context, client *model.Client) error
	FindClientById(ctx context.Context, clientId uuid.UUID, userId uuid.UUID) (*model.Client, error)
	FindClient(ctx context.Context, queryData *model.ClientSearchCriteria) ([]*model.Client, error)
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

func (clientRepo *ClientRepository) FindClientById(ctx context.Context, clientId uuid.UUID, userId uuid.UUID) (*model.Client, error) {
	var client model.Client

	err := clientRepo.DB.WithContext(ctx).Where("id = ? AND user_id = ?", clientId, userId).First(&client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &client, nil
}

func (clientRepo *ClientRepository) FindClient(ctx context.Context, queryData *model.ClientSearchCriteria) ([]*model.Client, error) {
	var clients []*model.Client
	query := clientRepo.DB.WithContext(ctx).Where("user_id = ?", queryData.UserId)

	if queryData.CPF != "" {
		query = query.Where("cpf = ?", queryData.CPF)
	}
	if queryData.Name != "" {
		query = query.Where("name = ?", queryData.Name)
	}
	if queryData.ID != "" {
		idUUID, err := uuid.Parse(queryData.ID)
		if err == nil {
			query = query.Where("id = ?", idUUID)
		}
	}
	if err := query.Find(&clients).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return clients, nil
}
