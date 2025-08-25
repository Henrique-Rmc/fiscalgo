package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/service"
	"github.com/Henrique-Rmc/fiscalgo/utils"
	"github.com/gofiber/fiber/v2"
)

type InvoiceHandlerInterface interface {
	CreateInvoiceHandler(c *fiber.Ctx) error
}

type InvoiceHandler struct {
	InvoiceService service.InvoiceServiceInterface
}

func NewInvoiceHandler(invoiceService service.InvoiceServiceInterface) InvoiceHandlerInterface {
	return &InvoiceHandler{InvoiceService: invoiceService}
}

func (handler *InvoiceHandler) CreateInvoiceHandler(c *fiber.Ctx) error {
	body := new(model.InvoiceBody)
	metadata := c.FormValue("metadata")

	if err := json.Unmarshal([]byte(metadata), &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Erro ao acessar metadata",
		})
	}
	fileHeader, err := c.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Erro ao acessar campo do arquivo image",
		})
	}
	imageData, err := utils.ExtractImageData(fileHeader)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if imageData != nil && imageData.FileCloser != nil {
		defer imageData.FileCloser.Close()
	}
	fmt.Println(body.UserId, body.Value, body.AccessKey, body.Description, body.ExpenseCategory)
	invoice, err := handler.InvoiceService.CreateInvoice(c.Context(), body, imageData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "JSON invalido",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(invoice)
}
