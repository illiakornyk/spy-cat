package missions

import (
	"log/slog"
	"net/http"

	"github.com/illiakornyk/spy-cat/internal/common"
	"github.com/illiakornyk/spy-cat/internal/utils"
)

type MissionLister interface {
	GetAllMissions() ([]common.Mission, error)
}

type MissionResponse struct {
			ID       int64             `json:"id"`
			CatID    *int64            `json:"cat_id,omitempty"`
			Complete bool              `json:"complete"`
			Targets  []common.Target   `json:"targets"`
		}

func GetAllHandler(logger *slog.Logger, missionLister MissionLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.missions.list"

		logger = logger.With(slog.String("op", op))

		missions, err := missionLister.GetAllMissions()
		if err != nil {
			logger.Error("failed to list missions", slog.Any("error", err))
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		var response []MissionResponse
		for _, mission := range missions {
			var catID *int64
			if mission.CatID.Valid {
				catID = &mission.CatID.Int64
			}
			response = append(response, MissionResponse{
				ID:       mission.ID,
				CatID:    catID,
				Complete: mission.Complete,
				Targets:  mission.Targets,
			})
		}

		logger.Info("missions listed successfully", slog.Int("count", len(missions)))

		utils.WriteJSON(w, http.StatusOK, response)
	}
}
