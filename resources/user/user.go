package user

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/quinntas/go-fiber-template/database"
	"github.com/quinntas/go-fiber-template/database/repository"
)

type CreateUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(c *fiber.Ctx) error {
	userDTO := new(CreateUserDTO)
	if err := c.BodyParser(userDTO); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "failed to parse body",
		})
	}
	createUser, err := database.Repo.CreateUser(context.Background(), repository.CreateUserParams{
		Email:    userDTO.Email,
		Password: userDTO.Password,
	})
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(createUser)
}
