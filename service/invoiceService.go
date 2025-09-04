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
	CreateInvoice(ctx context.Context, invoiceBody *model.InvoiceBody, imageData *model.ImageHeader) (*model.Invoice, error)
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
func (service *InvoiceService) CreateInvoice(ctx context.Context, invoiceBody *model.InvoiceBody, imageHeader *model.ImageHeader) (*model.Invoice, error) {
	var user *model.User
	user, err := service.UserRepo.FindUserById(ctx, invoiceBody.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("Falha ao Encontrar o usuário com o ID correspondente")
	}
	invoiceId := uuid.New()

	objectName := fmt.Sprintf("%s/%s", user.ID.String(), invoiceId.String())

	if imageHeader != nil {
		err = service.ImageService.UploadImageService(ctx, imageHeader, objectName)
		if err != nil {
			return nil, fmt.Errorf("Falha ao subir a imagem no serviço de Bucket")
		}
	}
	
	invoiceStruct := &model.Invoice{
		ID:              invoiceId,
		UserID:          user.ID,
		Description:     invoiceBody.Description,
		Value:           invoiceBody.Value,
		ExpenseCategory: invoiceBody.ExpenseCategory,
		AccessKey:       invoiceBody.AccessKey,
		ImageURL:        objectName,
		IssueDate:       invoiceBody.IssueDate,
	}

	invoice, err := service.InvoiceRepo.CreateInvoice(ctx, invoiceStruct)
	if err != nil {
		fmt.Println("Erro ao Cirar invoice no banco")
		return nil, fmt.Errorf("Falha ao criar Invoice No banco de dados")
	}
	return invoice, nil
}
