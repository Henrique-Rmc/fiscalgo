package routes

import (
	"github.com/Henrique-Rmc/fiscalgo/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupInvoiceRoutes(app *fiber.App, handler handler.InvoiceHandlerInterface){
	app.Post("/invoice/create", handler.CreateInvoiceHandler)
}
