package main

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ranggakrisnaa/sharing-vision-backend/internal/article"
	"github.com/ranggakrisnaa/sharing-vision-backend/internal/router"
	"github.com/ranggakrisnaa/sharing-vision-backend/pkg/config"
	"github.com/ranggakrisnaa/sharing-vision-backend/pkg/database"
	"github.com/ranggakrisnaa/sharing-vision-backend/pkg/logger"
	validatorpkg "github.com/ranggakrisnaa/sharing-vision-backend/pkg/validator"
)

func main() {
	cfg := config.Load()
	logger.Init()

	db, err := database.NewMySQL(cfg.DatabaseURL)
	if err != nil {
		logger.Log.WithError(err).Fatal("failed connect DB")
	}
	defer db.Close()

	app := fiber.New()

	// Register routes
	articleRepository := article.NewMySQLRepository(db)
	articleService := article.NewService(articleRepository)
	articleHandler := article.NewHandler(articleService, validatorpkg.NewValidator())
	router.Register(app, articleHandler)

	port := cfg.Port
	logger.Log.WithField("port", port).Info("server listening")
	if err := app.Listen(":" + port); err != nil {
		logger.Log.WithError(err).Fatal("server crashed")
	}
}
