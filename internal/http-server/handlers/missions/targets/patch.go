package targets

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/illiakornyk/spy-cat/internal/utils"
)

type UpdateRequest struct {
	Notes    *string `json:"notes,omitempty" validate:"omitempty,max=500"`
	Complete *bool   `json:"complete,omitempty"`
}

type TargetUpdater interface {
	UpdateNotes(targetID int64, notes string) error
	UpdateCompleteStatus(targetID int64, complete bool) error
	MissionExists(missionID int64) (bool, error)
	TargetExists(targetID int64) (bool, error)
}
func UpdateTargetHandler(logger *slog.Logger, targetUpdater TargetUpdater) http.HandlerFunc {
	validate := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.missions.targets.update"
		logger = logger.With(slog.String("op", op))

		missionIDStr := chi.URLParam(r, "missionID")
		missionID, err := strconv.ParseInt(missionIDStr, 10, 64)
		if err != nil {
			logger.Error("invalid mission id", slog.Any("error", err))
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		exists, err := targetUpdater.MissionExists(missionID)
		if err != nil {
			logger.Error("failed to check if mission exists", slog.Any("error", err))
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		if !exists {
			logger.Error("mission does not exist", slog.Int64("missionID", missionID))
			utils.WriteError(w, http.StatusNotFound, errors.New("mission does not exist"))
			return
		}

		targetIDStr := chi.URLParam(r, "targetID")
		targetID, err := strconv.ParseInt(targetIDStr, 10, 64)
		if err != nil {
			logger.Error("invalid target id", slog.Any("error", err))
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		exists, err = targetUpdater.TargetExists(targetID)
		if err != nil {
			logger.Error("failed to check if target exists", slog.Any("error", err))
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		if !exists {
			logger.Error("target does not exist", slog.Int64("targetID", targetID))
			utils.WriteError(w, http.StatusNotFound, errors.New("target does not exist"))
			return
		}

		var req UpdateRequest
		err = utils.ParseJSON(r, &req)
		if errors.Is(err, io.EOF) {
			logger.Error("request body is empty")
			utils.WriteError(w, http.StatusBadRequest, errors.New("empty request"))
			return
		}
		if err != nil {
			logger.Error("failed to decode request body", slog.Any("error", err))
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		logger.Info("request body decoded", slog.Any("req", req))

		err = validate.Struct(req)
		if err != nil {
			logger.Error("validation failed", slog.Any("error", err))
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		if req.Notes != nil && req.Complete != nil {
			logger.Error("cannot update both notes and complete status at the same time")
			utils.WriteError(w, http.StatusBadRequest, errors.New("cannot update both notes and complete status at the same time"))
			return
		}

		if req.Notes != nil {
			updateNotes(w, r, targetID, *req.Notes, logger, targetUpdater)
		} else if req.Complete != nil {
			updateCompleteStatus(w, r, targetID, *req.Complete, logger, targetUpdater)
		}
	}
}

func updateNotes(w http.ResponseWriter, r *http.Request, targetID int64, notes string, logger *slog.Logger, targetUpdater TargetUpdater) {
	const op = "handlers.targets.updateNotes"
	logger = logger.With(slog.String("op", op))

	err := targetUpdater.UpdateNotes(targetID, notes)
	if err != nil {
		logger.Error("failed to update notes", slog.Any("error", err))
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	logger.Info("notes updated successfully", slog.Int64("id", targetID))
	w.WriteHeader(http.StatusNoContent)
}

func updateCompleteStatus(w http.ResponseWriter, r *http.Request, targetID int64, complete bool, logger *slog.Logger, targetUpdater TargetUpdater) {
	const op = "handlers.targets.updateCompleteStatus"
	logger = logger.With(slog.String("op", op))

	err := targetUpdater.UpdateCompleteStatus(targetID, complete)
	if err != nil {
		logger.Error("failed to update complete status", slog.Any("error", err))
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	logger.Info("complete status updated successfully", slog.Int64("id", targetID))
	w.WriteHeader(http.StatusNoContent)
}
