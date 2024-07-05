package spycat

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/illiakornyk/spy-cat/internal/breeds"
	"github.com/illiakornyk/spy-cat/internal/lib/api/response"
)

type CreateRequest struct {
    Name              string  `json:"name" validate:"required,min=1,max=100"`
    YearsOfExperience int     `json:"years_of_experience" validate:"min=0"`
    Breed             string  `json:"breed" validate:"required,min=1,max=100"`
    Salary            float64 `json:"salary" validate:"required,gt=0"`
}

type CreateResponse struct {
	response.Response
	ID     int64  `json:"id,omitempty"`
}


type SpyCatCreator interface {
    CreateCat(name string, yearsOfExperience int, breed string, salary float64) (int64, error)
}

func New(logger *slog.Logger, spyCatCreator SpyCatCreator) http.HandlerFunc {
	validate := validator.New()

    return func(w http.ResponseWriter, r *http.Request) {
        const op = "handlers.spycat.create"

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

		// Validate the breed
        if !breeds.IsValidBreed(req.Breed) {
            logger.Error("invalid breed", slog.String("breed", req.Breed))

            json.NewEncoder(w).Encode(response.Response{
                Status: response.StatusError,
                Error:  "invalid breed",
            })
            return
        }


        id, err := spyCatCreator.CreateCat(req.Name, req.YearsOfExperience, req.Breed, req.Salary)
        if err != nil {
            logger.Error("failed to create spy cat", slog.Any("error", err))

            json.NewEncoder(w).Encode(response.Response{
                Status: response.StatusError,
                Error:  "failed to create spy cat",
            })
            return
        }

        logger.Info("spy cat created successfully", slog.Int64("id", id))

        json.NewEncoder(w).Encode(CreateResponse{
            Response: response.Response{
                Status: response.StatusOK,
            },
            ID: id,
        })
    }
	}
