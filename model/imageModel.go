package model

import (
	"io"

	"github.com/google/uuid"
)

/*
Ajustar o tipo do OwnerId que no momento esta sendo definido estaticamente com um tipo inconsistente
*/
type Image struct {
	InvoiceId  uuid.UUID `gorm:"type:uuid;primaryKey" json:"invoice_id"`
	ObjectPath string    `gorm:"not null;unique" json:"object_path"`
}

// type ImageData struct {
// 	FileName      string
// 	FileExtension string
// 	ContentType   string
// 	FileSize      int64
// 	InvoiceId     uuid.UUID
// 	File          io.Reader
// }

type ImageDto struct {
	FileName      string        `validate:"required"`
	FileExtension string        `validate:"required,oneof=.jpg .jpeg .png"`
	ContentType   string        `validate:"required"`
	FileSize      int64         `validate:"required,lte=10485760"` // lte = Less Than or Equal (10MB)
	File          io.ReadCloser // Unificado, pois ReadCloser já é um Reader
}