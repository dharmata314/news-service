package newshandler

import (
	"log/slog"
	"net/http"
	"news-service/api/response"
	errMsg "news-service/internal/err"
	"news-service/internal/models"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type NewsItem struct {
	ID         int    `json:"Id"`
	Title      string `json:"Title"`
	Content    string `json:"Content"`
	Categories []int  `json:"Categories"`
}

type ResponseNewsList struct {
	Success bool       `json:"Success"`
	News    []NewsItem `json:"News"`
}

func ListAllNews(log *slog.Logger, newsRepository models.NewsRepository, newsCategoriesRepository models.NewsCategoriesRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const loggerOptions = "handlers.listAllNews"
		log = log.With(
			slog.Any("options", loggerOptions),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		newsArray, err := newsRepository.ListNews(r.Context())
		if err != nil {
			log.Error("Failed to retrieve news", errMsg.Err(err))
			render.JSON(w, r, response.Error("Failed to retrieve news"))
			return
		}
		result := make([]NewsItem, len(newsArray))
		for i, news := range newsArray {
			var categories []int

			categories, err := newsCategoriesRepository.ListCategories(r.Context(), news.ID)
			if err != nil {
				log.Error("Failed to retrieve news", errMsg.Err(err))
				render.JSON(w, r, response.Error("Failed to retrieve news"))
				return
			}

			result[i] = NewsItem{
				ID:         news.ID,
				Title:      news.Title,
				Content:    news.Content,
				Categories: categories,
			}
		}

		responseOKgetNews(w, r, result)
	}
}

func responseOKgetNews(w http.ResponseWriter, r *http.Request, news []NewsItem) {
	render.JSON(w, r, ResponseNewsList{
		News:    news,
		Success: true,
	})
}
