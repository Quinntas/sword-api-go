package resources

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quinntas/go-fiber-template/resources/healthCheck"
	"github.com/quinntas/go-fiber-template/resources/task"
	"github.com/quinntas/go-fiber-template/resources/user"
)

func setupV1Router(app *fiber.App) {
	v1Router := app.Group("/api/v1")

	healthCheck.SetupRoutes(v1Router)
	user.SetupRoutes(v1Router)
	task.SetupRoutes(v1Router)
}

func SetupRouter(app *fiber.App) {
	setupV1Router(app)
}
