package categoriesrepo

import (
	"context"
	"log/slog"
	"news-service/internal/entities"
	errMsg "news-service/internal/err"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoriesRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewCategoriesRepository(db *pgxpool.Pool, log *slog.Logger) *CategoriesRepository {
	return &CategoriesRepository{db: db, log: log}
}

func (c *CategoriesRepository) CreateCategorie(ctx context.Context, categorie *entities.Categorie) error {
	err := c.db.QueryRow(ctx, `INSERT INTO Categories (name) VALUES ($1) ON CONFLICT (name) DO UPDATE 
    SET name = EXCLUDED.name
	RETURNING id`, categorie.Name).Scan(&categorie.ID)
	if err != nil {
		c.log.Error("failed to create news", errMsg.Err(err))
		return err
	}
	return nil
}
