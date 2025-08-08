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
	AsksInvoice bool    `json:"asks_invoice"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type ClientData struct {
	Name        string
	Cpf         string
	Phone       string
	Email       string
	AsksInvoice bool
}
