package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceServiceInterface interface {
	CreateInvoice(ctx context.Context, InvoiceDto *model.InvoiceDto, imageData *model.ImageDto) (*model.Invoice, error)
	// FindInvoice(ctx context.Context, )
}

type InvoiceService struct {
	InvoiceRepo  repository.InvoiceRepositoryInterface
	UserRepo     repository.UserRepositoryInterface
	ImageService ImageServiceInterface
}

func NewInvoiceService(invoiceRepo repository.InvoiceRepositoryInterface, userRepo repository.UserRepositoryInterface, imageService ImageServiceInterface) InvoiceServiceInterface {
	return &InvoiceService{InvoiceRepo: invoiceRepo, UserRepo: userRepo, ImageService: imageService}
}

/*
Eu deveria criar 2 funções no service? Uma chamada quando existe imagem e uma chamada quando não existe?
*/
func (service *InvoiceService) CreateInvoice(ctx context.Context, InvoiceDto *model.InvoiceDto, ImageDto *model.ImageDto) (*model.Invoice, error) {
	var user *model.User
	user, err := service.UserRepo.FindUserById(ctx, InvoiceDto.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("Falha ao Encontrar o usuário com o ID correspondente")
	}
	invoiceId := uuid.New()

	objectName := fmt.Sprintf("%s/%s", user.ID.String(), invoiceId.String())

	if ImageDto != nil {
		err = service.ImageService.UploadImageService(ctx, ImageDto, objectName)
		if err != nil {
			return nil, fmt.Errorf("Falha ao subir a imagem no serviço de Bucket")
		}
	}

	invoiceStruct := &model.Invoice{
		ID:              invoiceId,
		UserID:          user.ID,
		Description:     InvoiceDto.Description,
		Value:           InvoiceDto.Value,
		ExpenseCategory: InvoiceDto.ExpenseCategory,
		IsDeclared:      false,
		AccessKey:       InvoiceDto.AccessKey,
		ImageURL:        objectName,
		IssueDate:       InvoiceDto.IssueDate,
	}

	invoice, err := service.InvoiceRepo.CreateInvoice(ctx, invoiceStruct)
	if err != nil {
		fmt.Println("Erro ao Cirar invoice no banco")
		return nil, fmt.Errorf("Falha ao criar Invoice No banco de dados")
	}
	return invoice, nil
}
