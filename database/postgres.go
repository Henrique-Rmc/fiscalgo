package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres" // Driver PostgreSQL para GORM
	"gorm.io/gorm"            // O pacote principal do GORM
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	connStr := "user=postgres password=0102 host=localhost port=5432 dbname=postgres sslmode=disable"

	var err error

	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir o banco de dados PostgreSQL: %w", err)
	}

	sqlDB, err := DB.DB()

	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados PostgreSQL: %w", err)
	}
	sqlDB.SetConnMaxLifetime(time.Minute * 3) // Define o tempo máximo de vida da conexão
	sqlDB.SetMaxOpenConns(10)                 // Define o número máximo de conexões abertas
	sqlDB.SetMaxIdleConns(10)                 // Define o número máximo de conexões ociosas

	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("Erro ao conectar ao banco Postgres: %w", err)
	}
	fmt.Println("Conexão com o banco de dados PostgreSQL estabelecida via GORM.")

	return DB, nil
}

func CloseDB(db *gorm.DB) {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			fmt.Printf("Aviso: Erro ao obter conexão SQL subjacente para fechar: %v\n", err)
			return
		}
		sqlDB.Close()
		fmt.Println("Conexão com o banco de dados PostgreSQL (GORM) fechada.")

	}
}
