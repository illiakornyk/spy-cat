package spycat

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/illiakornyk/spy-cat/internal/lib/api/response"
)

type Request struct {
	Name              string  `json:"name" validate:"required"`
	YearsOfExperience int     `json:"years_of_experience" validate:"required"`
	Breed             string  `json:"breed" validate:"required"`
	Salary            float64 `json:"salary" validate:"required"`
}

type Response struct {
	response.Response
	ID     int64  `json:"id,omitempty"`
}


type SpyCatSaver interface {
    SaveCat(name string, yearsOfExperience int, breed string, salary float64) (int64, error)
}


func New(logger *slog.Logger, spyCatSaver SpyCatSaver) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const op = "handlers.spycat.save.New"

        logger = logger.With(slog.String("op", op))

        var req Request

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

        // Log the decoded request body
        logger.Info("request body decoded", slog.Any("req", req))

        id, err := spyCatSaver.SaveCat(req.Name, req.YearsOfExperience, req.Breed, req.Salary)
        if err != nil {
            logger.Error("failed to save spy cat", slog.Any("error", err))

            json.NewEncoder(w).Encode(response.Response{
                Status: response.StatusError,
                Error:  "failed to save spy cat",
            })
            return
        }

        logger.Info("spy cat saved successfully", slog.Int64("id", id))

        json.NewEncoder(w).Encode(Response{
            Response: response.Response{
                Status: response.StatusOK,
            },
            ID: id,
        })
    }
}
