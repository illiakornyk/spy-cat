package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/illiakornyk/spy-cat/internal/breeds"
	"github.com/illiakornyk/spy-cat/internal/config"
	"github.com/illiakornyk/spy-cat/internal/http-server/handlers/spycat"
	mwLogger "github.com/illiakornyk/spy-cat/internal/http-server/middleware/logger"
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



	storage, err := sqlite.New(cfg.StoragePath)
    if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}

	breeds.StartBreedCache(24*time.Hour)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)


    router.Route("/api/v1/spy-cats", func(r chi.Router) {
		r.Get("/", spycat.GetAllHandler(logger, storage))
        r.Post("/", spycat.CreateHandler(logger, storage))
        r.Delete("/{id}", spycat.DeleteHandler(logger, storage))
        r.Patch("/{id}", spycat.PatchHandler(logger, storage))
		r.Get("/{id}", spycat.GetOneHandler(logger, storage))
    })

	log.Printf("Starting server at %s...", cfg.HTTPServer.Address)

	if err := http.ListenAndServe(cfg.HTTPServer.Address, router); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
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
