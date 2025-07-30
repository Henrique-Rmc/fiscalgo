package model

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `gorm:"type:uuid; primaryKey" json:"id"`
	Name         string    `gorm:"not null" json:"name"`
	Email        string    `gorm:"not null;unique" json:"email"`
	PasswordHash string    `gorm:"column:password_hash;not null" json:"-"`
	Occupation   string    `json:"occupation"`
}

type UserData struct {
	Name       string
	Email      string
	Password   string
	Occupation string
}
