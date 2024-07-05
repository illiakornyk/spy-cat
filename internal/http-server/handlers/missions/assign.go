package missions

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

type AssignCatRequest struct {
	CatID int64 `json:"cat_id" validate:"required,min=1"`
}

type MissionAssigner interface {
	AssignCatToMission(missionID, catID int64) error
}

func AssignCatHandler(logger *slog.Logger, missionAssigner MissionAssigner) http.HandlerFunc {
	validate := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.missions.assign_cat"

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

		var req AssignCatRequest

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

		err = missionAssigner.AssignCatToMission(missionID, req.CatID)
		if err != nil {
			logger.Error("failed to assign cat to mission", slog.Any("error", err))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "failed to assign cat to mission",
			})
			return
		}

		logger.Info("cat assigned to mission successfully", slog.Int64("missionID", missionID))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.Response{
			Status: response.StatusOK,
		})
	}
}
