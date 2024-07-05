package targets

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/illiakornyk/spy-cat/internal/lib/api/response"
)

type TargetDeleter interface {
	DeleteTarget(targetID int64) error
}

func DeleteTargetHandler(logger *slog.Logger, targetDeleter TargetDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.missions.targets.delete"

		logger = logger.With(slog.String("op", op))

		missionIDStr := chi.URLParam(r, "missionID")
		_, err := strconv.ParseInt(missionIDStr, 10, 64)
		if err != nil {
			logger.Error("invalid mission id", slog.Any("error", err))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "invalid mission id",
			})
			return
		}

		targetIDStr := chi.URLParam(r, "targetID")
		targetID, err := strconv.ParseInt(targetIDStr, 10, 64)
		if err != nil {
			logger.Error("invalid target id", slog.Any("error", err))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "invalid target id",
			})
			return
		}

		err = targetDeleter.DeleteTarget(targetID)
		if err != nil {
			logger.Error("failed to delete target", slog.Any("error", err))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  err.Error(),
			})
			return
		}

		logger.Info("target deleted successfully", slog.Int64("targetID", targetID))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.Response{
			Status: response.StatusOK,
		})
	}
}
