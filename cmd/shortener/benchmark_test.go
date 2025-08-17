package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/IvanKondrashkov/go-shortener/internal/handlers"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
	"github.com/IvanKondrashkov/go-shortener/internal/service"
	customContext "github.com/IvanKondrashkov/go-shortener/internal/service/middleware/auth"
	"github.com/IvanKondrashkov/go-shortener/internal/storage/mem"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func BenchmarkShortenURL(b *testing.B) {
	repo := mem.NewRepository(nil)
	svc := service.NewService(nil, repo, repo)
	app := handlers.NewApp(svc, nil)

	u := "https://example.com/very/long/url/to/be/shortened"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(u))
		res := httptest.NewRecorder()
		app.ShortenURL(res, req)
	}
}

func BenchmarkShortenAPI(b *testing.B) {
	repo := mem.NewRepository(nil)
	svc := service.NewService(nil, repo, repo)
	app := handlers.NewApp(svc, nil)

	reqDto := models.RequestShortenAPI{URL: "https://example.com/very/long/url/to/be/shortened"}
	body, _ := json.Marshal(reqDto)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(body))
		res := httptest.NewRecorder()
		app.ShortenAPI(res, req)
	}
}

func BenchmarkGetURLByID(b *testing.B) {
	repo := mem.NewRepository(nil)
	svc := service.NewService(nil, repo, repo)
	app := handlers.NewApp(svc, nil)

	u, _ := url.Parse("https://example.com")
	id, _ := svc.Save(context.Background(), uuid.New(), u)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id.String())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/"+id.String(), nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()
		app.GetURLByID(res, req)
	}
}

func BenchmarkGetAllURLByUserID(b *testing.B) {
	repo := mem.NewRepository(nil)
	svc := service.NewService(nil, repo, repo)
	app := handlers.NewApp(svc, nil)

	ctx := customContext.SetContextUserID(context.Background(), uuid.New())
	for i := 0; i < 10; i++ {
		u, _ := url.Parse("https://example.com/" + strconv.Itoa(i))
		svc.Save(ctx, uuid.New(), u)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
		req = req.WithContext(ctx)
		res := httptest.NewRecorder()
		app.GetAllURLByUserID(res, req)
	}
}
