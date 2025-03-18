package task

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quinntas/go-fiber-template/eventEmitter"
	"github.com/quinntas/go-fiber-template/resources/user"
	"log"
)

func SetupRoutes(router fiber.Router) {
	taskRouter := router.Group("/tasks")

	taskRouter.Post("/", user.EnsureAuthenticated, CreateTask)
	taskRouter.Delete("/:taskPid", user.EnsureAuthenticated, DeleteTask)
	taskRouter.Put("/:taskPid", user.EnsureAuthenticated, UpdateTask)
	taskRouter.Get("/", user.EnsureAuthenticated, GetTasks)
}

func SetupEvents(manager *eventEmitter.ChannelManager) {
	err := manager.CreateQueueWithConsumer(eventEmitter.DefaultChannelName, OnTaskCompleteQueueName, OnTaskComplete)
	if err != nil {
		log.Fatal(err)
	}
}
