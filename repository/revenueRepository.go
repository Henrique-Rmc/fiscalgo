package repository

import (
	"context"
	"errors"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Erro personalizado para quando uma receita não é encontrada.
var ErrRevenueNotFound = errors.New("receita não encontrada")

// RevenueRepositoryInterface define o contrato para o repositório de receitas.
type RevenueRepositoryInterface interface {
	Create(ctx context.Context, revenue *model.Revenue) error
	FindByID(ctx context.Context, revenueID, userID uuid.UUID) (*model.Revenue, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Revenue, error)
	DeclareRevenue(ctx context.Context, userID uuid.UUID, revenueID uuid.UUID) error
	Update(ctx context.Context, revenue *model.Revenue) error
	Delete(ctx context.Context, revenueID, userID uuid.UUID) error
	Find(ctx context.Context, criteria *model.RevenueSearchCriteria) ([]*model.Revenue, error)
	// GetDeclaredSum(ctx context.Context, criteria *model.RevenueSearchCriteria) (float32, error)
}

// RevenueRepository é a implementação da interface.
type RevenueRepository struct {
	DB *gorm.DB
}

// NewRevenueRepository é o construtor do repositório.
func NewRevenueRepository(db *gorm.DB) RevenueRepositoryInterface {
	return &RevenueRepository{DB: db}
}

/*
-paciente inicia um pagamento
-Total -400 paciente paga 200
-Valor restante = 400
Quando uma receita é declarada, eu subtraio do debito
O valor pago é ZERADO
Para saber o valor pago, Value - Debit = Total Paid
*/

func (r *RevenueRepository) DeclareRevenue(ctx context.Context, userID uuid.UUID, revenueID uuid.UUID) error {
	result := r.DB.WithContext(ctx).
		Model(&model.Revenue{}).
		Where("user_id = ? AND id = ?", userID, revenueID).Update("is_declared = ?, total_paid = 0", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRevenueNotFound
	}
	return nil
}

// Create insere uma nova receita no banco de dados.
func (r *RevenueRepository) Create(ctx context.Context, revenue *model.Revenue) error {
	if err := r.DB.WithContext(ctx).Create(revenue).Error; err != nil {
		// Adicione aqui a lógica para tratar erros específicos do DB, se necessário.
		return err
	}
	return nil
}

// FindByID busca uma única receita pelo seu ID, garantindo que ela pertença ao usuário correto.
func (r *RevenueRepository) FindByID(ctx context.Context, revenueID, userID uuid.UUID) (*model.Revenue, error) {
	var revenue model.Revenue
	err := r.DB.WithContext(ctx).
		Where("id = ? AND user_id = ?", revenueID, userID).
		First(&revenue).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRevenueNotFound
		}
		return nil, err
	}
	return &revenue, nil
}

// FindAllByUserID busca todas as receitas de um usuário específico.
func (r *RevenueRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Revenue, error) {
	var revenues []*model.Revenue
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("issue_date DESC"). // Ordena as receitas da mais recente para a mais antiga.
		Find(&revenues).Error

	if err != nil {
		return nil, err
	}
	return revenues, nil
}

// Update salva as alterações de uma receita existente.
func (r *RevenueRepository) Update(ctx context.Context, revenue *model.Revenue) error {
	// O método Save atualiza todos os campos se o ID já existir.
	if err := r.DB.WithContext(ctx).Save(revenue).Error; err != nil {
		return err
	}
	return nil
}

// Delete remove uma receita do banco de dados, garantindo que o dono correto a está a apagar.
func (r *RevenueRepository) Delete(ctx context.Context, revenueID, userID uuid.UUID) error {
	result := r.DB.WithContext(ctx).
		Where("id = ? AND user_id = ?", revenueID, userID).
		Delete(&model.Revenue{})

	if result.Error != nil {
		return result.Error
	}

	// Verifique se alguma linha foi realmente afetada. Se não, a receita não existia.
	if result.RowsAffected == 0 {
		return ErrRevenueNotFound
	}

	return nil
}

/*
O sistema vai receber uma data de inicio e uma data de fim para filtrar a soma que busca

*/
// func (r *RevenueRepository)GetDeclaredSum(ctx context.Context, criteria *model.RevenueSearchCriteria) (float32, error){
// 	query := r.DB.WithContext(ctx).Where("user_id = ? AND is_declares = true", criteria.UserID)
// 	if criteria.StartDate != ""{
// 		query = query.Where("issue_date")
// 	}
// }

func (r *RevenueRepository) Find(ctx context.Context, criteria *model.RevenueSearchCriteria) ([]*model.Revenue, error) {
	var revenues []*model.Revenue

	// 1. Começa a consulta, sempre filtrando pelo dono (user_id). ESSENCIAL PARA SEGURANÇA.
	query := r.DB.WithContext(ctx).Where("user_id = ?", criteria.UserID)

	// 2. Adiciona os filtros opcionais condicionalmente.
	if criteria.ClientID != "" {
		// Valida e converte o ClientID de string para UUID
		clientUUID, err := uuid.Parse(criteria.ClientID)
		if err == nil { // Só aplica o filtro se for um UUID válido
			query = query.Where("client_id = ?", clientUUID)
		}
	}

	if criteria.ProcedureType != "" {
		// Usa ILIKE para uma busca de texto flexível e case-insensitive.
		query = query.Where("description ILIKE ?", "%"+criteria.ProcedureType+"%")
	}
	if criteria.OnlyInDebt {
		query = query.Where("total_paid < value")
	}
	if criteria.IsDeclared {
		query = query.Where("is_declared = ?", criteria.IsDeclared)
	}
	if criteria.StartDate != "" {
		// Adiciona um filtro de data "maior ou igual a".
		query = query.Where("issue_date >= ?", criteria.StartDate)
	}

	if criteria.EndDate != "" {
		query = query.Where("issue_date <= ?", criteria.EndDate)
	}

	if err := query.Order("issue_date DESC").Find(&revenues).Error; err != nil {
		return nil, err
	}

	return revenues, nil
}
