package database

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	log.Printf("DB_HOST: %s", dbHost)
	log.Printf("DB_PORT: %s", dbPort)
	log.Printf("DB_USER: %s", dbUser)
	log.Printf("DB_PASSWORD: %s", dbPassword)
	log.Printf("DB_NAME: %s", dbName)

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		return nil, fmt.Errorf("Erro: Uma ou mais variáveis de ambiente do banco de dados estão vazias.")
	}
	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Erro ao conectar ao banco de dados com GORM: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("Erro ao obter a conexão SQL subjacente do GORM: %w", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("Erro ao testar a conexão com o banco de dados: %w", err)
	}

	fmt.Println("Conexão com o banco de dados PostgreSQL via GORM estabelecida com sucesso!")
	return db, nil
}

func RunMigrations(migrationPath string) error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		return fmt.Errorf("Erro: Uma ou mais variáveis de ambiente do banco de dados estão vazias para migrations.")
	}

	// DSN para o golang-migrate (formato ligeiramente diferente do GORM)
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	log.Printf("Executando migrations de %s para %s", migrationPath, databaseURL)

	m, err := migrate.New(
		migrationPath,
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("Falha ao criar instância de migrate: %w", err)
	}
	defer m.Close() // Garante que a instância de migrate seja fechada

	// Aplica todas as migrations pendentes
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Falha ao aplicar migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("Nenhuma nova migration para aplicar.")
	} else {
		log.Println("Migrations aplicadas com sucesso!")
	}

	return nil
}
