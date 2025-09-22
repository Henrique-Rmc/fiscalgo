package service

import (
	"context"

	"github.com/Henrique-Rmc/fiscalgo/apperror"
	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/repository"
	"github.com/google/uuid"
)

// RevenueServiceInterface define o contrato do serviço.
type RevenueServiceInterface interface {
	Create(ctx context.Context, loggedInUserID uuid.UUID, data *model.RevenueDto) (*model.Revenue, *apperror.AppError)
	Find(ctx context.Context, data *model.RevenueSearchCriteria) ([]*model.Revenue, *apperror.AppError)
}

// revenueService é a implementação.
type revenueService struct {
	RevenueRepo repository.RevenueRepositoryInterface
	ClientRepo  repository.ClientRepositoryInterface // Dependência para validar o cliente
}

// NewRevenueService é o construtor.
func NewRevenueService(revenueRepo repository.RevenueRepositoryInterface, clientRepo repository.ClientRepositoryInterface) RevenueServiceInterface {
	return &revenueService{
		RevenueRepo: revenueRepo,
		ClientRepo:  clientRepo,
	}
}

func (service *revenueService) Find(ctx context.Context, data *model.RevenueSearchCriteria) ([]*model.Revenue, *apperror.AppError) {
	revenues, err := service.RevenueRepo.Find(ctx, data)
	if err != nil {
		// Se o repositório retornar um erro, nós o "embrulhamos" num erro de aplicação.
		return nil, apperror.InternalServer("Ocorreu um erro ao buscar as receitas.", err)
	}
	return revenues, nil
}

/*
O usuario acessa a área de clientes, seleciona um cliente e adiciona uma revenue para ele, dessa forma, o revenue é passado no handler
*/
func (service *revenueService) Create(ctx context.Context, loggedInUserID uuid.UUID, data *model.RevenueDto) (*model.Revenue, *apperror.AppError) {

	parsedClientID, err := uuid.Parse(*data.ClientID)
	if err != nil {
		return nil, apperror.UnprocessableEntity("O client_id fornecido não é um UUID válido.", err)
	}
	_, err = service.ClientRepo.FindClientById(ctx, parsedClientID, loggedInUserID)
	if err != nil {
		return nil, apperror.UnprocessableEntity("O cliente especificado não existe ou não pertence a este usuário.", err)
	}
	clientID := &parsedClientID

	newRevenue := &model.Revenue{
		ID:                 uuid.New(),
		UserID:             loggedInUserID,
		ClientID:           clientID,
		BeneficiaryCpfCnpj: data.BeneficiaryCpfCnpj,
		ProcedureType:      data.ProcedureType,
		Value:              data.Value,
		Debit:              data.Value,
		TotalPaid:          data.TotalPaid,
		Description:        data.Description,
		IsDeclared:         data.IsDeclared,
		IssueDate:          data.IssueDate,
	}

	// 4. Peça ao repositório para salvar no banco.
	if err := service.RevenueRepo.Create(ctx, newRevenue); err != nil {
		// AQUI você pode tratar erros específicos do banco, como chaves duplicadas.
		return nil, apperror.InternalServer("Não foi possível registar a receita.", err)
	}

	return newRevenue, nil
}
