package routes

import (
	"fiscalgo/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupImageRoutes(app *fiber.App) {
	app.Post("/upload", handler.UploadImageHandler)
}
