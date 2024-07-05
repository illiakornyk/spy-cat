package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/illiakornyk/spy-cat/internal/config"
	"github.com/illiakornyk/spy-cat/internal/storage/sqlite"
)

const (
    envLocal = "local"
    envDev   = "dev"
    envProd  = "prod"
)


func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)
    logger = logger.With(slog.String("env", cfg.Env))

    logger.Info("initializing server", slog.String("address", cfg.Address))
    logger.Debug("logger debug mode enabled")



	_, err := sqlite.New(cfg.StoragePath)
    if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
}



func setupLogger(env string) *slog.Logger {
    var log *slog.Logger

    switch env {
    case envLocal:
        log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    case envDev:
        log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    case envProd:
        log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
    }

    return log
}
