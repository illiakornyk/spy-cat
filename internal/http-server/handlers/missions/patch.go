package missions

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/illiakornyk/spy-cat/internal/lib/api/response"
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

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "invalid mission id",
			})
			return
		}

		exists, err := missionUpdater.MissionExists(id)
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

		var req UpdateMissionRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Error("failed to decode request body", slog.Any("error", err))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "failed to decode request",
			})
			return
		}

		if req.Complete == nil && req.CatID == nil {
			logger.Error("no update fields provided")

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "no update fields provided",
			})
			return
		}

		if req.Complete != nil && req.CatID != nil {
			logger.Error("cannot update both complete status and cat_id at the same time")

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "cannot update both complete status and cat_id at the same time",
			})
			return
		}

		if req.Complete != nil {
			err = missionUpdater.UpdateMissionCompleteStatus(id, *req.Complete)
			if err != nil {
				logger.Error("failed to update mission complete status", slog.Any("error", err))

				json.NewEncoder(w).Encode(response.Response{
					Status: response.StatusError,
					Error:  "failed to update mission complete status",
				})
				return
			}
			logger.Info("mission complete status updated successfully", slog.Int64("id", id), slog.Bool("complete", *req.Complete))
		}

		if req.CatID != nil {
			err = validate.Var(req.CatID, "required,min=1")
			if err != nil {
				logger.Error("validation failed", slog.Any("error", err))

				json.NewEncoder(w).Encode(response.Response{
					Status: response.StatusError,
					Error:  "validation failed",
				})
				return
			}

			err = missionUpdater.AssignCatToMission(id, *req.CatID)
			if err != nil {
				logger.Error("failed to assign cat to mission", slog.Any("error", err))

				json.NewEncoder(w).Encode(response.Response{
					Status: response.StatusError,
					Error:  "failed to assign cat to mission",
				})
				return
			}
			logger.Info("cat assigned to mission successfully", slog.Int64("missionID", id), slog.Int64("catID", *req.CatID))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.Response{
			Status: response.StatusOK,
		})
	}
}
