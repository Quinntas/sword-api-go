package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"os"
	"time"
)

func setupMiddlewares(app *fiber.App) {
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(requestid.New())

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        20,
		Expiration: 30 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"message": "Too many requests",
			})
		},
	}))

	app.Get("/metrics", monitor.New())
}

func Create() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:       "go-fiber-template",
		ServerHeader:  "go-fiber-template",
		CaseSensitive: false,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		},
	})

	setupMiddlewares(app)

	return app
}

func Listen(app *fiber.App) error {
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "not found",
		})
	})

	serverHost := os.Getenv("HOST")
	serverPort := os.Getenv("PORT")

	return app.Listen(fmt.Sprintf("%s:%s", serverHost, serverPort))
}

/*
userRouter := v1Router.Group("/user")

	userRouter.Use(func(c *fiber.Ctx) error {
		c.Locals("user", "admin")
		return c.Next()
	})

	userRouter.Get("/", func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}
		fmt.Println(user)
		return c.JSON(fiber.Map{
			"message": "ok",
		})
	})
*/
