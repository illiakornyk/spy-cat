package missions

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/illiakornyk/spy-cat/internal/common"
	"github.com/illiakornyk/spy-cat/internal/lib/api/response"
)

type MissionLister interface {
	GetAllMissions() ([]common.Mission, error)
}

func GetAllHandler(logger *slog.Logger, missionLister MissionLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.missions.list"

		logger = logger.With(slog.String("op", op))

		missions, err := missionLister.GetAllMissions()
		if err != nil {
			logger.Error("failed to list missions", slog.Any("error", err))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "failed to list missions",
			})
			return
		}

		logger.Info("missions listed successfully", slog.Int("count", len(missions)))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(missions)
	}
}
