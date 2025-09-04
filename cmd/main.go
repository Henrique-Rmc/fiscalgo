package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Henrique-Rmc/fiscalgo/database"
	"github.com/Henrique-Rmc/fiscalgo/database/seed"
	"github.com/Henrique-Rmc/fiscalgo/handler"
	"github.com/Henrique-Rmc/fiscalgo/repository"
	"github.com/Henrique-Rmc/fiscalgo/routes"
	"github.com/Henrique-Rmc/fiscalgo/service"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
)

func main() {
	err := godotenv.Load()
	/*Migrations*/
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
	/*Variaveis de Ambiente*/
	ctx := context.Background()
	minioEndPoint := os.Getenv("MINIO_ENDPOINT")
	minioAcessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	bucketName := os.Getenv("BUCKET_NAME")
	redisAddr := os.Getenv("REDIS_ADDR")

	/*Redis*/
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})

	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Não foi possivel conectar ao Redis pelo erro : %v", err)
	}
	log.Println("Conectado ao Redis com sucesso!")

	/*Minio*/
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
	/*Database*/
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

	/*Seeders*/
	seedFlag := flag.Bool("seed", false, "Se verdadeiro, roda o seeder do banco de dados")
	flag.Parse()
	if *seedFlag {
		log.Println("Iniciando o seeder...")

		// Crie uma instância do seeder e chame o método Seed()
		seeder := seed.NewSeeder(db)
		seeder.Seed()

		log.Println("Seeder finalizado com sucesso.")

		// Termina a aplicação após rodar o seeder
		os.Exit(0)
	}
	/*Iniciando Aplicação*/
	app := fiber.New()
	app.Use(
		func(c *fiber.Ctx) error {
			c.Locals("db", db)
			return c.Next()
		})
	userRepository := repository.NewUserRepo(db)
	imageRepository := repository.NewImageRepo(db)
	clientRepository := repository.NewClientRepository(db)
	invoiceRepository := repository.NewInvoiceRepository(db)
	revenueRepository := repository.NewRevenueRepository(db)

	userService := service.NewUserService(userRepository, rdb)
	imageService := service.NewImageService(imageRepository, userRepository, minioClient, bucketName)
	clientService := service.NewClientService(clientRepository, userRepository, rdb)
	invoiceService := service.NewInvoiceService(invoiceRepository, userRepository, imageService)
	revenueService := service.NewRevenueService(revenueRepository, clientRepository)

	userHandler := handler.NewUserHandler(userService)
	invoiceHandler := handler.NewInvoiceHandler(invoiceService)
	clientHandler := handler.NewClientHandler(clientService)
	revenueHandler := handler.NewRevenueHandler(revenueService)

	routes.SetupClientRoutes(app, clientHandler)
	routes.SetupUserRoutes(app, userHandler)
	routes.SetupInvoiceRoutes(app, invoiceHandler)
	routes.SetupRevenueRoutes(app, revenueHandler)

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
