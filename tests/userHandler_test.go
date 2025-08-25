package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Henrique-Rmc/fiscalgo/handler"
	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

//Criamos um mock do Service com as funções de criar e resgatar user
/*
Cria um mock do service e recebe os mesmos parametros
cria uma variavel args que vai ser do tipo m(mockService) passando os dados que passamos
Nesse momento os dados sao atribuidos ao args(0)
Verificamos se algo realmente foi inserido no args(0) e se foi, acessamos o retorno esperado
*/
func (m *MockUserService) CreateUser(ctx context.Context, data *model.UserData) (*model.User, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) GetUserById(ctx context.Context, userId uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

/*
Meu Create User vai receber um body simples e processar seus dados normalmente
O Gemini sugeriu que eu devo criar um model com os dados que pretendo inserir com o Handler
O model é um objeto e o que o json.Marshal faz é converter uma estrutura de dados em um objeto json
*/

// Devemos criar a simulação da aplicação desde o inicio para inicializar o handler
func TestUserApi(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUserService)
	userHandler := handler.NewUserHandler(mockService)

	api := app.Group("api/users")
	api.Post("/create", userHandler.CreateUser)
	api.Get("/:userId", userHandler.GetUserById)

	t.Run("Should Create User Successfully", func(t *testing.T) {
		createUserData := model.UserData{
			Name:                 "João Teste",
			Email:                "joao.teste@email.com",
			CPF:                  "11122233344",
			Password:             "senha123",
			Occupation:           "Tester",
			ProfessionalRegistry: "TEST-123",
		}
		requestedBody, _ := json.Marshal(createUserData)

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
		//Uma função ON do mockService, essa função recebe uma string que parece ser
		//Ele diz que deve retornar um mockUser ao passar um createUserData na função CreateUser
		mockService.On("CreateUser", mock.Anything, &createUserData).Return(mockUser, nil).Once()
		//Vai criar uma nova requisição e definir o header como um json
		req := httptest.NewRequest(http.MethodPost, "/api/users/create", bytes.NewReader(requestedBody))
		req.Header.Set("Content-Type", "application/json")

		//Captura a resposta
		resp, err := app.Test(req)
		//Valida que foi feita uma requisição com os status esperados
		assert.NoError(t, err, "Should not Fail")
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Status Code Should be 201 created")
		//Recria o objeto recebido e o compara com o que esperamos
		bodyBytes, _ := io.ReadAll(resp.Body)
		var userResponse model.User
		json.Unmarshal(bodyBytes, &userResponse)

		assert.Equal(t, mockUser.ID, userResponse.ID, "The User Id should be the same as the one returned by the mock")
		assert.Equal(t, mockUser.Name, userResponse.Name, "The User name should be the same as the one returned by the mock")
		createdUserId := userResponse.ID

		t.Run("Should Find User Successfully", func(t *testing.T) {
			mockService.On("GetUserById", mock.Anything, createdUserId).Return(mockUser, nil).Once()
			req := httptest.NewRequest(http.MethodGet, "/api/users/"+createdUserId.String(), nil)
			resp, err := app.Test(req)
			assert.NoError(t, err, "Should not Fail")
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Status Code Should be 200 found")
			bodyBytes, _ := io.ReadAll(resp.Body)
			var userResponse model.User
			json.Unmarshal(bodyBytes, &userResponse)

			assert.Equal(t, createdUserId, userResponse.ID, "The User Id should be the same as the one returned by the mock")
			assert.Equal(t, "João Teste", userResponse.Name, "The User name should be the same as the one returned by the mock")
		})
	})

}
