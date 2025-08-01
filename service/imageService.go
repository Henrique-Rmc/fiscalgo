package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/repository"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type ImageServiceInterface interface {
	UploadImageService(ctx context.Context, data model.ImageData) (*model.Image, error)
	DownloadImageService(ctx context.Context, userId uuid.UUID, uniqueFileName string) (string, error)
}

type ImageService struct {
	ImageRepo   repository.ImageRepositoryInterface
	UserRepo    repository.UserRepositoryInterface
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
func NewImageService(imageRepo repository.ImageRepositoryInterface, userRepo repository.UserRepositoryInterface, minioC *minio.Client, bucketName string) ImageServiceInterface {
	return &ImageService{
		ImageRepo:   imageRepo,
		UserRepo:    userRepo,
		MinioClient: minioC,
		BucketName:  bucketName,
	}
}

func (service *ImageService) UploadImageService(ctx context.Context, data model.ImageData) (*model.Image, error) {
	var user *model.User

	user, err := service.UserRepo.FindUserById(ctx, data.Body.OwnerId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	objectName := fmt.Sprintf("users/%s/%s", user.ID.String(), data.FileName)

	_, err = service.MinioClient.PutObject(
		ctx,
		service.BucketName,
		objectName,
		data.File,
		data.FileSize,
		minio.PutObjectOptions{ContentType: data.ContentType})
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/%s/%s", os.Getenv("MINIO_ENDPOINT"), service.BucketName, data.FileName)

	newUUID := uuid.New()
	image := model.Image{
		ID:             newUUID,
		OwnerId:        data.Body.OwnerId,
		UniqueFileName: data.FileName + (newUUID).String(),
		Tags:           data.Body.Tags,
		Description:    data.Body.Description,
		Url:            url,
	}
	if err := service.ImageRepo.CreateImage(ctx, &image); err != nil {
		fmt.Printf("Erro ao Inserir Imagem no Banco")
		return nil, err
	}
	return &image, nil

}

func (service *ImageService) DownloadImageService(ctx context.Context, userId uuid.UUID, uniqueFileName string) (string, error) {
	/*
		Recebendo o id, buscar o nome da imagem que possui aquele id
	*/
	image, err := service.ImageRepo.FindByUniqueFileName(ctx, uniqueFileName)
	if err != nil {
		return "Imagem Buscada não existe no Banco de dados", err
	}
	user, err := service.UserRepo.FindUserById(ctx, userId)
	if err != nil {
		return fmt.Sprintf("Usuário com id %s não existe no Banco de dados", userId), err
	}
	if image.OwnerId.String() != user.ID.String() {
		fmt.Printf("Id informado = %s \n Id do dono da Image = %s", userId, image.OwnerId.String())
		return "Id do Usuário informado não corresponde ao Id no Banco de dados", err
	}
	filePath := fmt.Sprintf("users/%s/%s", user.ID, image.OwnerId)
	expireTime := 1 * time.Minute
	presignedURL, err := service.MinioClient.PresignedGetObject(
		ctx,
		service.BucketName,
		filePath,
		expireTime,
		nil,
	)
	if err != nil {
		log.Printf("Erro ao gerar URL pré-assinada: %v", err)
		return "", err
	}

	return presignedURL.String(), nil
}
