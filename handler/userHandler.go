package handler

import (
	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/service"
	"github.com/gofiber/fiber/v2"
)

type UserHandlerInterface interface {
	CreateUser(c *fiber.Ctx) error
}

type UserHandler struct {
	UserService service.UserServiceInterface
}

func NewUserHandler(userService service.UserServiceInterface) UserHandlerInterface {
	return &UserHandler{UserService: userService}
}

func (handler *UserHandler) CreateUser(c *fiber.Ctx) error {
	body := new(model.UserData)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "JSON invalido",
		})
	}
	user, err := handler.UserService.CreateUser(c.Context(), body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "JSON invalido",
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)

}
