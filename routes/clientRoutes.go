package routes

import (
	"github.com/Henrique-Rmc/fiscalgo/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupClientRoutes(app *fiber.App, handler handler.ClientHandlerInterface) {
	app.Post("/client/create", handler.CreateClient)
}
