package missions

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/illiakornyk/spy-cat/internal/utils"
)

type UpdateMissionRequest struct {
	Complete *bool `json:"complete,omitempty"`
	CatID    *int64 `json:"cat_id,omitempty"`
}

type MissionUpdater interface {
	UpdateMissionCompleteStatus(id int64, complete bool) error
	AssignCatToMission(missionID, catID int64) error
	MissionExists(id int64) (bool, error)
}
func UpdateHandler(logger *slog.Logger, missionUpdater MissionUpdater) http.HandlerFunc {
	validate := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.missions.update"
		logger = logger.With(slog.String("op", op))

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			logger.Error("invalid mission id", slog.Any("error", err))
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid mission id"))
			return
		}

		exists, err := missionUpdater.MissionExists(id)
		if err != nil {
			logger.Error("failed to check if mission exists", slog.Any("error", err))
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to check if mission exists"))
			return
		}
		if !exists {
			logger.Error("mission does not exist", slog.Int64("missionID", id))
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("mission does not exist"))
			return
		}

		var req UpdateMissionRequest
		err = utils.ParseJSON(r, &req)
		if err != nil {
			logger.Error("failed to decode request body", slog.Any("error", err))
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to decode request"))
			return
		}

		if req.Complete == nil && req.CatID == nil {
			logger.Error("no update fields provided")
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("no update fields provided"))
			return
		}

		if req.Complete != nil && req.CatID != nil {
			logger.Error("cannot update both complete status and cat_id at the same time")
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("cannot update both complete status and cat_id at the same time"))
			return
		}

		if req.Complete != nil {
			updateCompleteStatus(w, r, id, *req.Complete, logger, missionUpdater)
		} else if req.CatID != nil {
			assignCat(w, r, id, *req.CatID, logger, missionUpdater, validate)
		}
	}
}

func updateCompleteStatus(w http.ResponseWriter, r *http.Request, id int64, complete bool, logger *slog.Logger, missionUpdater MissionUpdater) {
	const op = "handlers.missions.updateCompleteStatus"
	logger = logger.With(slog.String("op", op))

	err := missionUpdater.UpdateMissionCompleteStatus(id, complete)
	if err != nil {
		logger.Error("failed to update mission complete status", slog.Any("error", err))
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to update mission complete status"))
		return
	}

	logger.Info("mission complete status updated successfully", slog.Int64("id", id), slog.Bool("complete", complete))
	w.WriteHeader(http.StatusNoContent)
}

func assignCat(w http.ResponseWriter, r *http.Request, id int64, catID int64, logger *slog.Logger, missionUpdater MissionUpdater, validate *validator.Validate) {
	const op = "handlers.missions.assignCat"
	logger = logger.With(slog.String("op", op))

	err := validate.Var(catID, "required,min=1")
	if err != nil {
		logger.Error("validation failed", slog.Any("error", err))
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("validation failed: %w", err))
		return
	}

	err = missionUpdater.AssignCatToMission(id, catID)
	if err != nil {
		logger.Error("failed to assign cat to mission", slog.Any("error", err))
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to assign cat to mission"))
		return
	}

	logger.Info("cat assigned to mission successfully", slog.Int64("missionID", id), slog.Int64("catID", catID))
	w.WriteHeader(http.StatusNoContent)
}
