package routes

import (
	"fiscalgo/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupImageRoutes(app *fiber.App, imageHandler *handler.ImageHander) {
	app.Post("/upload", imageHandler.UploadImageHandler)
}
