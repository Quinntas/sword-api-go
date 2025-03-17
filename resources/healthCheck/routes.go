package healthCheck

import "github.com/gofiber/fiber/v2"

func SetupRoutes(router fiber.Router) {
	router.Get("/", Controller)
}
