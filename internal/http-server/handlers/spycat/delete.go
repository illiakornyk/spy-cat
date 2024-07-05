package spycat

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/illiakornyk/spy-cat/internal/lib/api/response"

	"log/slog"

	"github.com/go-chi/chi/v5"
)

type SpyCatDeleter interface {
    DeleteCat(id int64) error
	CatExists(id int64) (bool, error)
}

type DeleteResponse struct {
    response.Response
}

func DeleteHandler(logger *slog.Logger, spyCatDeleter SpyCatDeleter) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const op = "handlers.spycat.delete"

        logger = logger.With(slog.String("op", op))

        idStr := chi.URLParam(r, "id")
		        logger.Info("Extracted ID from URL", slog.String("idStr", idStr)) // Add log

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

        exists, err := spyCatDeleter.CatExists(id)
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

        err = spyCatDeleter.DeleteCat(id)
        if err != nil {
            logger.Error("failed to delete spy cat", slog.Any("error", err))

            json.NewEncoder(w).Encode(response.Response{
                Status: response.StatusError,
                Error:  "failed to delete spy cat",
            })
            return
        }

        logger.Info("spy cat deleted successfully", slog.Int64("id", id))

		w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(DeleteResponse{
            Response: response.Response{
                Status: response.StatusOK,
            },
        })
    }
}
