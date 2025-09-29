package service

import (
	"context"

	"github.com/Henrique-Rmc/fiscalgo/apperror"
	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/repository"
	"github.com/google/uuid"
)

type PaymentServiceInterface interface {
	CreatePaymentService(ctx context.Context, paymentDto *model.PaymentDto) (*model.Payment, *apperror.AppError)
}

type paymentService struct {
	PaymentRepo repository.PaymentRepositoryInterface
	RevenueRepo repository.RevenueRepositoryInterface
}

func NewPaymentService(paymentRepo repository.PaymentRepositoryInterface, revenueRepo repository.RevenueRepositoryInterface) PaymentServiceInterface {
	return &paymentService{
		PaymentRepo: paymentRepo,
		RevenueRepo: revenueRepo,
	}
}

//Garantir que só o usuário logado pode realizar ações sobre o pagamento e seus clientes
//Essa validação deve ter que ser feita provavelmente por um middleware ou no controller
func (paymentService *paymentService) CreatePaymentService(ctx context.Context, paymentDto *model.PaymentDto) (*model.Payment, *apperror.AppError) {
	if paymentDto.Debit < paymentDto.ValuePaid{
		return nil, apperror.UnprocessableEntity("O valor do débito não pode ser menor que o valor Pago", nil)
	}	

	newPayment := &model.Payment{
		ID: uuid.New(),
		RevenueId: paymentDto.ID,
		Debit: paymentDto.Debit,
		ValuePaid: paymentDto.ValuePaid,
		PaymentDate: paymentDto.PaymentDate,
	}
	
	if err := paymentService.PaymentRepo.Create(ctx, newPayment); err != nil {
		return nil, apperror.InternalServer("Não foi possível registar o pagamento.", err)
	}

	return newPayment, nil
}
