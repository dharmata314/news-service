package database

import (
	"context"
	"fmt"
	"log/slog"
	"news-service/internal/config"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Db     *pgxpool.Pool
	log    *slog.Logger
	Config *config.Config
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

func NewPG(ctx context.Context, connString string, log *slog.Logger, cfg *config.Config) (*Postgres, error) {
	var err error

	pgOnce.Do(func() {
		var db *pgxpool.Pool
		db, err = pgxpool.New(ctx, connString)

		if err != nil {
			log.Error("unable to create connection pool", slog.String("error", err.Error()))
			err = fmt.Errorf("unable to create connection pool: %w", err)
			return
		}

		pgInstance = &Postgres{db, log, cfg}

		if err = CreateTable(ctx, db, log, cfg); err != nil {
			log.Error("failed to create tables", slog.String("error", err.Error()))
			return
		}
	})

	if err != nil {
		return nil, err
	}
	return pgInstance, nil
}

func CreateTable(ctx context.Context, db *pgxpool.Pool, log *slog.Logger, cfg *config.Config) error {

	_, err := db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS News (
	id SERIAL PRIMARY KEY,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Error("failed to create news table", slog.String("error", err.Error()))
		return fmt.Errorf("failed to create news table")
	}

	_, err = db.Exec(ctx, `CREATE TABLE IF NOT EXISTS Categories (
	id SERIAL PRIMARY KEY,
	name INT NOT NULL UNIQUE 
	)`)
	if err != nil {
		log.Error("failed to create categories table", slog.String("error", err.Error()))
		return fmt.Errorf("failed to create categories table")
	}

	_, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS NewsCategories (
    news_id INT REFERENCES News(id) ON DELETE CASCADE,
    category_id INT REFERENCES Categories(id) ON DELETE CASCADE,
    PRIMARY KEY (news_id, category_id)
	)`)
	if err != nil {
		log.Error("failed to create newsCategories table", slog.String("error", err.Error()))
		return fmt.Errorf("failed to create newsCategories table")
	}

	_, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS Users (
	    id SERIAL PRIMARY KEY, 
	    email VARCHAR(100) UNIQUE NOT NULL, 
	    password VARCHAR(255) NOT NULL,
	    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP)
	`)
	if err != nil {
		log.Error("failed to create users table", slog.String("error", err.Error()))
		return fmt.Errorf("failed to create users table: %w", err)
	}
	return nil

}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.Db.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.Db.Close()
}
