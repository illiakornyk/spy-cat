package missions

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/illiakornyk/spy-cat/internal/utils"
)

type MissionDeleter interface {
	DeleteMission(id int64) error
	MissionExists(missionID int64) (bool, error)
}

func DeleteHandler(logger *slog.Logger, missionDeleter MissionDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.missions.delete"
		logger = logger.With(slog.String("op", op))

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			logger.Error("invalid mission id", slog.Any("error", err))
			utils.WriteError(w, http.StatusBadRequest, errors.New("invalid mission id"))
			return
		}

		exists, err := missionDeleter.MissionExists(id)
		if err != nil {
			logger.Error("failed to check if mission exists", slog.Any("error", err))
			utils.WriteError(w, http.StatusInternalServerError, errors.New("failed to check if mission exists"))
			return
		}
		if !exists {
			logger.Error("mission does not exist", slog.Int64("missionID", id))
			utils.WriteError(w, http.StatusNotFound, errors.New("mission does not exist"))
			return
		}

		err = missionDeleter.DeleteMission(id)
		if err != nil {
			logger.Error("failed to delete mission", slog.Any("error", err))
			utils.WriteError(w, http.StatusInternalServerError, errors.New("failed to delete mission"))
			return
		}

		logger.Info("mission deleted successfully", slog.Int64("id", id))
		w.WriteHeader(http.StatusNoContent)
	}
}
