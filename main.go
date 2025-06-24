package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Post("/upload", uploadImageHandler)

	port := ":8080"
	fmt.Printf("Servidor Iniciado em http://localhost%s\n", port)
	err := app.Listen(port)
	if err != nil {
		fmt.Printf("Erro ao iniciar servidor Fiber: %v\n", err)
		os.Exit(1)
	}
}

func uploadImageHandler(c *fiber.Ctx) error {
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
	const maxFileSize = 5 * 1024 * 1024

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

	filepath := filepath.Join(uploadDir, file.Filename)

	if err := c.SaveFile(file, filepath); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Erro ao salvar o arquivo na pasta de destino")
	}
	fmt.Printf("Imagem salva com sucesso")
	return c.Status(fiber.StatusOK).SendString("Imagem salva")

}
