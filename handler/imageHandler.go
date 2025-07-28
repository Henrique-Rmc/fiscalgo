package handler

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/service"
	"github.com/gofiber/fiber/v2"
)

type ImageHandlerInterface interface {
	UploadImageHandler(c *fiber.Ctx) error
}

// *objeto imageHandler serve apenas para unir a interface aos metodos*/
type ImageHandler struct {
	ImageService service.ImageServiceInterface
}

func NewImageHandler(imageService service.ImageServiceInterface) ImageHandlerInterface {
	return &ImageHandler{ImageService: imageService}
}

func (handler *ImageHandler) UploadImageHandler(c *fiber.Ctx) error {

	/**Preciso receber o pedaço do multipart que contem o body json com meus dados e ebtão converter no meu
	formato de imageBdy*/
	body := new(model.ImageBody)
	metadata := c.FormValue("metadata")

	if err := json.Unmarshal([]byte(metadata), &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Erro ao acessar matadata")
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Erro ao processar o upload do arquivo.",
			"message": "Verifique se o arquivo foi enviado corretamente no campo 'image'.",
			"details": err.Error(),
		})
	}
	/**Extrai fileExtension do file*/
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
	src, err := file.Open()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erro ao abrir o arquivo."})
	}
	defer src.Close()

	imageData := model.ImageData{
		FileName:      file.Filename,
		FileExtension: lowerFileExtension,
		ContentType:   file.Header.Get("Content-Type"),
		FileSize:      file.Size,
		Body:          *body,
		File:          src,
	}

	savedImage, err := handler.ImageService.UploadImageService(c.Context(), imageData)
	if err != nil {
		fmt.Printf("Erro ao salvar imagem: %v\n", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro ao salvar imagem.",
			"message": "Ocorreu um erro interno ao registrar a imagem.",
			"details": err.Error(),
		})
	}

	fmt.Printf("Imagem salva com sucesso")
	return c.Status(fiber.StatusOK).JSON(savedImage)

}
