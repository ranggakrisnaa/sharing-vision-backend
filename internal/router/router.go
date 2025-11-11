package router

import (
	"time"

	"github.com/ranggakrisnaa/sharing-vision-backend/internal/article"
	"github.com/ranggakrisnaa/sharing-vision-backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Register(app *fiber.App, articleHandler *article.Handler) {
	// Global middlewares
	app.Use(recover.New())
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)
		logger.Log.WithFields(map[string]interface{}{
			"ip":         c.IP(),
			"method":     c.Method(),
			"path":       c.Path(),
			"status":     c.Response().StatusCode(),
			"latency_ms": latency.Milliseconds(),
		}).Info("http_request")
		return err
	})

	// check health
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Register article routes
	articleGroup := app.Group("/articles")
	articleHandler.Register(articleGroup)
}
