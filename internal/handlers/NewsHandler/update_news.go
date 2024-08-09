package newshandler

import (
	"log/slog"
	"net/http"
	"news-service/api/response"
	"news-service/internal/entities"
	errMsg "news-service/internal/err"
	"news-service/internal/models"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type RequestUpdateNews struct {
	ID         int    `json:"Id" validate:"required"`
	Title      string `json:"Title"`
	Content    string `json:"Content"`
	Categories []int  `json:"Categories"`
}

func UpdateNews(log *slog.Logger, NewsRepository models.NewsRepository, CategoriesRepository models.CategoriesRepository, NewsCategoriesRepository models.NewsCategoriesRepository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		newsID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("failed to convert request parameter id", errMsg.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		const loggerOptions = "handlers.UpdateNews"

		log = log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestUpdateNews
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", errMsg.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}
		log.Info("request body request", slog.Any("request", req))

		news, err := NewsRepository.FindNewsByID(r.Context(), newsID)
		if err != nil {
			render.Status(r, http.StatusNotFound)
			log.Error("Failed to find news")
			render.JSON(w, r, response.Error("news not found"))
			return
		}
		news.ID = newsID
		news.Content = req.Content
		news.Title = req.Title
		err = NewsRepository.UpdateNews(r.Context(), &news)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			log.Error("Failed to update news")
			render.JSON(w, r, response.Error("Failed to update news"))
			return
		}
		NewsCategoriesRepository.DeleteCategories(r.Context(), news.ID)
		for _, categorie_name := range req.Categories {
			categorie := entities.Categorie{Name: categorie_name}
			err = CategoriesRepository.CreateCategorie(r.Context(), &categorie)
			if err != nil {
				log.Error("could not update categorie to the table", errMsg.Err(err))
				return
			}

			err = NewsCategoriesRepository.UpdateNewsCategories(r.Context(), categorie.ID, newsID)
			if err != nil {
				log.Error("could not update Newscategories relation to the table", errMsg.Err(err))
				return
			}
		}

		log.Info("news updated")

		render.JSON(w, r, response.OK())

	}

}
