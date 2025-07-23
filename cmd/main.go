package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Henrique-Rmc/fiscalgo/database"
	"github.com/Henrique-Rmc/fiscalgo/handler"
	"github.com/Henrique-Rmc/fiscalgo/repository"
	"github.com/Henrique-Rmc/fiscalgo/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		log.Println("Modo de execução: Migrations")

		err := database.RunMigrations("file://./migrations")
		if err != nil {
			log.Fatalf("Erro ao executar migrations: %v", err)
		}
		log.Println("Migrations concluídas. Encerrando o serviço de migration.")
		return
	}

	log.Println("Modo de execução: Aplicação principal")
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Erro ao obter a conexão SQL subjacente para fechar: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Erro ao obter a conexão SQL subjacente para fechar: %v", err)
	}
	defer sqlDB.Close()

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
