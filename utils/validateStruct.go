package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

//Valida se um objeto está correto com base nas refras definidas pelas tags no model

var validate = validator.New()

func ValidateStruct(s interface{}) error{
	err := validate.Struct(s)

	if err != nil{
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrors{

				return fmt.Errorf("Erro de validação no campo '%s': a regra '%s' falhou", fieldErr.Field(), fieldErr.Tag())
			}
		}
		return err
	}
	return nil
}