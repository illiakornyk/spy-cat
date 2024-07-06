package spycat

import (
	"fmt"
	"net/http"
	"strconv"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/illiakornyk/spy-cat/internal/utils"
)

type SpyCatUpdater interface {
    UpdateCatSalary(id int64, salary float64) error
	CatExists(id int64) (bool, error)
}

type PatchRequest struct {
    Salary float64 `json:"salary" validate:"required,gt=0"`
}


func PatchHandler(logger *slog.Logger, spyCatUpdater SpyCatUpdater) http.HandlerFunc {
	validate := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.spycat.patch"
		logger = logger.With(slog.String("op", op))

		idStr := chi.URLParam(r, "id")
		logger.Info("Extracted ID from URL", slog.String("idStr", idStr))

		if idStr == "" {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id path parameter is missing"))
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id < 1 {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id path parameter"))
			return
		}

		exists, err := spyCatUpdater.CatExists(id)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to check if cat exists"))
			return
		}

		if !exists {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("cat not found"))
			return
		}

		var req PatchRequest
		err = utils.ParseJSON(r, &req)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to decode request body"))
			return
		}

		err = validate.Struct(req)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("validation failed: %w", err))
			return
		}

		err = spyCatUpdater.UpdateCatSalary(id, req.Salary)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to update spy cat salary"))
			return
		}

		logger.Info("spy cat salary updated successfully", slog.Int64("id", id), slog.Float64("salary", req.Salary))
		w.WriteHeader(http.StatusNoContent)
	}
}
