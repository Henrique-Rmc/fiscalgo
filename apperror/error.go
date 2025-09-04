package apperror

import (
	"fmt"
	"net/http"
)

type OperationError struct {
	OriginalErr error
	IsFatal     bool
}

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *OperationError) Error() string {
	if e.OriginalErr == nil {
		return ""
	}
	return e.OriginalErr.Error()
}

func NotFound(resource string, err error) *AppError {
	return &AppError{
		Code:    http.StatusNotFound, // 404
		Message: fmt.Sprintf("%s não encontrado(a).", resource),
		Err:     err,
	}
}

// Forbidden cria um erro para acesso não permitido (mapeia para HTTP 403).
func Forbidden(message string, err error) *AppError {
	if message == "" {
		message = "Você não tem permissão para executar esta ação."
	}
	return &AppError{
		Code:    http.StatusForbidden, // 403
		Message: message,
		Err:     err,
	}
}

// UnprocessableEntity cria um erro para dados que são semanticamente inválidos (mapeia para HTTP 422).
// A sintaxe está correta, mas as regras de negócio foram violadas.
func UnprocessableEntity(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusUnprocessableEntity, // 422
		Message: message,
		Err:     err,
	}
}

// InternalServer cria um erro para falhas inesperadas (mapeia para HTTP 500).
func InternalServer(message string, err error) *AppError {
	if message == "" {
		message = "Ocorreu um erro interno inesperado."
	}
	return &AppError{
		Code:    http.StatusInternalServerError, // 500
		Message: message,
		Err:     err,
	}
}