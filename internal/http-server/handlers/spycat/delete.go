package spycat

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/illiakornyk/spy-cat/internal/utils"

	"log/slog"

	"github.com/go-chi/chi/v5"
)

type SpyCatDeleter interface {
    DeleteCat(id int64) error
	CatExists(id int64) (bool, error)
}

func DeleteHandler(logger *slog.Logger, spyCatDeleter SpyCatDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.spycat.delete"

		logger = logger.With(slog.String("op", op))

		idStr := chi.URLParam(r, "id")
		logger.Info("Extracted ID from URL", slog.String("idStr", idStr)) // Add log

		if idStr == "" {
			logger.Error("id path parameter is missing")
			utils.WriteError(w, http.StatusBadRequest, errors.New("id path parameter is missing"))
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id < 1 {
			logger.Error("invalid id path parameter", slog.Any("error", err))
			utils.WriteError(w, http.StatusBadRequest, errors.New("invalid id path parameter"))
			return
		}

		exists, err := spyCatDeleter.CatExists(id)
		if err != nil {
			logger.Error("failed to check if cat exists", slog.Any("error", err))
			utils.WriteError(w, http.StatusInternalServerError, errors.New("failed to check if cat exists"))
			return
		}

		if !exists {
			logger.Error("cat not found", slog.Int64("id", id))
			utils.WriteError(w, http.StatusNotFound, errors.New("cat not found"))
			return
		}

		err = spyCatDeleter.DeleteCat(id)
		if err != nil {
			logger.Error("failed to delete spy cat", slog.Any("error", err))
			utils.WriteError(w, http.StatusInternalServerError, errors.New("failed to delete spy cat"))
			return
		}

		logger.Info("spy cat deleted successfully", slog.Int64("id", id))
		w.WriteHeader(http.StatusNoContent)
	}
}
