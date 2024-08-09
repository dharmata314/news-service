package newshandler

import (
	"log/slog"
	"net/http"
	"news-service/api/response"
	"news-service/internal/entities"
	errMsg "news-service/internal/err"
	"news-service/internal/models"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type RequestNews struct {
	Title      string `json:"Title"`
	Content    string `json:"Content"`
	Categories []int  `json:"Categories"`
}

type ResponseNews struct {
	response.Response
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Categories []int  `json:"categories"`
}

func NewNews(log *slog.Logger, NewsRepository models.NewsRepository, CategoriesRepository models.CategoriesRepository, NewsCategoriesRepository models.NewsCategoriesRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.createNews.New"

		log = log.With(
			slog.String("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestNews
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", errMsg.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}
		log.Info("request body request", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid requets", errMsg.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		news := entities.News{Title: req.Title, Content: req.Content}
		err = NewsRepository.CreateNews(r.Context(), &news)
		if err != nil {
			log.Error("failed to create news", errMsg.Err((err)))
			render.JSON(w, r, response.Error("failed to create  news"))
			return
		}
		for _, categorie_name := range req.Categories {
			categorie := entities.Categorie{Name: categorie_name}
			err = CategoriesRepository.CreateCategorie(r.Context(), &categorie)
			if err != nil {
				log.Error("could not add categorie to the table", errMsg.Err(err))
				return
			}

			newsCategories := entities.NewsCategories{CategoryID: categorie.ID, NewsID: news.ID}
			err = NewsCategoriesRepository.Create(r.Context(), &newsCategories)
			if err != nil {
				log.Error("could not add Newscategories relation to the table", errMsg.Err(err))
				return
			}
		}
		log.Info("news added to postgres")
		responseOK(w, r, news.ID, news.Title, news.Content, req.Categories)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id int, title, content string, categories []int) {
	render.JSON(w, r, ResponseNews{
		response.OK(),
		id,
		title,
		content,
		categories,
	})
}
