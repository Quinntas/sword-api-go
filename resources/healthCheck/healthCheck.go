package healthCheck

import "github.com/gofiber/fiber/v2"

func Controller(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ok",
	})
}
