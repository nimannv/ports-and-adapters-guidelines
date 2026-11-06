package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// NewServer builds and configures a Fiber HTTP application with registered route handlers.
func NewServer(
	boardHandler *BoardHandler,
	taskHandler *TaskHandler,
) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Simple Task Management API",
	})

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New())

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})

	apiV1 := app.Group("/api/v1")
	boardHandler.Register(apiV1)
	taskHandler.Register(apiV1)

	return app
}

