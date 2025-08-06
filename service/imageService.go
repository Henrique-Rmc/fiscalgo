package service

import (
	"context"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/repository"
	"github.com/minio/minio-go/v7"
)

type ImageServiceInterface interface {
	UploadImageService(ctx context.Context, data *model.ImageHeader, objectName string) error
	// DownloadImageService(ctx context.Context, userId uuid.UUID, uniqueFileName string) (string, error)
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

/*O service de invoice vai chaamr o uploadImage, dessa forma,id do invoice já vai ser passado diretamente*/
func (service *ImageService) UploadImageService(ctx context.Context, data *model.ImageHeader, objectName string) error {

	_, err := service.MinioClient.PutObject(
		ctx,
		service.BucketName,
		objectName,
		data.File,
		data.FileSize,
		minio.PutObjectOptions{ContentType: data.ContentType})
	if err != nil {
		return err
	}

	return nil

}

// func (service *ImageService) DownloadImageService(ctx context.Context, userId uuid.UUID, uniqueFileName string) (string, error) {
// 	/*
// 		Recebendo o id, buscar o nome da imagem que possui aquele id
// 	*/
// 	image, err := service.ImageRepo.FindByUniqueFileName(ctx, uniqueFileName)
// 	if err != nil {
// 		return "Imagem Buscada não existe no Banco de dados", err
// 	}
// 	user, err := service.UserRepo.FindUserById(ctx, userId)
// 	if err != nil {
// 		return fmt.Sprintf("Usuário com id %s não existe no Banco de dados", userId), err
// 	}
// 	if image.OwnerId.String() != user.ID.String() {
// 		fmt.Printf("Id informado = %s \n Id do dono da Image = %s", userId, image.OwnerId.String())
// 		return "Id do Usuário informado não corresponde ao Id no Banco de dados", err
// 	}
// 	filePath := fmt.Sprintf("users/%s/%s", user.ID, image.OwnerId)
// 	expireTime := 1 * time.Minute
// 	presignedURL, err := service.MinioClient.PresignedGetObject(
// 		ctx,
// 		service.BucketName,
// 		filePath,
// 		expireTime,
// 		nil,
// 	)
// 	if err != nil {
// 		log.Printf("Erro ao gerar URL pré-assinada: %v", err)
// 		return "", err
// 	}

// 	return presignedURL.String(), nil
// }
