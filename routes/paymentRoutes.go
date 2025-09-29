package routes

import (
	"github.com/Henrique-Rmc/fiscalgo/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupPaymentRoutes(app *fiber.App, paymentHandler handler.PaymentHandlerInterface) {
	paymentRoutes := app.Group("/api/payments")
	paymentRoutes.Post("/create", paymentHandler.CreatePayment)
}
