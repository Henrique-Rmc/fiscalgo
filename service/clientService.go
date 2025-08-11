package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ClientServiceInterface interface {
	CreateClient(ctx context.Context, clientData *model.ClientData, idUser uuid.UUID) (*model.Client, error)
	FindClient(ctx context.Context, queryData *model.ClientSearchCriteria) ([]*model.Client, error)
	GetById(ctx context.Context, clientId uuid.UUID, userId uuid.UUID) (*model.Client, error)
}

type clientService struct {
	ClientRepo repository.ClientRepositoryInterface
	UserRepo   repository.UserRepositoryInterface
	RedisC     *redis.Client
}

/*Sempre recebe interfaces, nunca o objeto literal*/
/*Sempre retorna interface, nunca o objeto criado*/
func NewClientService(clientRepo repository.ClientRepositoryInterface, userRepo repository.UserRepositoryInterface, redisC *redis.Client) ClientServiceInterface {
	return &clientService{
		ClientRepo: clientRepo,
		UserRepo:   userRepo,
		RedisC:     redisC,
	}
}

func (clientService *clientService) CreateClient(ctx context.Context, clientData *model.ClientData, idUser uuid.UUID) (*model.Client, error) {

	user, err := clientService.UserRepo.FindUserById(ctx, idUser)
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

func (clientService *clientService) FindClient(ctx context.Context, queryData *model.ClientSearchCriteria) ([]*model.Client, error) {
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

func (clientService *clientService) GetById(ctx context.Context, clientId uuid.UUID, userId uuid.UUID) (*model.Client, error) {
	key := fmt.Sprintf("client:%s", clientId.String())

	cachedClientJson, err := clientService.RedisC.Get(ctx, key).Result()
	if err == nil {
		log.Println("Cache Hit para a chave:", key)
		var client model.Client
		err := json.Unmarshal([]byte(cachedClientJson), &client)
		if err == nil {
			if client.UserId == userId {
				return &client, nil
			}
		}
	}
	if err != nil && err != redis.Nil {
		log.Printf("Erro ao acessar o cache Redis: %v", err)
	}
	log.Println("CACHE MISS para a chave:", key)
	clientFromDb, err := clientService.ClientRepo.FindClientById(ctx, clientId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	clientJSON, err := json.Marshal(clientFromDb)
	if err != nil {
		return nil, err
	} else {
		err = clientService.RedisC.Set(ctx, key, clientJSON, 10*time.Minute).Err()
		if err != nil {
			log.Printf("Erro ao salvar no cache Redis: %v", err)
		}
	}

	return clientFromDb, nil

}
