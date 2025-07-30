package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Henrique-Rmc/fiscalgo/database"
	"github.com/Henrique-Rmc/fiscalgo/handler"
	"github.com/Henrique-Rmc/fiscalgo/repository"
	"github.com/Henrique-Rmc/fiscalgo/routes"
	"github.com/Henrique-Rmc/fiscalgo/service"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

	ctx := context.Background()
	minioEndPoint := os.Getenv("MINIO_ENDPOINT")
	minioAcessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	bucketName := os.Getenv("BUCKET_NAME")

	minioClient, err := minio.New(minioEndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAcessKey, minioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln("Erro ao inicializar Minio", err)
	}
	found, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatalln("Erro ao verificar bucket no MinIO:", err)
	}
	if !found {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalln("Erro ao criar bucket no MinIO:", err)
		}
		log.Printf("Bucket '%s' criado com sucesso.", bucketName)
	} else {
		log.Printf("Conectado ao bucket '%s'.", bucketName)
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
	userRepository := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	imageRepository := repository.NewImageRepo(db)
	imageService := service.NewImageService(imageRepository, minioClient, bucketName)
	imageHandler := handler.NewImageHandler(imageService)

	routes.SetupImageRoutes(app, imageHandler)
	routes.SetupUserRoutes(app, userHandler)
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
