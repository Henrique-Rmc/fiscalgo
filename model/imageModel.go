package model

import (
	"io"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

/*
Ajustar o tipo do OwnerId que no momento esta sendo definido estaticamente com um tipo inconsistente
*/
type Image struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OwnerId        uuid.UUID      `gorm:"type:uuid;foreignKey" json:"owner_id"`
	UniqueFileName string         `gorm:"not null;unique" json:"unique_file_name"`
	Tags           pq.StringArray `gorm:"type:text[]" json:"tags"`
	Description    string         `json:"description"`
	Url            string         `json:"url"`
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
	OwnerId     uuid.UUID `json:"ownerId" xml:"ownerId" form:"ownerId"`
	Description string    `json:"description" xml:"description" form:"description"`
	Tags        []string  `json:"tags" xml:"tags" form:"tags"`
}
