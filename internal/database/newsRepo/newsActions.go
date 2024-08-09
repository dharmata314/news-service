package newsrepo

import (
	"context"
	"fmt"
	"log/slog"
	"news-service/internal/entities"
	errMsg "news-service/internal/err"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewNewsRepository(db *pgxpool.Pool, log *slog.Logger) *NewsRepository {
	return &NewsRepository{db: db, log: log}
}

func (n *NewsRepository) CreateNews(ctx context.Context, news *entities.News) error {
	err := n.db.QueryRow(ctx, `INSERT INTO News (content, title) VALUES ($1, $2) RETURNING id`, news.Content, news.Title).Scan(&news.ID)
	if err != nil {
		n.log.Error("failed to create news", errMsg.Err(err))
		return err
	}
	return nil
}

func (n *NewsRepository) ListNews(ctx context.Context) ([]entities.News, error) {
	query, err := n.db.Query(ctx, `SELECT id, title, content FROM News ORDER BY id DESC`)
	if err != nil {
		n.log.Error("Error querying news", errMsg.Err(err))
		return nil, err
	}
	defer query.Close()

	var newsArray []entities.News
	for query.Next() {
		var news entities.News
		err := query.Scan(&news.ID, &news.Title, &news.Content)
		if err != nil {
			n.log.Error("Error scanning news", errMsg.Err(err))
			return nil, err
		}
		newsArray = append(newsArray, news)
	}

	if err := query.Err(); err != nil {
		n.log.Error("Error iterating over news", errMsg.Err(err))
		return nil, err
	}

	return newsArray, nil

}

func (n *NewsRepository) UpdateNews(ctx context.Context, news *entities.News) error {
	_, err := n.db.Exec(ctx, `UPDATE News SET content = $1, title = $2 WHERE id = $3`, news.Content, news.Title, news.ID)
	if err != nil {
		n.log.Error("failed to update news", errMsg.Err(err))
		return err
	}

	return nil

}

func (n *NewsRepository) FindNewsByID(ctx context.Context, id int) (entities.News, error) {
	query, err := n.db.Query(ctx, `SELECT content, title FROM News WHERE id = $1`, id)
	if err != nil {
		n.log.Error("error querying news", errMsg.Err(err))
		return entities.News{}, err
	}
	defer query.Close()
	row := entities.News{}
	if !query.Next() {
		n.log.Error("news not found")
		return entities.News{}, fmt.Errorf("news not found")
	} else {
		err := query.Scan(&row.Content, &row.Title)
		if err != nil {
			n.log.Error("error scanning news", errMsg.Err(err))
			return entities.News{}, err
		}
	}
	return row, nil
}
