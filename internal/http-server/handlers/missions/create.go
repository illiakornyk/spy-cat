package missions

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/illiakornyk/spy-cat/internal/common"
	"github.com/illiakornyk/spy-cat/internal/utils"
)

type CreateRequest struct {
	CatID    *int64         `json:"cat_id,omitempty"`
	Targets  []common.Target `json:"targets" validate:"required,dive"`
	Complete bool           `json:"complete"`
}

type CreateResponse struct {
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

		err := utils.ParseJSON(r, &req)
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

		var catID sql.NullInt64
		if req.CatID != nil {
			catID = sql.NullInt64{Int64: *req.CatID, Valid: true}
		} else {
			catID = sql.NullInt64{Valid: false}
		}

		id, err := missionCreator.CreateMission(catID, req.Targets, req.Complete)
		if err != nil {
			logger.Error("failed to create mission", slog.Any("error", err))
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		logger.Info("mission created successfully", slog.Int64("id", id))

		utils.WriteJSON(w, http.StatusCreated, CreateResponse{ID: id})
	}
}
