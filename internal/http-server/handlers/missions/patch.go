package missions

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/illiakornyk/spy-cat/internal/lib/api/response"
)

type UpdateCompleteStatusRequest struct {
	Complete bool `json:"complete"`
}

type MissionCompleter interface {
	UpdateMissionCompleteStatus(id int64, complete bool) error
}

func UpdateCompleteStatusHandler(logger *slog.Logger, missionCompleter MissionCompleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.missions.update_complete_status"

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

		var req UpdateCompleteStatusRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Error("failed to decode request body", slog.Any("error", err))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "failed to decode request",
			})
			return
		}

		err = missionCompleter.UpdateMissionCompleteStatus(id, req.Complete)
		if err != nil {
			logger.Error("failed to update mission complete status", slog.Any("error", err))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "failed to update mission complete status",
			})
			return
		}

		logger.Info("mission complete status updated successfully", slog.Int64("id", id), slog.Bool("complete", req.Complete))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.Response{
			Status: response.StatusOK,
		})
	}
}
