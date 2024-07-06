package targets

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/illiakornyk/spy-cat/internal/utils"
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
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		targetIDStr := chi.URLParam(r, "targetID")
		targetID, err := strconv.ParseInt(targetIDStr, 10, 64)
		if err != nil {
			logger.Error("invalid target id", slog.Any("error", err))
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		err = targetDeleter.DeleteTarget(targetID)
		if err != nil {
			logger.Error("failed to delete target", slog.Any("error", err))
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		logger.Info("target deleted successfully", slog.Int64("targetID", targetID))
		w.WriteHeader(http.StatusNoContent)
	}
}
