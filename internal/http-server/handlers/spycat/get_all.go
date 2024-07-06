package spycat

import (
	"encoding/json"
	"net/http"

	"github.com/illiakornyk/spy-cat/internal/lib/api/response"

	"log/slog"
)

type SpyCatsGetter interface {
    GetAllCats() ([]SpyCat, error)
}

type SpyCat struct {
    ID               int64   `json:"id"`
    Name             string  `json:"name"`
    YearsOfExperience int     `json:"years_of_experience"`
    Breed            string  `json:"breed"`
    Salary           float64 `json:"salary"`
}

type GetAllResponse struct {
    response.Response
    Cats []SpyCat `json:"cats"`
}

func GetAllHandler(logger *slog.Logger, spyCatGetter SpyCatsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.spycat.get_all"

		logger = logger.With(slog.String("op", op))

		cats, err := spyCatGetter.GetAllCats()
		if err != nil {
			logger.Error("failed to get all spy cats", slog.Any("error", err))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Response{
				Status: response.StatusError,
				Error:  "failed to get all spy cats",
			})
			return
		}

		logger.Info("retrieved all spy cats successfully", slog.Int("count", len(cats)))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GetAllResponse{
			Response: response.Response{
				Status: response.StatusOK,
			},
			Cats: cats,
		})
	}
}
