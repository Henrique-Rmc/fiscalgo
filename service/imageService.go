package service

import (
	"context"
	"fmt"
	"os"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/repository"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type ImageServiceInterface interface {
	UploadImageService(ctx context.Context, data model.ImageData) (*model.Image, error)
}

type ImageService struct {
	ImageRepo   repository.ImageRepositoryInterface
	MinioClient *minio.Client
	BucketName  string
}

/*
*
A interface define todos os metodos que o objeto vai implementar
a Struct é a classe e ela precisa do imageRepo para ser construida
O NewImageHandler é o construtor que vai montar um ImageHandler ao receber um repository, sendo assim não precisa ser mapeado
pela interface
*
*/
func NewImageService(imageRepo repository.ImageRepositoryInterface, minioC *minio.Client, bucketName string) ImageServiceInterface {
	return &ImageService{
		ImageRepo:   imageRepo,
		MinioClient: minioC,
		BucketName:  bucketName,
	}
}

func (service *ImageService) UploadImageService(ctx context.Context, data model.ImageData) (*model.Image, error) {
	_, err := service.MinioClient.PutObject(
		ctx,
		service.BucketName,
		data.FileName,
		data.File,
		data.FileSize,
		minio.PutObjectOptions{ContentType: data.ContentType})
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/%s/%s", os.Getenv("MINIO_ENDPOINT"), service.BucketName, data.FileName)

	newUUID := uuid.New()
	image := model.Image{
		OwnerId:        data.Body.OwnerId,
		UniqueFileName: data.FileName + (newUUID).String(),
		Tags:           data.Body.Tags,
		Description:    data.Body.Description,
		Url:            url,
	}
	if err := service.ImageRepo.Create(&image); err != nil {
		fmt.Printf("Erro ao Inserir Imagem no Banco")
		return nil, err
	}
	return &image, nil

}
