package tests

import (
	"context"

	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/stretchr/testify/mock"
)

type MockClientRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateClient(ctx context.Context, ClientDto *model.ClientDto)
