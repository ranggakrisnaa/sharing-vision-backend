package config

import (
    "log"
    "os"
    "strings"

    "github.com/joho/godotenv"
)

type Config struct {
    Port        string
    DatabaseURL string
    AppEnv      string
}

func Load() Config {
    _ = godotenv.Load()
    cfg := Config{
        Port:        getEnv("PORT", "8080"),
        DatabaseURL: getEnv("DATABASE_URL", ""),
        AppEnv:      getEnv("APP_ENV", "development"),
    }
    if strings.TrimSpace(cfg.DatabaseURL) == "" {
        log.Println("Warning: DATABASE_URL is empty")
    }
    return cfg
}

func getEnv(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}