package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"news-service/internal/config"
	"news-service/internal/database"
	categoriesrepo "news-service/internal/database/categoriesRepo"
	newscategoriesrepo "news-service/internal/database/newsCategoriesRepo"
	newsrepo "news-service/internal/database/newsRepo"
	usersrepo "news-service/internal/database/usersRepo"
	newshandler "news-service/internal/handlers/NewsHandler"
	userhandlers "news-service/internal/handlers/userHandler"
	"news-service/internal/jwt"
	"os"

	errMsg "news-service/internal/err"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger()
	log.Debug("debug messages are active")

	pg, err := connectToPostgres(cfg, log)
	if err != nil {
		log.Error("failed to create postgres db", errMsg.Err(err))
		os.Exit(1)
	}

	log.Info("connecting to postgres")

	defer pg.Close()
	if pg == nil {
		log.Error("failed to connect to postgres")
		os.Exit(1)
	}

	if err := pg.Ping(context.Background()); err != nil {
		log.Error("failed to ping postgres db", errMsg.Err(err))
		os.Exit(1)
	}

	log.Info("postgres db connected successfully")

	log.Info("application started")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	newsRepository := newsrepo.NewNewsRepository(pg.Db, log)
	categoriesRepository := categoriesrepo.NewCategoriesRepository(pg.Db, log)
	newsCategoriesRepository := newscategoriesrepo.NewNewsCategoriesRepository(pg.Db, log)
	userRepository := usersrepo.NewUserRepository(pg.Db, log)

	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, log)

	router.Post("/users/new", userhandlers.NewUser(log, userRepository))
	router.Post("/login", userhandlers.LoginFunc(log, userRepository, jwtManager))

	router.With(func(next http.Handler) http.Handler {
		return jwt.TokenAuthMiddleware(jwtManager, next)
	}).Post("/news", newshandler.NewNews(log, newsRepository, categoriesRepository, newsCategoriesRepository))

	router.With(func(next http.Handler) http.Handler {
		return jwt.TokenAuthMiddleware(jwtManager, next)
	}).Get("/list", newshandler.ListAllNews(log, newsRepository, newsCategoriesRepository))

	router.With(func(next http.Handler) http.Handler {
		return jwt.TokenAuthMiddleware(jwtManager, next)
	}).Patch("/news/edit/{id}", newshandler.UpdateNews(log, newsRepository, categoriesRepository, newsCategoriesRepository))

	server := &http.Server{
		Addr:              cfg.HTTPServer.Addr,
		Handler:           router,
		ReadHeaderTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout:      cfg.HTTPServer.Timeout,
		IdleTimeout:       cfg.HTTPServer.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server", errMsg.Err(err))
	}

}

func setupLogger() *slog.Logger {
	var log *slog.Logger = slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}

func connectToPostgres(cfg *config.Config, log *slog.Logger) (*database.Postgres, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)
	pg, err := database.NewPG(context.Background(), connString, log, cfg)
	return pg, err
}
