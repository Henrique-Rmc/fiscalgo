package routes

import (
	"github.com/Henrique-Rmc/fiscalgo/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, handler handler.UserHandlerInterface) {
	userRoutes := app.Group("/api/users")
	userRoutes.Post("/create", handler.CreateUser)
	userRoutes.Get("/:userId", handler.GetUserById)
}
