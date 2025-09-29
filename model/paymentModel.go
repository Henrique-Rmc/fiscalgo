package model

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct{
	ID uuid.UUID `gorm:"type:uuid; primaryKey" json:"id"`
	RevenueId uuid.UUID `gorm:"type:uuid;not null;collumn:revenue_id" json:"revenue_id"`
	Debit float32 `json:"debit"`
	ValuePaid  float32 `json:"value_paid"`
	PaymentDate time.Time `json:"payment_date"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

type PaymentDto struct {
	ID uuid.UUID `json:"id"`
	Debit float32 `json:"debit"`
	ValuePaid  float32 `json:"value_paid"`
	PaymentDate time.Time   `json:"payment_date"`
}