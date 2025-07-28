package model

import (
	"io"
	"time"

	"github.com/lib/pq"
)

type Image struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	OwnerId        uint           `gorm:"foreignKey" json:"owner_id"`
	UniqueFileName string         `gorm:"not null;unique" json:"unique_file_name"`
	Tags           pq.StringArray `gorm:"type:text[]" json:"tags"`
	Description    string         `json:"description"`
	Url            string         `json:"url"`
	UploadedAt     time.Time      `gorm:"column:uploaded_at;default:CURRENT_TIMESTAMP" json:"uploaded_at"`
}

type ImageData struct {
	FileName      string
	FileExtension string
	ContentType   string
	FileSize      int64
	Body          ImageBody
	File          io.Reader
}

type ImageBody struct {
	OwnerId     uint     `json:"ownerId" xml:"ownerId" form:"ownerId"`
	Description string   `json:"description" xml:"description" form:"description"`
	Tags        []string `json:"tags" xml:"tags" form:"tags"`
}
