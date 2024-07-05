package missions

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/illiakornyk/spy-cat/internal/common"
	"github.com/illiakornyk/spy-cat/internal/lib/api/response"
)

type GetAllResponse struct {
	response.Response
	Missions []common.Mission `json:"missions"`
}



type MissionGetter interface {
	GetAllMissions() ([]common.Mission, error)
}

func GetAllHandler(logger *slog.Logger, missionGetter MissionGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.missions.get_all"

		logger = logger.With(slog.String("op", op))

		missions, err := missionGetter.GetAllMissions()
		if err != nil {
			logger.Error("failed to get all missions", slog.Any("error", err))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "failed to get all missions",
			})
			return
		}

		logger.Info("retrieved all missions successfully", slog.Int("count", len(missions)))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GetAllResponse{
			Response: response.Response{
				Status: response.StatusOK,
			},
			Missions: missions,
		})
	}
}
