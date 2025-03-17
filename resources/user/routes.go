package user

import "github.com/gofiber/fiber/v2"

func SetupRoutes(router fiber.Router) {
	userRouter := router.Group("/users")

	userRouter.Post("/", CreateUser)
}
