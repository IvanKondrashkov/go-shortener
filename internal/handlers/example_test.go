package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
	"github.com/IvanKondrashkov/go-shortener/internal/service"
	"github.com/IvanKondrashkov/go-shortener/internal/service/worker"
	"github.com/IvanKondrashkov/go-shortener/internal/storage/mem"

	"github.com/go-chi/chi/v5"
)

func setUpApp() *App {
	zl, _ := logger.NewZapLogger(config.LogLevel)

	var newRepository service.Repository
	var newRunner service.Runner

	newRepository = mem.NewRepository(zl)
	newRunner = newRepository
	// В реальном коде используйте NewSuite для инициализации
	newService := service.NewService(zl, newRunner, newRepository)
	newWorker := worker.NewWorker(context.Background(), config.WorkerCount, zl, newService)
	return NewApp(newService, newWorker)
}

// Пример использования ShortenURL (текстовый формат)
func ExampleApp_ShortenURL() {
	app := setUpApp()

	// Создаем тестовый запрос
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com/long/url"))
	w := httptest.NewRecorder()

	// Вызываем хендлер
	app.ShortenURL(w, req)

	// Выводим результат
	fmt.Println("Status:", w.Code)
	fmt.Println("Short URL:", w.Body.String())
	// Output:
	// Status: 201
	// Short URL: http://localhost:8080/6b88474e-b4d5-552e-8d1c-eb63852e8b75
}

// Пример использования ShortenAPI (JSON формат)
func ExampleApp_ShortenAPI() {
	app := setUpApp()

	request := models.RequestShortenAPI{
		URL: "https://example.com/json-api",
	}
	body, _ := json.Marshal(request)

	// Создаем тестовый запрос
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Вызываем хендлер
	app.ShortenAPI(w, req)

	var response models.ResponseShortenAPI
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	// Выводим результат
	fmt.Println("Status:", w.Code)
	fmt.Println("Result:", response.Result)
	// Output:
	// Status: 201
	// Result: http://localhost:8080/7ac2c9ff-7699-5787-9e11-8af9321965d7
}

// Пример использования GetURLByID (текстовый формат)
func ExampleApp_GetURLByID() {
	app := setUpApp()

	// Предварительно создаем URL
	createReq := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
	createResp := httptest.NewRecorder()
	app.ShortenURL(createResp, createReq)
	shortURL := createResp.Body.String()
	id := strings.TrimPrefix(shortURL, app.URL)

	// Создаем тестовый запрос
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	req := httptest.NewRequest(http.MethodGet, "/"+id, nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	// Вызываем хендлер
	app.GetURLByID(w, req)

	// Вывод результатов
	fmt.Println("Status:", w.Code)
	fmt.Println("Location:", w.Header().Get("Location"))
	// Output:
	// Status: 307
	// Location: https://example.com
}
