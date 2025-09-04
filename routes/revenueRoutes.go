package routes

import (
	"github.com/Henrique-Rmc/fiscalgo/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRevenueRoutes(app *fiber.App, handler handler.RevenueHandlerInterface){
	revenueRoutes := app.Group("/api/revenues")
	revenueRoutes.Post("/create", handler.CreateRevenueHandler)
}
