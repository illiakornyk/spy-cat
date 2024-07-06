package targets

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/illiakornyk/spy-cat/internal/utils"
)

type AddTargetRequest struct {
	Name    string `json:"name" validate:"required,min=1,max=100"`
	Country string `json:"country" validate:"required,min=1,max=100"`
	Notes   string `json:"notes" validate:"omitempty,max=500"`
}

type TargetAdder interface {
	AddTarget(missionID int64, name, country, notes string) (int64, error)
	MissionExists(missionID int64) (bool, error)
}

func AddTargetHandler(logger *slog.Logger, targetAdder TargetAdder) http.HandlerFunc {
	validate := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.missions.targets.add"
		logger = logger.With(slog.String("op", op))

		missionIDStr := chi.URLParam(r, "missionID")
		missionID, err := strconv.ParseInt(missionIDStr, 10, 64)
		if err != nil {
			logger.Error("invalid mission id", slog.Any("error", err))
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		exists, err := targetAdder.MissionExists(missionID)
		if err != nil {
			logger.Error("failed to check if mission exists", slog.Any("error", err))
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		if !exists {
			logger.Error("mission does not exist", slog.Int64("missionID", missionID))
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("mission with ID %d does not exist", missionID))
			return
		}

		var req AddTargetRequest
		err = utils.ParseJSON(r, &req)
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

		targetID, err := targetAdder.AddTarget(missionID, req.Name, req.Country, req.Notes)
		if err != nil {
			logger.Error("failed to add target", slog.Any("error", err))
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		logger.Info("target added successfully", slog.Int64("targetID", targetID))
		utils.WriteJSON(w, http.StatusCreated, map[string]int64{"targetID": targetID})
	}
}
