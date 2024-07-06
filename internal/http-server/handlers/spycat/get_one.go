package spycat

import (
	"fmt"
	"net/http"
	"strconv"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/illiakornyk/spy-cat/internal/utils"
)

type SpyCatGetter interface {
    GetCatByID(id int64) (*SpyCat, error)
	CatExists(id int64) (bool, error)

}

type GetOneResponse struct {
    Cat *SpyCat `json:"cat,omitempty"`
}

func GetOneHandler(logger *slog.Logger, spyCatGetter SpyCatGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.spycat.get_one"
		logger = logger.With(slog.String("op", op))

		idStr := chi.URLParam(r, "id")
		logger.Info("Extracted ID from URL", slog.String("idStr", idStr))

		if idStr == "" {
			logger.Error("id path parameter is missing")
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id path parameter is missing"))
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id < 1 {
			logger.Error("invalid id path parameter", slog.Any("error", err))
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id path parameter"))
			return
		}

		exists, err := spyCatGetter.CatExists(id)
		if err != nil {
			logger.Error("failed to check if cat exists", slog.Any("error", err))
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to check if cat exists"))
			return
		}

		if !exists {
			logger.Error("cat not found", slog.Int64("id", id))
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("cat not found"))
			return
		}

		cat, err := spyCatGetter.GetCatByID(id)
		if err != nil {
			logger.Error("failed to get spy cat by id", slog.Any("error", err))
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get spy cat by id"))
			return
		}

		logger.Info("retrieved spy cat successfully", slog.Int64("id", id))
		utils.WriteJSON(w, http.StatusOK, GetOneResponse{Cat: cat})
	}
}
