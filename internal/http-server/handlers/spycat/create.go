package spycat

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/illiakornyk/spy-cat/internal/breeds"
	"github.com/illiakornyk/spy-cat/internal/utils"
)

type CreateRequest struct {
    Name              string  `json:"name" validate:"required,min=1,max=100"`
    YearsOfExperience int     `json:"years_of_experience" validate:"min=0"`
    Breed             string  `json:"breed" validate:"required,min=1,max=100"`
    Salary            float64 `json:"salary" validate:"required,gt=0"`
}

type CreateResponse struct {
	ID     int64  `json:"id,omitempty"`
}


type SpyCatCreator interface {
    CreateCat(name string, yearsOfExperience int, breed string, salary float64) (int64, error)
}

func CreateHandler(logger *slog.Logger, spyCatCreator SpyCatCreator) http.HandlerFunc {
	validate := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.spycat.create"

		logger = logger.With(slog.String("op", op))

		var req CreateRequest

		err := utils.ParseJSON(r, &req)
		if errors.Is(err, io.EOF) {
			logger.Error("request body is empty")
			utils.WriteError(w, http.StatusBadRequest, errors.New("empty request"))
			return
		}
		if err != nil {
			logger.Error("failed to decode request body", slog.Any("error", err))
			utils.WriteError(w, http.StatusBadRequest, errors.New("failed to decode request"))
			return
		}

		logger.Info("request body decoded", slog.Any("req", req))

		err = validate.Struct(req)
		if err != nil {
			logger.Error("validation failed", slog.Any("error", err))
			utils.WriteError(w, http.StatusBadRequest, errors.New("validation failed"))
			return
		}

		// Validate the breed
		if !breeds.IsValidBreed(req.Breed) {
			logger.Error("invalid breed", slog.String("breed", req.Breed))
			utils.WriteError(w, http.StatusBadRequest, errors.New("invalid breed"))
			return
		}

		id, err := spyCatCreator.CreateCat(req.Name, req.YearsOfExperience, req.Breed, req.Salary)
		if err != nil {
			logger.Error("failed to create spy cat", slog.Any("error", err))
			utils.WriteError(w, http.StatusInternalServerError, errors.New("failed to create spy cat"))
			return
		}

		logger.Info("spy cat created successfully", slog.Int64("id", id))
		utils.WriteJSON(w, http.StatusCreated, CreateResponse{
			ID: id,
		})
	}
}
