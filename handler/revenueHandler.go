package handler

import (

	// Supondo que tenha um pacote de autenticação
	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/service"
	"github.com/Henrique-Rmc/fiscalgo/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RevenueHandlerInterface interface {
	CreateRevenueHandler(c *fiber.Ctx) error
}

type RevenueHandler struct {
	RevenueService service.RevenueServiceInterface
}

func NewRevenueHandler(revenueService service.RevenueServiceInterface) RevenueHandlerInterface {
	return &RevenueHandler{RevenueService: revenueService}
}

func (h *RevenueHandler) CreateRevenueHandler(c *fiber.Ctx) error {
	loggedInUserID := "6daa7ce0-6594-43ed-b583-c74bd6aa1a13"
	userUuid, _ := uuid.Parse(loggedInUserID)

	// loggedInUserID, err := auth.GetUserIDFromContext(c)
	// if err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Não autenticado."})
	// }

	revenueDto := new(model.RevenueDto)
	if err := c.BodyParser(revenueDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Corpo da requisição inválido."})
	}
	if err := utils.ValidateStruct(revenueDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Dados de entrada inválidos.",
			"details": err.Error(),
		})
	}
	newRevenue, appErr := h.RevenueService.Create(c.Context(), userUuid, revenueDto)
	if appErr != nil {
		return c.Status(appErr.StatusCode).JSON(fiber.Map{
			"error": appErr.Message,
			"details": appErr.Details,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(newRevenue)
}
