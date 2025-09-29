package apperror

import (
	"fmt"
	"net/http"
)

// OperationError é um wrapper para erros internos que podem ou não ser fatais.
type OperationError struct {
	OriginalErr error
	IsFatal     bool
}

func (e *OperationError) Error() string {
	if e.OriginalErr == nil {
		return ""
	}
	return e.OriginalErr.Error()
}

// AppError representa um erro de aplicação padronizado.
type AppError struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Err        error       `json:"-"`
	Details    interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// NotFound cria um erro para recurso não encontrado (HTTP 404).
func NotFound(resource string, err error) *AppError {
	return &AppError{
		StatusCode: http.StatusNotFound,
		Message:    fmt.Sprintf("%s não encontrado(a).", resource),
		Err:        err,
	}
}

// Forbidden cria um erro para acesso não permitido (HTTP 403).
func Forbidden(message string, err error) *AppError {
	if message == "" {
		message = "Você não tem permissão para executar esta ação."
	}
	return &AppError{
		StatusCode: http.StatusForbidden,
		Message:    message,
		Err:        err,
	}
}

// UnprocessableEntity cria um erro para dados inválidos segundo regras de negócio (HTTP 422).
func UnprocessableEntity(message string, err error) *AppError {
	return &AppError{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    message,
		Err:        err,
	}
}

// InternalServer cria um erro para falhas inesperadas (HTTP 500).
func InternalServer(message string, err error) *AppError {
	if message == "" {
		message = "Ocorreu um erro interno inesperado."
	}
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
		Err:        err,
	}
}
