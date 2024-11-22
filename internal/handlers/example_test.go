package handlers_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/real-splendid/url-shortener-practicum/internal/handlers"
	"github.com/real-splendid/url-shortener-practicum/internal/middleware"
	"github.com/real-splendid/url-shortener-practicum/internal/storage"
	"go.uber.org/zap"
)

func ExampleMakeShortenHandler() {
	// Создаем тестовое хранилище
	testStorage := storage.NewMemoryStorage()

	// Создаем логгер
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	// Создаем хэндлер
	handler := handlers.MakeShortenHandler(testStorage, sugar, "http://localhost:8080")

	// Создаем тестовый запрос
	body := []byte("https://example.com")
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	mockUserID := "test-user-id"
	mockSignedCookie := middleware.SignCookie(mockUserID)
	req.AddCookie(&http.Cookie{Name: "user_id", Value: mockSignedCookie})

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Вызываем хэндлер
	handler.ServeHTTP(rr, req)

	fmt.Printf("Status: %d\n", rr.Code)

	// Output:
	// Status: 201
}

func ExampleMakeRedirectionHandler() {
	// Создаем тестовое хранилище
	testStorage := storage.NewMemoryStorage()
	originalURL := "https://ya.ru"
	key := "testtest"
	userID := "test-user-id"
	testStorage.Set(key, originalURL, userID)

	// Создаем логгер
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	// Создаем хэндлер
	handler := handlers.MakeRedirectionHandler(testStorage, sugar)

	// Создаем тестовый запрос
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", key), nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("key", key)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Вызываем хэндлер
	handler.ServeHTTP(rr, req)

	fmt.Printf("Status: %d\n", rr.Code)
	fmt.Printf("Location: %s\n", rr.Header().Get("Location"))

	// Output:
	// Status: 307
	// Location: https://ya.ru
}
