package missions

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/illiakornyk/spy-cat/internal/lib/api/response"
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

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "invalid mission id",
			})
			return
		}

		exists, err := missionDeleter.MissionExists(id)
		if err != nil {
			logger.Error("failed to check if mission exists", slog.Any("error", err))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "failed to check if mission exists",
			})
			return
		}
		if !exists {
			logger.Error("mission does not exist", slog.Int64("missionID", id))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "mission does not exist",
			})
			return
		}

		err = missionDeleter.DeleteMission(id)
		if err != nil {
			logger.Error("failed to delete mission", slog.Any("error", err))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "failed to delete mission",
			})
			return
		}

		logger.Info("mission deleted successfully", slog.Int64("id", id))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.Response{
			Status: response.StatusOK,
		})
	}
}
