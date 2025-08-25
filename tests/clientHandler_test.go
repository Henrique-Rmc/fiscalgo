package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Henrique-Rmc/fiscalgo/handler"
	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClientService struct {
	mock.Mock
}

//Criamos um mock do Service com as funções de criar e resgatar user
/*
Cria um mock do service e recebe os mesmos parametros
cria uma variavel args que vai ser do tipo m(mockService) passando os dados que passamos
Nesse momento os dados sao atribuidos ao args(0)
Verificamos se algo realmente foi inserido no args(0) e se foi, acessamos o retorno esperado
*/
func (m *MockClientService) CreateClient(ctx context.Context, clientData *model.ClientData, idUser uuid.UUID) (*model.Client, error) {
	args := m.Called(ctx, clientData, idUser)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Client), args.Error(1)
}

func (m *MockClientService) GetById(ctx context.Context, clientId uuid.UUID, userId uuid.UUID) (*model.Client, error) {
	args := m.Called(ctx, clientId, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Client), args.Error(1)
}

func (m *MockClientService) FindClient(ctx context.Context, queryData *model.ClientSearchCriteria) ([]*model.Client, error) {
	args := m.Called(ctx, queryData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Client), args.Error(1)
}

func TestClientApi(t *testing.T) {
	app := fiber.New()
	mockClientService := new(MockClientService)
	clientHandler := handler.NewClientHandler(mockClientService)

	api := app.Group("api/clients")
	api.Post("/create", clientHandler.CreateClient)
	api.Get("/:clientId", clientHandler.GetCliendById)
	api.Get("/", clientHandler.FindClient)

	email := "joao.teste@email.com"
	userUuid, _ := uuid.Parse("6daa7ce0-6594-43ed-b583-c74bd6aa1a13")
	createClientData := model.ClientData{
		Name:        "joao Teste",
		Email:       &email,
		Cpf:         "11122233344",
		Phone:       "8877777777",
		AsksInvoice: true,
	}
	mockClient := &model.Client{
		ID:          uuid.New(),
		UserId:      userUuid,
		Name:        createClientData.Name,
		Email:       *createClientData.Email,
		Cpf:         createClientData.Cpf,
		Phone:       createClientData.Phone,
		AsksInvoice: createClientData.AsksInvoice,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	t.Run("Should Create Client Successfully", func(t *testing.T) {

		requestedBody, _ := json.Marshal(createClientData)
		mockClientService.On("CreateClient", mock.Anything, &createClientData, userUuid).Return(mockClient, nil).Once()
		req := httptest.NewRequest(http.MethodPost, "/api/clients/create", bytes.NewReader(requestedBody))
		req.Header.Set("Content-Type", "application/json")

		req = req.WithContext(context.WithValue(req.Context(), "loggedInUserID", userUuid))

		resp, err := app.Test(req)
		assert.NoError(t, err, "Should not Fail")
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Status Code Should be 201 created")
		bodyBytes, _ := io.ReadAll(resp.Body)
		var clientResponse model.Client
		json.Unmarshal(bodyBytes, &clientResponse)
		assert.Equal(t, mockClient.ID, clientResponse.ID, "The Client Id should be the same as the one returned by the mock")
		assert.Equal(t, mockClient.Name, clientResponse.Name, "The Client name should be the same as the one returned by the mock")
	})
	t.Run("Should Find Client Successfully", func(t *testing.T) {
		mockClientService.On("GetById", mock.Anything, mockClient.ID, userUuid).Return(mockClient, nil).Once()
		req := httptest.NewRequest(http.MethodGet, "/api/clients/"+mockClient.ID.String(), nil)

		req = req.WithContext(context.WithValue(req.Context(), "loggedInUserID", userUuid))

		resp, err := app.Test(req)
		assert.NoError(t, err, "Should not Fail")
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Status Code Should be 200 found")
		bodyBytes, _ := io.ReadAll(resp.Body)
		var clientResponse model.Client
		json.Unmarshal(bodyBytes, &clientResponse)

		assert.Equal(t, mockClient.ID, clientResponse.ID, "The User Id should be the same as the one returned by the mock")
		assert.Equal(t, "joao Teste", clientResponse.Name, "The User name should be the same as the one returned by the mock")
	})
	t.Run("Should Find Client by CPF Successfully", func(t *testing.T) {
		mockClientList := []*model.Client{mockClient}
		mockClientService.On("FindClient", mock.Anything, mock.AnythingOfType("*model.ClientSearchCriteria")).Return(mockClientList, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/clients?cpf="+mockClient.Cpf, nil)

		req = req.WithContext(context.WithValue(req.Context(), "loggedInUserID", userUuid))

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		bodyBytes, _ := io.ReadAll(resp.Body)
		var clientResponse []*model.Client
		json.Unmarshal(bodyBytes, &clientResponse)

		assert.NotEmpty(t, clientResponse)
		assert.Equal(t, mockClient.ID, clientResponse[0].ID)
	})

	t.Run("Should Find Client by Name Successfully", func(t *testing.T) {
		mockClientList := []*model.Client{mockClient}
		mockClientService.On("FindClient", mock.Anything, mock.AnythingOfType("*model.ClientSearchCriteria")).Return(mockClientList, nil).Once()
		encodedName := url.QueryEscape(mockClient.Name)
		req := httptest.NewRequest(http.MethodGet, "/api/clients?name="+encodedName, nil)
		req = req.WithContext(context.WithValue(req.Context(), "loggedInUserID", userUuid))

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		bodyBytes, _ := io.ReadAll(resp.Body)
		var clientResponse []*model.Client
		json.Unmarshal(bodyBytes, &clientResponse)

		assert.NotEmpty(t, clientResponse)
		assert.Equal(t, mockClient.ID, clientResponse[0].ID)
	})

	t.Run("Should Find Client by ID using query param Successfully", func(t *testing.T) {
		mockClientList := []*model.Client{mockClient}
		mockClientService.On("FindClient", mock.Anything, mock.AnythingOfType("*model.ClientSearchCriteria")).Return(mockClientList, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/clients?id="+mockClient.ID.String(), nil)
		req = req.WithContext(context.WithValue(req.Context(), "loggedInUserID", userUuid))

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		bodyBytes, _ := io.ReadAll(resp.Body)
		var clientResponse []*model.Client
		json.Unmarshal(bodyBytes, &clientResponse)

		assert.NotEmpty(t, clientResponse)
		assert.Equal(t, mockClient.ID, clientResponse[0].ID)
	})

}
