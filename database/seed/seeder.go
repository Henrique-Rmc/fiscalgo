// database/seed/seeder.go
package seed

import (
	"log"

	"github.com/Henrique-Rmc/fiscalgo/model" // Ajuste o caminho do import
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type Seeder struct {
	DB *gorm.DB
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{DB: db}
}

func (s *Seeder) Seed() {
	log.Println("Iniciando o processo de seeding...")

	user, err := s.seedUser()
	if err != nil {
		log.Fatalf("Erro ao semear usuário: %v", err)
	}
	log.Printf("Usuário semeado/encontrado: %s", user.Email)

	for i := 0; i < 10; i++ {
		client, err := s.seedClient(user)
		if err != nil {
			log.Fatalf("Erro ao semear cliente: %v", err)
		}
		numRevenues := rand.Intn(5) + 1
	for j := 0; j < numRevenues; j++ {
		_, err := s.seedRevenue(user, client)
		if err != nil {
			log.Fatalf("Erro ao semear receita para o cliente %s: %v", client.Name, err)
		}
		
	}
	log.Println("10 clientes semeados com sucesso.")

	
}

}
func (s *Seeder) seedUser() (*model.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Email:                "dentista@email.com", 
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
		existingUser.ID, err = uuid.Parse("e6c45d98-e223-43de-a11c-de51cdfa6d1d")
		if err != nil{
			return nil, err
		}
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

// (NOVO) seedRevenue cria uma receita de teste para um cliente e usuário específicos.
func (s *Seeder) seedRevenue(user *model.User, client *model.Client) (*model.Revenue, error) {
// Lista de procedimentos para tornar os dados mais realistas
	procedures := []string{"Consulta de Rotina", "Limpeza e Profilaxia", "Restauração", "Extração de Siso", "Clareamento Dental"}

	// Gera um valor aleatório para a receita
	value := float64(rand.Intn(1400) + 100) // Valor entre 100 e 1500

	// Simula se a receita foi paga totalmente, parcialmente ou nada
	var totalPaid float64
	paymentStatus := rand.Intn(3) // Gera 0, 1, ou 2
	switch paymentStatus {
		case 0:
			totalPaid = value // Totalmente paga
		case 1:
			totalPaid = value * (rand.Float64()*0.5 + 0.2) // Parcialmente paga (20% a 70%)
		default:
			totalPaid = 0 // Não paga
		}

	revenue := model.Revenue{
		UserID:             user.ID,
		ClientID:           &client.ID, // Passa o ponteiro para o ID do cliente
		ProcedureType:      procedures[rand.Intn(len(procedures))],
		BeneficiaryCpfCnpj: client.Cpf,
		Value:              value,
		TotalPaid:          totalPaid,
		Description:        faker.Sentence(),
		IssueDate:          time.Now().AddDate(0, -rand.Intn(6), -rand.Intn(28)), // Data aleatória nos últimos 6 meses
		IsDeclared:         rand.Intn(2) == 1, // 50% de chance de ser true ou false
		}

	// Para idempotência, verificamos uma combinação de cliente e data para evitar duplicatas exatas.
	var existingRevenue model.Revenue
	err := s.DB.Where("client_id = ? AND issue_date = ? AND value = ?", client.ID, revenue.IssueDate.Format("2006-01-02"), revenue.Value).
	FirstOrInit(&existingRevenue, revenue).Error
		if err != nil {
		return nil, err
		}

		if existingRevenue.ID == uuid.Nil {
			existingRevenue.ID = uuid.New()
		if err := s.DB.Create(&existingRevenue).Error; err != nil {
			return nil, err
		}
		}

	return &existingRevenue, nil


}
