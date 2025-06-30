package main

import (
	"fiscalgo/database"
	"fiscalgo/handler"
	"fiscalgo/repository"
	"fiscalgo/routes"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

func main() {
	/**
	A primeira parte da função main realiza a inicialização do banco de dados, essa conexão vai ser injetada no contexto
	Fiber com a função app.Use()

	**/
	db, err := database.InitDB()
	if err != nil {
		fmt.Printf("Erro ao inicializar o banco de dados: %v\n", err)
		os.Exit(1)
	}

	defer database.CloseDB(db)

	app := fiber.New()
	app.Use(
		func(c *fiber.Ctx) error {
			c.Locals("db", db)
			return c.Next()
		})
	imageRepository := repository.NewImageRepo(db)
	imageHandler := handler.NewImageHandler(imageRepository)

	routes.SetupImageRoutes(app, imageHandler)

	port := ":8080"
	fmt.Printf("Servidor Iniciado em http://localhost%s\n", port)

	go func() {
		if err := app.Listen(port); err != nil {
			fmt.Printf("Erro ao iniciar servidor Fiber: %v\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("\nSinal de desligamento recebido. Encerrando servidor Fiber...")

	if err := app.Shutdown(); err != nil {
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Erro ao iniciar servidor Fiber: %v\n", err)
		os.Exit(1)
	}
}
