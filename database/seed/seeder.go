// database/seed/seeder.go
package seed

import (
	"log"

	"github.com/Henrique-Rmc/fiscalgo/model" // Ajuste o caminho do import
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seeder contém a conexão com o banco de dados.
type Seeder struct {
	DB *gorm.DB
}

// NewSeeder é o construtor do nosso seeder.
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{DB: db}
}

// Seed é o método principal que orquestra a criação dos dados.
func (s *Seeder) Seed() {
	log.Println("Iniciando o processo de seeding...")

	// 1. Crie um usuário. A função seedUser retornará o usuário criado ou já existente.
	user, err := s.seedUser()
	if err != nil {
		log.Fatalf("Erro ao semear usuário: %v", err)
	}
	log.Printf("Usuário semeado/encontrado: %s", user.Email)

	// 2. Crie 10 clientes (pacientes) para este usuário.
	for i := 0; i < 10; i++ {
		_, err := s.seedClient(user)
		if err != nil {
			log.Fatalf("Erro ao semear cliente: %v", err)
		}
	}
	log.Println("10 clientes semeados com sucesso.")
}

// seedUser cria um usuário de teste se ele não existir.
func (s *Seeder) seedUser() (*model.User, error) {
	// Gere um hash para a senha '123456'
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Crie o objeto de usuário com dados fixos/falsos
	user := model.User{
		Email:                "dentista@email.com", // Usaremos o email como chave para verificar a existência
		Name:                 "Dr. Exemplo da Silva",
		CPF:                  faker.CCNumber(),
		PasswordHash:         string(passwordHash),
		Occupation:           "Dentista",
		ProfessionalRegistry: "CRO-SP 12345",
	}

	// Lógica de idempotência:
	// Tenta encontrar um usuário com este email. Se não encontrar, prepara para criar.
	var existingUser model.User
	err = s.DB.Where("email = ?", user.Email).FirstOrInit(&existingUser, user).Error
	if err != nil {
		return nil, err
	}

	// Se o ID for zero, significa que o usuário não existia e o GORM o inicializou. Agora vamos criar.
	if existingUser.ID == uuid.Nil {
		// Antes de criar, precisamos gerar o UUID do ID
		existingUser.ID = uuid.New()
		if err := s.DB.Create(&existingUser).Error; err != nil {
			return nil, err
		}
		log.Println("Novo usuário de teste criado.")
	}

	return &existingUser, nil
}

// seedClient cria um cliente de teste para um usuário específico.
func (s *Seeder) seedClient(user *model.User) (*model.Client, error) {
	client := model.Client{
		Name:  faker.Name(),
		Cpf:   faker.CCNumber(),
		Phone: faker.Phonenumber(),
		Email: faker.Email(),
		// Associa o cliente ao usuário passado como parâmetro
		UserId: user.ID,
	}

	// Verifique se um cliente com este CPF já existe PARA ESTE USUÁRIO
	var existingClient model.Client
	err := s.DB.Where(model.Client{Cpf: client.Cpf, UserId: user.ID}).FirstOrInit(&existingClient, client).Error
	if err != nil {
		return nil, err
	}

	if existingClient.ID == uuid.Nil {
		existingClient.ID = uuid.New()
		if err := s.DB.Create(&existingClient).Error; err != nil {
			return nil, err
		}
	}

	return &existingClient, nil
}
