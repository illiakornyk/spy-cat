package targets

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/illiakornyk/spy-cat/internal/lib/api/response"
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

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "invalid mission id",
			})
			return
		}

		exists, err := targetUpdater.MissionExists(missionID)
		if err != nil {
			logger.Error("failed to check if mission exists", slog.Any("error", err))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "failed to check if mission exists",
			})
			return
		}
		if !exists {
			logger.Error("mission does not exist", slog.Int64("missionID", missionID))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "mission does not exist",
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

		exists, err = targetUpdater.TargetExists(targetID)
		if err != nil {
			logger.Error("failed to check if target exists", slog.Any("error", err))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "failed to check if target exists",
			})
			return
		}
		if !exists {
			logger.Error("target does not exist", slog.Int64("targetID", targetID))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "target does not exist",
			})
			return
		}

		var req UpdateRequest

		err = json.NewDecoder(r.Body).Decode(&req)
		if errors.Is(err, io.EOF) {
			logger.Error("request body is empty")

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "empty request",
			})
			return
		}
		if err != nil {
			logger.Error("failed to decode request body", slog.Any("error", err))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "failed to decode request",
			})
			return
		}

		logger.Info("request body decoded", slog.Any("req", req))

		err = validate.Struct(req)
		if err != nil {
			logger.Error("validation failed", slog.Any("error", err))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "validation failed",
			})
			return
		}

		if req.Notes != nil {
			err = targetUpdater.UpdateNotes(targetID, *req.Notes)
			if err != nil {
				logger.Error("failed to update notes", slog.Any("error", err))

				json.NewEncoder(w).Encode(response.Response{
					Status: response.StatusError,
					Error:  "failed to update notes",
				})
				return
			}
			logger.Info("notes updated successfully", slog.Int64("id", targetID))
		}

		if req.Complete != nil {
			err = targetUpdater.UpdateCompleteStatus(targetID, *req.Complete)
			if err != nil {
				logger.Error("failed to update complete status", slog.Any("error", err))

				json.NewEncoder(w).Encode(response.Response{
					Status: response.StatusError,
					Error:  "failed to update complete status",
				})
				return
			}
			logger.Info("complete status updated successfully", slog.Int64("id", targetID))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.Response{
			Status: response.StatusOK,
		})
	}
}
