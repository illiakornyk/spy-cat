package missions

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/illiakornyk/spy-cat/internal/common"
	"github.com/illiakornyk/spy-cat/internal/utils"
)

type MissionGetter interface {
	GetMission(id int64) (*common.Mission, error)
}

func GetOneHandler(logger *slog.Logger, missionGetter MissionGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.missions.get"

		logger = logger.With(slog.String("op", op))

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			logger.Error("invalid mission id", slog.Any("error", err))
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid mission id"))
			return
		}

		mission, err := missionGetter.GetMission(id)
		if err != nil {
			logger.Error("failed to get mission", slog.Any("error", err))
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get mission"))
			return
		}
		if mission == nil {
			logger.Error("mission not found", slog.Int64("missionID", id))
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("mission not found"))
			return
		}

		logger.Info("mission retrieved successfully", slog.Int64("missionID", id))

		utils.WriteJSON(w, http.StatusOK, mission)
	}
}
