package model

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID          uuid.UUID `gorm:"type:uuid; primaryKey" json:"id"`
	UserId      uuid.UUID `gorm:"type:uuid;not null;collumn:user_id" json:"user_id"`
	Name        string    `gorm:"not null" json:"name"`
	Cpf         string    `json:"cpf"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	AsksInvoice bool      `json:"asks_invoice"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type ClientData struct {
	Name        string  `json:"name" validate:"required,min=2"`
	Cpf         string  `json:"cpf" validate:"required,len=11"` // Exige um CPF com exatamente 11 dígitos
	Phone       string  `json:"phone" validate:"required"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email"` // Opcional, mas se for enviado, deve ser um email válido
	AsksInvoice bool    `json:"asks_invoice"`                               // Bools são validados como 'true' ou 'false'
}

type ClientSearchCriteria struct {
	UserId uuid.UUID
	CPF    string
	ID     string
	Name   string
}
