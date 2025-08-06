package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceRepositoryInterface interface {
	CreateInvoice(ctx context.Context, invoiceData *model.Invoice) (*model.Invoice, error)
}

type InvoiceRepository struct {
	DB *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepositoryInterface {
	return &InvoiceRepository{DB: db}
}

func (invoiceRepo *InvoiceRepository) CreateInvoice(ctx context.Context, invoiceData *model.Invoice) (*model.Invoice, error) {
	/*Vai salvar um invoice no banco de dados*/
	err := invoiceRepo.DB.WithContext(ctx).Create(invoiceData).Error
	if err != nil {
		return nil, err
	}
	return invoiceData, nil
}

func (invoiceRepo *InvoiceRepository) FindInvoiceById(ctx context.Context, id uuid.UUID) (*model.Invoice, error) {
	var invoice model.Invoice
	err := invoiceRepo.DB.WithContext(ctx).Where("id = ?", id).First(&invoice).Error
	if err != nil {
		log.Printf("ERRO DO GORM AO CRIAR INVOICE: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		log.Printf("ERRO DO GORM AO CRIAR INVOICE: %v", err)
		return nil, fmt.Errorf("Falha ao encontrar o invoice")
	}
	return &invoice, nil
}
