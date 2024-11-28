package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/IvanKondrashkov/go-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL(t *testing.T) {
	baseURL := "http://localhost:8080/"
	app := &App{
		BaseURL:    baseURL,
		repository: storage.NewMemRepositoryImpl(),
	}

	tests := []struct {
		name   string
		url    string
		status int
		want   []byte
	}{
		{
			name:   "is invalidate url",
			status: http.StatusBadRequest,
			url:    "://ya.ru/",
			want:   []byte("Url is invalidate!"),
		},
		{
			name:   "ok",
			status: http.StatusCreated,
			url:    "https://ya.ru/",
			want:   []byte(baseURL + uuid.NewSHA1(uuid.NameSpaceURL, []byte("https://ya.ru/")).String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBuffer([]byte(tt.url))
			req := httptest.NewRequest(http.MethodPost, baseURL, b)

			w := httptest.NewRecorder()

			app.ShortenURL(w, req)

			assert.Equal(t, tt.status, w.Code)
			assert.Equal(t, tt.want, w.Body.Bytes())
		})
	}
}

func TestGetURLByID(t *testing.T) {
	baseURL := "http://localhost:8080/"
	app := &App{
		BaseURL:    baseURL,
		repository: storage.NewMemRepositoryImpl(),
	}

	tests := []struct {
		name   string
		id     uuid.UUID
		status int
		want   string
	}{
		{
			name:   "id not found",
			status: http.StatusNotFound,
			id:     uuid.New(),
			want:   "Url by id not found!",
		},
		{
			name:   "ok",
			status: http.StatusTemporaryRedirect,
			id:     uuid.NewSHA1(uuid.NameSpaceURL, []byte("https://ya.ru/")),
			want:   "https://ya.ru/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, baseURL+tt.id.String(), nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id.String())
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			if tt.status == http.StatusTemporaryRedirect {
				u, _ := url.Parse(tt.want)
				app.repository.Save(tt.id, u)
				app.GetURLByID(w, req)

				assert.Equal(t, tt.status, w.Code)
				assert.Equal(t, tt.want, w.Header().Get("Location"))
			} else {
				app.GetURLByID(w, req)

				assert.Equal(t, tt.status, w.Code)
			}

		})
	}
}
