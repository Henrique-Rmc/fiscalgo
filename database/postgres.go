package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres" // Driver PostgreSQL para GORM
	"gorm.io/gorm"            // O pacote principal do GORM
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
	log.Printf("DB_PASSWORD: %s", dbPassword) // Cuidado ao imprimir senhas em logs de produção!
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
