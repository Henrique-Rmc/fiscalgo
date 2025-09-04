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
	Update(ctx context.Context, revenue *model.Revenue) error
	Delete(ctx context.Context, revenueID, userID uuid.UUID) error
	Find(ctx context.Context, criteria model.RevenueSearchCriteria) ([]*model.Revenue, error)
}

// RevenueRepository é a implementação da interface.
type RevenueRepository struct {
	DB *gorm.DB
}

// NewRevenueRepository é o construtor do repositório.
func NewRevenueRepository(db *gorm.DB) RevenueRepositoryInterface {
	return &RevenueRepository{DB: db}
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

func (r *RevenueRepository) Find(ctx context.Context, criteria model.RevenueSearchCriteria) ([]*model.Revenue, error) {
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

	if criteria.StartDate != "" {
		// Adiciona um filtro de data "maior ou igual a".
		query = query.Where("issue_date >= ?", criteria.StartDate)
	}

	if criteria.EndDate != "" {
		// Adiciona um filtro de data "menor ou igual a".
		query = query.Where("issue_date <= ?", criteria.EndDate)
	}

	// 3. Executa a consulta final e preenche o slice 'revenues'.
	//    Adicionamos uma ordenação padrão para resultados consistentes.
	if err := query.Order("issue_date DESC").Find(&revenues).Error; err != nil {
		return nil, err
	}

	// 4. Retorna o slice de receitas (pode estar vazio, não é um erro).
	return revenues, nil
}