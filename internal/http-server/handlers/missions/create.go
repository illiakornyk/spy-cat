package missions

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/illiakornyk/spy-cat/internal/common"
	"github.com/illiakornyk/spy-cat/internal/lib/api/response"
)

type CreateRequest struct {
	CatID    *int64         `json:"cat_id,omitempty"` // Changed to pointer to make it optional
	Targets  []common.Target `json:"targets" validate:"required,dive"`
	Complete bool           `json:"complete"`
}

type CreateResponse struct {
	response.Response
	ID int64 `json:"id,omitempty"`
}

type MissionCreator interface {
	CreateMission(catID sql.NullInt64, targets []common.Target, complete bool) (int64, error)
}

func CreateHandler(logger *slog.Logger, missionCreator MissionCreator) http.HandlerFunc {
	validate := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.missions.create"

		logger = logger.With(slog.String("op", op))

		var req CreateRequest

		err := json.NewDecoder(r.Body).Decode(&req)
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

		var catID sql.NullInt64
		if req.CatID != nil {
			catID = sql.NullInt64{Int64: *req.CatID, Valid: true}
		} else {
			catID = sql.NullInt64{Valid: false}
		}

		id, err := missionCreator.CreateMission(catID, req.Targets, req.Complete)
		if err != nil {
			logger.Error("failed to create mission", slog.Any("error", err))

			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "failed to create mission",
			})
			return
		}

		logger.Info("mission created successfully", slog.Int64("id", id))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateResponse{
			Response: response.Response{
				Status: response.StatusOK,
			},
			ID: id,
		})
	}
}
