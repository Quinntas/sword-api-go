package user

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/quinntas/go-fiber-template/database"
	"github.com/quinntas/go-fiber-template/database/repository"
	"github.com/quinntas/go-fiber-template/utils/crypto"
	jsonwebtoken "github.com/quinntas/go-fiber-template/utils/jwt"
	"os"
	"time"
)

const (
	MANAGER    = "MANAGER"
	TECHNICIAN = "TECHNICIAN"
	LOCAL_KEY  = "USER_LOCAL_KEY"
)

func CreateUser(c *fiber.Ctx) error {
	type CreateUserDTO struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	userDTO := new(CreateUserDTO)
	if err := c.BodyParser(userDTO); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "failed to parse body",
		})
	}

	password, err := crypto.EncryptValue(userDTO.Password, os.Getenv("PEPPER"))
	if err != nil {
		return err
	}

	pid := uuid.New()

	_, err = database.Repo.CreateUser(context.Background(), repository.CreateUserParams{
		Username: userDTO.Username,
		Password: password,
		Pid:      pid.String(),
		Role:     TECHNICIAN,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created",
	})
}

func Login(c *fiber.Ctx) error {
	type LoginDTO struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	loginDTO := new(LoginDTO)
	if err := c.BodyParser(loginDTO); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "failed to parse body",
		})
	}

	user, err := database.Repo.GetUserWithUsername(context.Background(), loginDTO.Username)
	if err != nil {
		return err
	}

	if passwordMatches, err := crypto.CompareHash(loginDTO.Password, user.Password, os.Getenv("PEPPER")); err != nil || !passwordMatches {
		return err
	}

	expiresIn := time.Hour * 24 * 30 // 30 days
	expiresAt := time.Now().Add(expiresIn)

	jwt, err := jsonwebtoken.Sign[repository.User](user, expiresIn, os.Getenv("JWT_SECRET"))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token":     jwt,
		"expiresAt": expiresAt.Format(time.RFC3339),
		"expiresIn": expiresIn.Seconds(),
	})
}

func EnsureAuthenticated(c *fiber.Ctx) error {
	bearerToken := c.Get("Authorization")
	if bearerToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	token := bearerToken[7:] // Remove "Bearer " prefix

	authedUser, err := jsonwebtoken.Decode[repository.User](token, os.Getenv("JWT_SECRET"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	c.Locals(LOCAL_KEY, authedUser)

	return c.Next()
}
