package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/real-splendid/url-shortener-practicum/internal"
	"github.com/real-splendid/url-shortener-practicum/internal/handlers"
	"github.com/real-splendid/url-shortener-practicum/internal/middleware"
)

// app представляет приложение.
type app struct {
	// router маршрутизатор приложения.
	router *chi.Mux
	// storage хранилище данных.
	storage internal.Storage
	// logger логгер приложения.
	logger *zap.SugaredLogger
	// baseURL базовый URL для коротких ссылок.
	baseURL string
}

// NewApp создает новый экземпляр приложения.
// Принимает хранилище данных, логгер, базовый URL и DSN базы данных.
// Возвращает указатель на экземпляр приложения.
func NewApp(
	storage internal.Storage,
	logger *zap.SugaredLogger,
	baseURL string,
	dDSN string,
) *app { // coverage-ignore
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
	router.Delete("/api/user/urls", handlers.MakeDeleteUserURLsHandler(storage, logger))
	return &app{
		router:  router,
		storage: storage,
		logger:  logger,
		baseURL: baseURL,
	}
}

// Serve запускает HTTP-сервер.
// Принимает указатель на строку, содержащую адрес для прослушивания.
// Возвращает ошибку, если произошла ошибка при запуске сервера.
func (a *app) Serve(address *string) error { // coverage-ignore
	defer a.logger.Sync()
	if err := http.ListenAndServe(*address, a.router); err != nil {
		a.logger.Error(err)
		return err
	}
	return nil
}
