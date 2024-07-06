package spycat

import (
	"net/http"

	"github.com/illiakornyk/spy-cat/internal/utils"

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
    Cats []SpyCat `json:"cats"`
}

func GetAllHandler(logger *slog.Logger, spyCatGetter SpyCatsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.spycat.get_all"

		logger = logger.With(slog.String("op", op))

		cats, err := spyCatGetter.GetAllCats()
		if err != nil {
			logger.Error("failed to get all spy cats", slog.Any("error", err))
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		logger.Info("retrieved all spy cats successfully", slog.Int("count", len(cats)))

		utils.WriteJSON(w, http.StatusOK, GetAllResponse{
			Cats: cats,
		})
	}
}
