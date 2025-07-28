package routes

import (
	"github.com/Henrique-Rmc/fiscalgo/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupImageRoutes(app *fiber.App, handler handler.ImageHandlerInterface) {
	app.Post("/upload", handler.UploadImageHandler)
	app.Post("/download", handler.DownloadImageHandler)
}
