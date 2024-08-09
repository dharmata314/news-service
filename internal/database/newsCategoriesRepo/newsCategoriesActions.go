package newscategoriesrepo

import (
	"context"
	"log/slog"
	"news-service/internal/entities"
	errMsg "news-service/internal/err"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsCategoriesRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewNewsCategoriesRepository(db *pgxpool.Pool, log *slog.Logger) *NewsCategoriesRepository {
	return &NewsCategoriesRepository{db, log}
}

func (n *NewsCategoriesRepository) Create(ctx context.Context, NC *entities.NewsCategories) error {
	_, err := n.db.Exec(ctx, `INSERT INTO NewsCategories (news_id, category_id) VALUES ($1, $2)`, NC.NewsID, NC.CategoryID)
	if err != nil {
		n.log.Error("failed to create newsCategorie", errMsg.Err(err))
		return err
	}
	return nil
}

func (n *NewsCategoriesRepository) ListCategories(ctx context.Context, id int) ([]int, error) {
	var arrayId []int
	err := n.db.QueryRow(ctx, `
 	SELECT array_agg(c.name)
	FROM Categories c
	WHERE c.id = ANY(
    	SELECT nc.category_id
   	 FROM NewsCategories nc
    	WHERE nc.news_id = $1)`, id).Scan(&arrayId)
	if err != nil {
		n.log.Error("failed to list categories", errMsg.Err(err))
		return nil, err
	}
	return arrayId, nil
}

func (n *NewsCategoriesRepository) UpdateNewsCategories(ctx context.Context, categoryID, newsID int) error {

	_, err := n.db.Exec(ctx, `INSERT INTO NewsCategories (news_id, category_id) VALUES ($1, $2)`, newsID, categoryID)
	if err != nil {
		n.log.Error("failed to insert new news category", errMsg.Err(err))
		return err
	}

	return nil

}

func (n *NewsCategoriesRepository) DeleteCategories(ctx context.Context, newsID int) error {
	_, err := n.db.Exec(ctx, `DELETE FROM NewsCategories WHERE news_id = $1`, newsID)
	if err != nil {
		n.log.Error("failed to delete existing news categories", errMsg.Err(err))
		return err
	}
	return nil
}
