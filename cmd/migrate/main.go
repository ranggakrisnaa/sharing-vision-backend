package main

import (
	"flag"
	"os"

	"database/sql"

	"github.com/ranggakrisnaa/sharing-vision-backend/pkg/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	logger.Init()

	dir := flag.String("dir", "migrations", "directory of migration files")
	action := flag.String("action", "up", "migration action: up or down")
	flag.Parse()

	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		logger.Log.Fatal("DATABASE_URL is empty")
	}

	logger.Log.WithFields(map[string]interface{}{"dir": *dir, "action": *action}).Info("starting migration")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Log.WithError(err).Fatal("open db")
	}
	defer db.Close()

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		logger.Log.WithError(err).Fatal("mysql driver init")
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+*dir, "mysql", driver)
	if err != nil {
		logger.Log.WithError(err).Fatal("migrate init")
	}

	switch *action {
	case "up":
		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				logger.Log.Info("migrate up: no change")
			} else {
				logger.Log.WithError(err).Fatal("migrate up failed")
			}
		} else {
			logger.Log.Info("migrate up: success")
		}
	case "down":
		if err := m.Down(); err != nil {
			if err == migrate.ErrNoChange {
				logger.Log.Info("migrate down: no change")
			} else {
				logger.Log.WithError(err).Fatal("migrate down failed")
			}
		} else {
			logger.Log.Info("migrate down: success")
		}
	default:
		logger.Log.WithField("action", *action).Fatal("unknown action")
	}

	logger.Log.Info("migration finished")
}
