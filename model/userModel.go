package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                   uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name                 string    `gorm:"type:varchar(255);not null" json:"name"`
	Email                string    `gorm:"type:varchar(255);not null;unique" json:"email"`
	CPF                  string    `gorm:"type:varchar(11);not null;unique" json:"cpf"`
	PasswordHash         string    `gorm:"type:varchar(255);not null;column:password_hash" json:"-"`
	Occupation           string    `gorm:"type:varchar(255);not null" json:"occupation"`
	ProfessionalRegistry string    `gorm:"type:varchar(255);not null;unique;column:professional_registry" json:"professional_registry"`
	CreatedAt            time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type UserData struct {
	Name                 string
	Email                string
	CPF                  string
	Password             string
	Occupation           string
	ProfessionalRegistry string
}
