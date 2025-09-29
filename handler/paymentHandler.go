package handler

import (
	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/service"
	"github.com/gofiber/fiber/v2"
)

type PaymentHandlerInterface interface {
	CreatePayment(c *fiber.Ctx) error
}

type PaymentHandler struct {
	PaymentService service.PaymentServiceInterface
}

func NewPaymentHandler(paymentService service.PaymentServiceInterface) PaymentHandlerInterface {
	return &PaymentHandler{
		PaymentService: paymentService,
	}
}

func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	var paymentDto model.PaymentDto

	if err := c.BodyParser(&paymentDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Erro ao processar o corpo da requisição",
		})
	}

	payment, appErr := h.PaymentService.CreatePaymentService(c.Context(), &paymentDto)
	if appErr != nil {
		return c.Status(appErr.StatusCode).JSON(fiber.Map{
			"error":   appErr.Message,
			"details": appErr.Details,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(payment)
}
