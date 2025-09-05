package model

import (
	"time"

	"github.com/google/uuid"
)

/*
The Revenue is the value that a client needs to pay after a procediment
The client doesnt need to pay all the value at once
*/
type Revenue struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID             uuid.UUID  `gorm:"type:uuid;not null;column:user_id" json:"user_id"`
	ClientID           *uuid.UUID `gorm:"type:uuid;column:client_id" json:"client_id"` // Ponteiro para ser opcional (NULL)
	ProcedureType      string     `gorm:"type:varchar(255);not null;column:procedure_type" json:"procedure_type"`
	BeneficiaryCpfCnpj string     `gorm:"type:varchar(18);not null;column:beneficiary_cpf_cnpj" json:"beneficiary_cpf_cnpj"`
	Value              float64    `gorm:"type:decimal(10,2);not null" json:"value"`
	TotalPaid          float64    `gorm:"type:decimal(10,2);not null;column:total_paid" json:"total_paid"`
	Description        string     `gorm:"type:text" json:"description"`
	IsDeclared         bool       `gorm:"not null;default:false;column:is_declared" json:"is_declared"`

	IssueDate time.Time `gorm:"type:date;not null" json:"issue_date"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

type RevenueDto struct {
	ClientID           *string   `json:"client_id,omitempty" validate:"omitempty,uuid"`
	ProcedureType      string    `json:"procedure_type" validate:"required"`
	BeneficiaryCpfCnpj string    `json:"beneficiary_cpf_cnpj" validate:"required"`
	Value              float64   `json:"value" validate:"required,gt=0"`
	TotalPaid          float64   `json:"total_paid" validate:"required,gte=0"`
	Description        string    `json:"description"`
	IsDeclared         bool      `json:"is_declared"`
	IssueDate          time.Time `json:"issue_date" validate:"required"`
}
type RevenueSearchCriteria struct {
	UserID uuid.UUID

	// Filtros Opcionais (virão dos query params da URL)
	ClientID      string // Para buscar todas as receitas de um cliente específico
	ProcedureType string // Para buscar receitas por texto na descrição
	StartDate     string // Para filtrar por um período (ex: "2025-01-01")
	EndDate       string // Para filtrar por um período (ex: "2025-12-31")
	OnlyInDebt    bool   // Um "interruptor" para ativar o filtro de dívida
	IsDeclared    bool
}
