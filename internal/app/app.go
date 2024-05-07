package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/real-splendid/url-shortener-practicum/internal"
	"github.com/real-splendid/url-shortener-practicum/internal/handlers"
	"github.com/real-splendid/url-shortener-practicum/internal/middleware"
)

type app struct {
	router  *chi.Mux
	storage internal.Storage
	logger  *zap.SugaredLogger
	baseURL string
}

func NewApp(storage internal.Storage, logger *zap.SugaredLogger, baseURL string, dDSN string) *app {
	router := chi.NewRouter()
	router.Use(middleware.MakeLogMiddleware(logger))
	router.Use(middleware.MakeGzipMiddleware(logger))
	router.Use(middleware.MakeAuthMiddleware(logger))
	router.Post("/", handlers.MakeShortenHandler(storage, logger, baseURL))
	router.Post("/api/shorten", handlers.MakeAPIShortenHandler(storage, logger, baseURL))
	router.Post("/api/shorten/batch", handlers.MakeAPIShortenBatchHandler(storage, logger, baseURL))
	router.Get("/{key}", handlers.MakeRedirectionHandler(storage, logger))
	router.Get("/ping", handlers.MakePingHandler(dDSN, logger))
	router.Get("/api/user/urls", handlers.MakeUserURLsHandler(storage, logger, baseURL))
	return &app{
		router:  router,
		storage: storage,
		logger:  logger,
		baseURL: baseURL,
	}
}

func (a *app) Serve(address *string) error {
	defer a.logger.Sync()
	if err := http.ListenAndServe(*address, a.router); err != nil {
		a.logger.Error(err)
		return err
	}
	return nil
}
