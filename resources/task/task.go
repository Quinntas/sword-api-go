package task

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/quinntas/go-fiber-template/database"
	"github.com/quinntas/go-fiber-template/database/repository"
	"github.com/quinntas/go-fiber-template/eventEmitter"
	"github.com/quinntas/go-fiber-template/resources/user"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

const (
	StatusPending   = "PENDING"
	StatusCompleted = "COMPLETED"
)

func CreateTask(c *fiber.Ctx) error {
	type CreateTaskDTO struct {
		Summary string `json:"summary"`
	}

	taskDTO := new(CreateTaskDTO)
	if err := c.BodyParser(taskDTO); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "failed to parse body",
		})
	}

	authedUser, ok := c.Locals(user.LocalKey).(*repository.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	pid := uuid.New()

	_, err := database.Repo.CreateTask(context.Background(), repository.CreateTaskParams{
		Pid:          pid.String(),
		Summary:      taskDTO.Summary,
		Status:       StatusPending,
		TechnicianID: authedUser.ID,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Task created",
	})
}

func DeleteTask(c *fiber.Ctx) error {
	authedUser, ok := c.Locals(user.LocalKey).(*repository.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	if authedUser.Role != user.RoleManager {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	taskPid := c.Params("taskPid")
	if taskPid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid taskPid",
		})
	}

	task, err := database.Repo.GetTaskWithPid(context.Background(), taskPid)
	if err != nil {
		return err
	}

	err = database.Repo.DeleteTask(context.Background(), task.Pid)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task deleted",
	})
}

func UpdateTask(c *fiber.Ctx) error {
	authedUser, ok := c.Locals(user.LocalKey).(*repository.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	taskPid := c.Params("taskPid")
	if taskPid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid taskPid",
		})
	}

	task, err := database.Repo.GetTaskWithPid(context.Background(), taskPid)
	if err != nil {
		return err
	}

	if authedUser.ID != task.TechnicianID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	type UpdateTaskDTO struct {
		Done bool `json:"done"`
	}

	updateTaskDTO := new(UpdateTaskDTO)
	if err := c.BodyParser(updateTaskDTO); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "failed to parse body",
		})
	}

	if updateTaskDTO.Done {
		err := database.Repo.CompleteTask(context.Background(), task.Pid)
		if err != nil {
			return err
		}

		eventMessage := fmt.Sprintf("The tech with pid %s completed the task with pid %s at %s",
			authedUser.Pid,
			task.Pid,
			time.Now().Format(time.RFC3339),
		)

		err = eventEmitter.Manager.PublishMessage(
			eventEmitter.DefaultChannelName,
			OnTaskCompleteQueueName,
			[]byte(eventMessage),
		)
		if err != nil {
			return err
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task updated",
	})
}

func GetTasks(c *fiber.Ctx) error {
	authedUser, ok := c.Locals(user.LocalKey).(*repository.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	var data []repository.Task
	var err error

	if authedUser.Role == user.RoleManager {
		data, err = database.Repo.GetAllTasks(context.Background())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to fetch tasks",
			})
		}
	} else if authedUser.Role == user.RoleTechnician {
		data, err = database.Repo.GetTaskWithTechId(context.Background(), authedUser.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to fetch tasks",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": data,
	})
}

const (
	OnTaskCompleteQueueName = "OnTaskComplete"
)

func OnTaskComplete(msg amqp.Delivery) {
	fmt.Println(string(msg.Body))
}
