package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/real-splendid/url-shortener-practicum/internal"
	"github.com/real-splendid/url-shortener-practicum/internal/handlers"
	"github.com/real-splendid/url-shortener-practicum/internal/middleware"
	"go.uber.org/zap"
)

type app struct {
	router  *chi.Mux
	storage internal.Storage
	logger  *zap.SugaredLogger
	baseURL string
}

func NewApp(storage internal.Storage, baseURL string, dDSN string) *app {
	rawLogger, _ := zap.NewDevelopment()
	logger := rawLogger.Sugar()
	router := chi.NewRouter()
	router.Use(middleware.MakeLogMiddleware(logger))
	router.Use(middleware.MakeGzipMiddleware(logger))
	router.Post("/", handlers.MakeShortenHandler(storage, logger, baseURL))
	router.Post("/api/shorten", handlers.MakeAPIShortenHandler(storage, logger, baseURL))
	router.Get("/{key}", handlers.MakeRedirectionHandler(storage, logger))
	router.Get("/ping", handlers.MakePingHandler(dDSN, logger))
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
