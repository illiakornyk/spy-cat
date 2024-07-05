package spycat

import (
	"encoding/json"
	"net/http"
	"strconv"

	"log/slog"

	"github.com/illiakornyk/spy-cat/internal/lib/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type SpyCatUpdater interface {
    UpdateCatSalary(id int64, salary float64) error
	CatExists(id int64) (bool, error)
}

type PatchRequest struct {
    Salary float64 `json:"salary" validate:"required,gt=0"`
}

type PatchResponse struct {
    response.Response
}

func PatchHandler(logger *slog.Logger, spyCatUpdater SpyCatUpdater) http.HandlerFunc {
    validate := validator.New()

    return func(w http.ResponseWriter, r *http.Request) {
        const op = "handlers.spycat.patch"

        logger = logger.With(slog.String("op", op))

        idStr := chi.URLParam(r, "id")
        logger.Info("Extracted ID from URL", slog.String("idStr", idStr))

        if idStr == "" {
            logger.Error("id path parameter is missing")

            json.NewEncoder(w).Encode(response.Response{
                Status: response.StatusError,
                Error:  "id path parameter is missing",
            })
            return
        }

        id, err := strconv.ParseInt(idStr, 10, 64)
        if err != nil || id < 1 {
            logger.Error("invalid id path parameter", slog.Any("error", err))

            json.NewEncoder(w).Encode(response.Response{
                Status: response.StatusError,
                Error:  "invalid id path parameter",
            })
            return
        }

        exists, err := spyCatUpdater.CatExists(id)
        if err != nil {
            logger.Error("failed to check if cat exists", slog.Any("error", err))

            json.NewEncoder(w).Encode(response.Response{
                Status: response.StatusError,
                Error:  "failed to check if cat exists",
            })
            return
        }

        if !exists {
            logger.Error("cat not found", slog.Int64("id", id))

            json.NewEncoder(w).Encode(response.Response{
                Status: response.StatusError,
                Error:  "cat not found",
            })
            return
        }


        var req PatchRequest
        err = json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
            logger.Error("failed to decode request body", slog.Any("error", err))

            json.NewEncoder(w).Encode(response.Response{
                Status: response.StatusError,
                Error:  "failed to decode request body",
            })
            return
        }

        err = validate.Struct(req)
        if err != nil {
            logger.Error("validation failed", slog.Any("error", err))

            json.NewEncoder(w).Encode(response.Response{
                Status: response.StatusError,
                Error:  "validation failed",
            })
            return
        }

        err = spyCatUpdater.UpdateCatSalary(id, req.Salary)
        if err != nil {
            logger.Error("failed to update spy cat salary", slog.Any("error", err))

            json.NewEncoder(w).Encode(response.Response{
                Status: response.StatusError,
                Error:  "failed to update spy cat salary",
            })
            return
        }

        logger.Info("spy cat salary updated successfully", slog.Int64("id", id), slog.Float64("salary", req.Salary))

		w.Header().Set("Content-Type", "application/json")

        json.NewEncoder(w).Encode(PatchResponse{
            Response: response.Response{
                Status: response.StatusOK,
            },
        })
    }
}
