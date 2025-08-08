package handler

import (
	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ClientHandlerInterface interface {
	CreateClient(c *fiber.Ctx) error
	FindClient(c *fiber.Ctx) error
}

type ClientHandler struct {
	ClientService service.ClientServiceInterface
}

func NewClientHandler(clientService service.ClientServiceInterface) ClientHandlerInterface {
	return &ClientHandler{ClientService: clientService}
}

func (clientHandler *ClientHandler) CreateClient(c *fiber.Ctx) error {
	/*Quando Criar Login, ajustar para resgatar id pelo usuario logado*/
	idUser := "6daa7ce0-6594-43ed-b583-c74bd6aa1a13"
	clientData := new(model.ClientData)
	if err := c.BodyParser(clientData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "JSON invalido",
		})
	}
	client, err := clientHandler.ClientService.CreateClient(c.Context(), clientData, idUser)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(client)
}

func (clientHandler *ClientHandler) FindClient(c *fiber.Ctx) error {
	idUser := "6daa7ce0-6594-43ed-b583-c74bd6aa1a13"
	userUUID, err := uuid.Parse(idUser)
	if err != nil {
		return err
	}
	criteria := model.ClientSearchCriteria{
		UserId: userUUID,
		CPF:    c.Query("cpf"),
		Name:   c.Query("name"),
		ID:     c.Query("id"),
	}
	client, err := clientHandler.ClientService.FindClient(c.Context(), &criteria)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "JSON invalido",
		})
	}
	return c.Status(fiber.StatusOK).JSON(client)

}
