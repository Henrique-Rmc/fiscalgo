package main

import (
	"fiscalgo/routes"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	routes.SetupImageRoutes(app)

	port := ":8080"
	fmt.Printf("Servidor Iniciado em http://localhost%s\n", port)
	err := app.Listen(port)
	if err != nil {
		fmt.Printf("Erro ao iniciar servidor Fiber: %v\n", err)
		os.Exit(1)
	}
}
