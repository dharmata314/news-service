package models

import (
	"context"
	"news-service/internal/entities"
)

type NewsRepository interface {
	CreateNews(ctx context.Context, news *entities.News) error
	ListNews(ctx context.Context) ([]entities.News, error)
	UpdateNews(ctx context.Context, news *entities.News) error
	FindNewsByID(ctx context.Context, id int) (entities.News, error)
}

type CategoriesRepository interface {
	CreateCategorie(ctx context.Context, categorie *entities.Categorie) error
}

type NewsCategoriesRepository interface {
	Create(ctx context.Context, NC *entities.NewsCategories) error
	ListCategories(ctx context.Context, id int) ([]int, error)
	UpdateNewsCategories(ctx context.Context, categoryID, newsID int) error
	DeleteCategories(ctx context.Context, newsID int) error
}
