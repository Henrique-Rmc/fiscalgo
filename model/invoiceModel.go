package model

import (
	"time"

	"github.com/google/uuid"
)

type Invoice struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;column:user_id" json:"user_id"`
	Description     string    `gorm:"type:text;not null" json:"description"`
	Value           float64   `gorm:"type:decimal(10,2);not null" json:"value"`
	ExpenseCategory string    `gorm:"type:varchar(255);not null;column:expense_category" json:"expense_category"`
	AccessKey       string    `gorm:"type:varchar(44)" json:"access_key,omitempty"`
	ImageURL        string    `gorm:"type:varchar(255);column:image_url" json:"image_url,omitempty"`
	IssueDate       time.Time `gorm:"type:date;not null;column:issue_date" json:"issue_date"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type InvoiceBody struct {
	UserId          uuid.UUID
	Description     string
	Value           float64
	ExpenseCategory string
	AccessKey       string
	IssueDate       time.Time
}
