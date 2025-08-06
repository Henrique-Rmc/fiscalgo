package utils

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/Henrique-Rmc/fiscalgo/model"
)

func ExtractImageData(file *multipart.FileHeader) (*model.ImageHeader, error) {
	if file == nil {
		return nil, nil
	}
	fileExtension := filepath.Ext(file.Filename)
	lowerFileExtension := strings.ToLower(fileExtension)
	if lowerFileExtension != ".jpg" && lowerFileExtension != ".png" {
		return nil, errors.New("formato de arquivo inválido. Apenas .jpg e .png são permitidos")
	}
	const maxFileSize = 10 * 1024 * 1024

	if file.Size > maxFileSize {
		return nil, errors.New("o arquivo excede o tamanho máximo de 10MB")
	}
	src, err := file.Open()
	if err != nil {
		return nil, errors.New("não foi possível ler o arquivo enviado")
	}
	imageData := model.ImageHeader{
		FileName:      file.Filename,
		FileExtension: lowerFileExtension,
		ContentType:   file.Header.Get("Content-Type"),
		FileSize:      file.Size,
		File:          src,
	}
	return &imageData, nil
}
