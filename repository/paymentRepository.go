package repository

import (
	"context"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"gorm.io/gorm"
)

type PaymentRepositoryInterface interface {
	Create(c context.Context, payment *model.Payment) error
}

type PaymentRepository struct {
	DB *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepositoryInterface {
	return &PaymentRepository{DB: db}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	if err := r.DB.WithContext(ctx).Create(payment).Error; err != nil {
		return err
	}
	return nil
}
