package routes

import (
	"github.com/Henrique-Rmc/fiscalgo/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupClientRoutes(app *fiber.App, clientHandler handler.ClientHandlerInterface) {
	clientRoutes := app.Group("/api/clients")
	clientRoutes.Post("/create", clientHandler.CreateClient)
	clientRoutes.Get("/", clientHandler.FindClient)
	clientRoutes.Get("/:clientId", clientHandler.GetCliendById)
}
