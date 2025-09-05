package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	// Retornamos o que foi programado no teste
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindUserById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// --- 2. Funções de Teste ---

func TestUserService(t *testing.T) {
	t.Run("CreateUser", func(t *testing.T) {
		t.Run("should create user successfully", func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			userService := service.NewUserService(mockRepo, nil) 

			createUserData := &model.UserDto{
				Name:                 "João Teste",
				Email:                "joao.teste@email.com",
				CPF:                  "11122233344",
				Password:             "senha123",
				Occupation:           "Tester",
				ProfessionalRegistry: "TEST-123",
			}

			mockUser := &model.User{
			ID:                   uuid.New(),
			Name:                 createUserData.Name,
			Email:                createUserData.Email,
			CPF:                  createUserData.CPF,
			PasswordHash:         "senha123",
			Occupation:           createUserData.Occupation,
			ProfessionalRegistry: createUserData.ProfessionalRegistry,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}

			mockRepo.On(
				"CreateUser",
				mock.Anything,
				mock.MatchedBy(func(user *model.User) bool {
					assert.NotEmpty(t, user.ID, "O ID deve ser gerado pelo serviço")
					assert.NotEqual(t, createUserData.Password, user.PasswordHash, "A senha deve ser encriptada")
					return user.Name == createUserData.Name 
				}),
			).Return(mockUser,nil).Once()

			createdUser, err := userService.CreateUser(context.Background(), createUserData)

			assert.NoError(t, err)
			assert.NotNil(t, createdUser)
			assert.Equal(t, createUserData.Email, createdUser.Email)
			mockRepo.AssertExpectations(t) 
		})

		t.Run("should return error when repository fails", func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			userService := service.NewUserService(mockRepo, nil)
			userData := &model.UserDto{Password: "123"}
			expectedError := errors.New("database error")

			mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User")).
				Return(nil, expectedError).Once()

			createdUser, err := userService.CreateUser(context.Background(), userData)

			assert.Error(t, err)
			assert.Nil(t, createdUser)
			assert.Equal(t, expectedError, err)
			mockRepo.AssertExpectations(t)
		})
	})

	t.Run("GetUserById", func(t *testing.T) {
		t.Run("should get user by id successfully", func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			userService := service.NewUserService(mockRepo, nil)

			testID := uuid.New()
			mockUser := &model.User{ID: testID, Name: "João Encontrado"}

			mockRepo.On("FindUserById", mock.Anything, testID).Return(mockUser, nil).Once()

			foundUser, err := userService.GetUserById(context.Background(), testID)

			assert.NoError(t, err)
			assert.NotNil(t, foundUser)
			assert.Equal(t, testID, foundUser.ID)
			mockRepo.AssertExpectations(t)
		})

		t.Run("should return not found error", func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			userService := service.NewUserService(mockRepo, nil)
			testID := uuid.New()

			mockRepo.On("FindUserById", mock.Anything, testID).Return(nil, gorm.ErrRecordNotFound).Once()

			foundUser, err := userService.GetUserById(context.Background(), testID)

			assert.Error(t, err)
			assert.Nil(t, foundUser)
			assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
			mockRepo.AssertExpectations(t)
		})
	})

}
