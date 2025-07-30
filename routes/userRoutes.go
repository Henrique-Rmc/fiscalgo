package routes

import (
	"github.com/Henrique-Rmc/fiscalgo/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, handler handler.UserHandlerInterface) {
	app.Post("/user/create", handler.CreateUser)
}
