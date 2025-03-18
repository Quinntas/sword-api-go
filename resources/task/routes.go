package task

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quinntas/go-fiber-template/resources/user"
)

func SetupRoutes(router fiber.Router) {
	taskRouter := router.Group("/tasks")

	taskRouter.Post("/", user.EnsureAuthenticated, CreateTask)
	taskRouter.Delete("/:taskPid", user.EnsureAuthenticated, DeleteTask)
	taskRouter.Put("/:taskPid", user.EnsureAuthenticated, UpdateTask)
	taskRouter.Get("/", user.EnsureAuthenticated, GetTasks)
}
