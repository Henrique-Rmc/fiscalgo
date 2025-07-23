package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ImageHandlerInterface interface {
	UploadImageHandler(c *fiber.Ctx) error
}
type ImageHandler struct {
	ImageRepo repository.ImageRepositoryInterface
}

func NewImageHandler(imageRepo repository.ImageRepositoryInterface) ImageHandlerInterface {
	return &ImageHandler{ImageRepo: imageRepo}
}

func (h *ImageHandler) UploadImageHandler(c *fiber.Ctx) error {
	//O primeiro passo é extrair a imagem recbida do body, para isso
	//devo usar a função c.FormFile("image") que vai extrair do formulario o parametro image
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Erro ao processar o upload do arquivo.",
			"message": "Verifique se o arquivo foi enviado corretamente no campo 'image'.",
			"details": err.Error(),
		})
	}

	fileExtension := filepath.Ext(file.Filename)
	lowerFileExtension := strings.ToLower(fileExtension)

	if lowerFileExtension != ".jpg" && lowerFileExtension != ".png" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "O formato da imagem selecionada não é valido.",
			"message": "Verifique o formato da imagem e o reenvie.",
		})
	}
	const maxFileSize = 10 * 1024 * 1024

	if file.Size > maxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "O tamanho da imagem enviada ultrapassa o limite suportado.",
			"message": "Verifique o tamanho da imagem e o reenvie.",
		})
	}
	//Agora que tenho a imagem, devo verificar se já existe um diretório que vai
	//salvar a imagem localmente
	uploadDir := "./uploads"
	//os.Stat verifica se o diretorio existe e gera o erro
	//os.IsNotExist verifica se o erro retornado é de que o diretorio nao existe
	//se o erro for de diretorio inexistente, criamos um diretorio
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		fmt.Printf("O diretorio %s não existe...Criando agora", uploadDir)
		//esse formato é a mesma coisa de
		// _, err := os.Mkdir(uploadDir)
		//if err != nil.....

		if err = os.Mkdir(uploadDir, 0755); err != nil {
			fmt.Printf("Erro ao criar o diretorio")
			return c.Status(fiber.StatusInternalServerError).SendString("Erro Ao criar diretorio")
		}
	} else if err != nil {
		fmt.Printf("Erro inesperado ao verificar diretorio")
		return c.Status(fiber.StatusInternalServerError).SendString("Erro ao verificar diretorio de upload")
	}
	newUUID := uuid.New()

	newFileName := newUUID.String() + lowerFileExtension

	filepath := filepath.Join(uploadDir, newFileName)

	if err := c.SaveFile(file, filepath); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Erro ao salvar o arquivo na pasta de destino")
	}

	url := "www.imagem.com"

	imageToSave := model.Image{
		// ID:             uint(0), // GORM geralmente preenche o ID para você em chaves primárias autoincrement
		// Se seu ID for UUID no DB, você precisaria de um tipo string para UniqueFileName
		// e criar um UUID para o ID. No seu modelo, ID é `uint`.
		OwnerId:        1, // Exemplo: Substitua por um ID de usuário real (e.g., de autenticação)
		UniqueFileName: newFileName,
		Tags:           pq.StringArray{"alimentacao"}, // Exemplo: Em um app real, isso viria de um campo de formulário
		Description:    "Nota Fiscal de Almoço",
		Url:            url, // Para ambiente local, a URL pode ser o caminho. Para cloud, seria a URL do bucket.
		UploadedAt:     time.Now(),
	}

	if err := h.ImageRepo.Create(&imageToSave); err != nil {
		fmt.Printf("Erro ao salvar os metadados da imagem no banco de dados: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao salvar os dados da imagem no banco de dados.",
			"message": "Ocorreu um erro interno ao registrar a imagem.",
			"details": err.Error(),
		})
	}

	fmt.Printf("Imagem salva com sucesso")
	return c.Status(fiber.StatusOK).SendString("Imagem salva")

}
