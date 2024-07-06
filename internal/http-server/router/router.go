package router

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/illiakornyk/spy-cat/internal/http-server/handlers/missions"
	"github.com/illiakornyk/spy-cat/internal/http-server/handlers/missions/targets"
	"github.com/illiakornyk/spy-cat/internal/http-server/handlers/spycat"
	mwLogger "github.com/illiakornyk/spy-cat/internal/http-server/middleware/logger"
	"github.com/illiakornyk/spy-cat/internal/storage/sqlite"

	"github.com/go-chi/chi/middleware"
)

func SetupRouter(logger *slog.Logger, storage *sqlite.Storage) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	setupRoutes(router, logger, storage)

	return router
}

func setupRoutes(router *chi.Mux, logger *slog.Logger, storage *sqlite.Storage) {
	router.Route("/api/v1/spy-cats", func(r chi.Router) {
		r.Get("/", spycat.GetAllHandler(logger, storage))
		r.Post("/", spycat.CreateHandler(logger, storage))
		r.Delete("/{id}", spycat.DeleteHandler(logger, storage))
		r.Patch("/{id}", spycat.PatchHandler(logger, storage))
		r.Get("/{id}", spycat.GetOneHandler(logger, storage))
	})

	router.Route("/api/v1/missions", func(r chi.Router) {
		r.Post("/", missions.CreateHandler(logger, storage))
		r.Get("/", missions.GetAllHandler(logger, storage))
		r.Get("/{id}", missions.GetOneHandler(logger, storage))
		r.Patch("/{id}", missions.UpdateHandler(logger, storage))
		r.Delete("/{id}", missions.DeleteHandler(logger, storage))

		// Target routes
		r.Route("/{missionID}/targets", func(r chi.Router) {
			r.Patch("/{targetID}", targets.UpdateTargetHandler(logger, storage))
			r.Delete("/{targetID}", targets.DeleteTargetHandler(logger, storage))
			r.Post("/", targets.AddTargetHandler(logger, storage))
		})
	})
}

func StartServer(address string, router *chi.Mux, logger *slog.Logger) {
	logger.Info("Starting server", slog.String("address", address))
	if err := http.ListenAndServe(address, router); err != nil {
		logger.Error("Failed to start server", slog.Any("error", err))
	}
}
