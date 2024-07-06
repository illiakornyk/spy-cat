package main

import (
	"log/slog"
	"time"

	"github.com/illiakornyk/spy-cat/internal/breeds"
	"github.com/illiakornyk/spy-cat/internal/config"
	"github.com/illiakornyk/spy-cat/internal/http-server/router"
	"github.com/illiakornyk/spy-cat/internal/logger"
	"github.com/illiakornyk/spy-cat/internal/storage/initializer"
)

func main() {
	cfg := config.MustLoad()

	logger := logger.SetupLogger(cfg.Env)
	logger = logger.With(slog.String("env", cfg.Env))

	storage := initializer.InitializeStorage(cfg.StoragePath, logger)
	breeds.StartBreedCache(24 * time.Hour)

	r := router.SetupRouter(logger, storage)

	router.StartServer(cfg.HTTPServer.Address, r, logger)
}
