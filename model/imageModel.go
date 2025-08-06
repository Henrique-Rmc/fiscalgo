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

type ImageHeader struct {
	FileName      string
	FileExtension string
	ContentType   string
	FileSize      int64
	File          io.Reader
	FileCloser    io.ReadCloser
}
