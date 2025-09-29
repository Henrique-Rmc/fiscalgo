package middleware

import (
	"github.com/Henrique-Rmc/fiscalgo/apperror"
	"github.com/gofiber/fiber/v2"
	"log"
)

// ErrorHandlerMiddleware é usado no Fiber.Config
func ErrorHandlerMiddleware(c *fiber.Ctx, err error) error {
	// Se o erro for do tipo AppError, retorna padronizado
	if appErr, ok := err.(*apperror.AppError); ok {
		return c.Status(appErr.StatusCode).JSON(fiber.Map{
			"error":   appErr.Message,
			"details": appErr.Details,
		})
	}

	// Para erros não tratados, loga e retorna 500
	log.Printf("Erro inesperado: %v", err)
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "Erro interno inesperado.",
	})
}
