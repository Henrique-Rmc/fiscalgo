package service

import (
	"context"
	"errors"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientServiceInterface interface {
	CreateClient(ctx context.Context, clientData *model.ClientData, idUser string) (*model.Client, error)
	FindClient(ctx context.Context, queryData *model.ClientSearchCriteria) ([]*model.Client, error)
}

type ClientService struct {
	ClientRepo repository.ClientRepositoryInterface
	UserRepo   repository.UserRepositoryInterface
}

/*Sempre recebe interfaces, nunca o objeto literal*/
/*Sempre retorna interface, nunca o objeto criado*/
func NewClientService(clientRepo repository.ClientRepositoryInterface, userRepo repository.UserRepositoryInterface) ClientServiceInterface {
	return &ClientService{
		ClientRepo: clientRepo,
		UserRepo:   userRepo,
	}
}

func (clientService *ClientService) CreateClient(ctx context.Context, clientData *model.ClientData, idUser string) (*model.Client, error) {
	userUUID, err := uuid.Parse(idUser)
	if err != nil {
		return nil, err
	}
	user, err := clientService.UserRepo.FindUserById(ctx, userUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
	}

	clientToSave := &model.Client{
		ID:          uuid.New(),
		UserId:      user.ID,
		Name:        clientData.Name,
		Cpf:         clientData.Cpf,
		Phone:       clientData.Phone,
		Email:       clientData.Email,
		AsksInvoice: clientData.AsksInvoice,
	}
	err = clientService.ClientRepo.CreateClient(ctx, clientToSave)
	if err != nil {
		return nil, err
	}
	return clientToSave, nil
}

func (clientService *ClientService) FindClient(ctx context.Context, queryData *model.ClientSearchCriteria) ([]*model.Client, error) {
	/*Adicionar Regras de Negocio Futuramente*/
	client, err := clientService.ClientRepo.FindClient(ctx, queryData)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return client, nil
}
