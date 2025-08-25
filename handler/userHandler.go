package handler

import (
	"github.com/Henrique-Rmc/fiscalgo/model"
	"github.com/Henrique-Rmc/fiscalgo/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandlerInterface interface {
	CreateUser(c *fiber.Ctx) error
	GetUserById(c *fiber.Ctx) error
}

type UserHandler struct {
	UserService service.UserServiceInterface
}

func NewUserHandler(userService service.UserServiceInterface) UserHandlerInterface {
	return &UserHandler{UserService: userService}
}

func (userHandler *UserHandler) CreateUser(c *fiber.Ctx) error {
	body := new(model.UserData)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "JSON invalido",
		})
	}
	user, err := userHandler.UserService.CreateUser(c.Context(), body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "JSON invalido",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(user)

}

func (userHandler *UserHandler) GetUserById(c *fiber.Ctx) error {
	//Capturar o user Id vindo dos params
	userId := c.Params("userId")
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Params",
		})
	}
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error when converting params to UUID",
		})
	}
	user, err := userHandler.UserService.GetUserById(c.Context(), userUUID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Coundn't find user with corresponding ID",
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}
