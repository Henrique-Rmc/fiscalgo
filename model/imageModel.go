package model

import "time"

type Image struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	OwnerId        uint      `gorm:"foreignKey" json:"owner_id"`
	UniqueFileName string    `gorm:"not null;unique" json:"unique_file_name"`
	Tags           []string  `gorm:"type:text[]" json:"tags"`
	Description    string    `json:"description"`
	Url            string    `json:"url"`
	UploadedAt     time.Time `gorm:"column:uploaded_at;default:CURRENT_TIMESTAMP" json:"uploaded_at"`
}
