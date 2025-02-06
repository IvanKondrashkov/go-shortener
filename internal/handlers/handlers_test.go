package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	customContext "github.com/IvanKondrashkov/go-shortener/internal/context"
	"github.com/IvanKondrashkov/go-shortener/internal/handlers/mock"
	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/service"
	"github.com/IvanKondrashkov/go-shortener/internal/storage"
	"github.com/IvanKondrashkov/go-shortener/internal/worker"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type Suite struct {
	*testing.T
	app *App
}

func New(t *testing.T) *Suite {
	t.Helper()
	t.Parallel()

	zl, _ := logger.NewZapLogger(config.LogLevel)
	memRepositoryImpl := storage.NewMemRepositoryImpl(zl)
	newService := service.NewService(zl, memRepositoryImpl)
	newWorker := worker.NewWorker(config.WorkerCount, zl, newService)
	app := &App{
		URL:     config.URL,
		service: newService,
		worker:  newWorker,
	}

	return &Suite{
		T:   t,
		app: app,
	}
}

func TestShortenURL(t *testing.T) {
	tc := New(t)
	tests := []struct {
		name    string
		payload string
		status  int
		want    []byte
	}{
		{
			name:    "is invalidate url",
			payload: "://ya.ru/",
			status:  http.StatusBadRequest,
			want:    []byte("Url is invalidate!"),
		},
		{
			name:    "ok",
			payload: "https://ya.ru/",
			status:  http.StatusCreated,
			want:    []byte(tc.app.URL + uuid.NewSHA1(uuid.NameSpaceURL, []byte("https://ya.ru/")).String()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBuffer([]byte(tt.payload))
			req := httptest.NewRequest(http.MethodPost, tc.app.URL, b)

			w := httptest.NewRecorder()
			tc.app.ShortenURL(w, req)

			assert.Equal(t, tt.status, w.Code)
			assert.Equal(t, tt.want, w.Body.Bytes())
		})
	}
}

func TestShortenAPI(t *testing.T) {
	tc := New(t)
	tests := []struct {
		name    string
		payload []byte
		status  int
		want    []byte
	}{
		{
			name:    "is invalidate url",
			payload: []byte("{\"url\":\"://ya.ru/\"}"),
			status:  http.StatusBadRequest,
			want:    []byte("Url is invalidate!"),
		},
		{
			name:    "ok",
			payload: []byte("{\"url\":\"https://ya.ru/\"}"),
			status:  http.StatusCreated,
			want:    []byte("{\"result\":\"" + (tc.app.URL + uuid.NewSHA1(uuid.NameSpaceURL, []byte("https://ya.ru/")).String()) + "\"}\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBuffer(tt.payload)
			req := httptest.NewRequest(http.MethodPost, tc.app.URL+"api/shorten", b)

			w := httptest.NewRecorder()
			tc.app.ShortenAPI(w, req)

			assert.Equal(t, tt.status, w.Code)
			assert.Equal(t, tt.want, w.Body.Bytes())
		})
	}
}

func TestShortenAPIBatch(t *testing.T) {
	tc := New(t)
	tests := []struct {
		name    string
		payload []byte
		status  int
		want    []byte
	}{
		{
			name:    "is invalidate url",
			payload: []byte("[{\"correlation_id\":\"eefbcef4-3940-5a38-b2f0-877152a6d470\",\"original_url\":\"://ya.ru/\"}]"),
			status:  http.StatusBadRequest,
			want:    []byte("Save batch error!"),
		},
		{
			name:    "ok",
			payload: []byte("[{\"correlation_id\":\"eefbcef4-3940-5a38-b2f0-877152a6d470\",\"original_url\":\"https://ya.ru/\"}]"),
			status:  http.StatusCreated,
			want:    []byte("[{\"correlation_id\":\"eefbcef4-3940-5a38-b2f0-877152a6d470\",\"short_url\":\"" + tc.app.URL + "eefbcef4-3940-5a38-b2f0-877152a6d470\"}]\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBuffer(tt.payload)
			req := httptest.NewRequest(http.MethodPost, tc.app.URL+"api/shorten/batch", b)

			w := httptest.NewRecorder()
			tc.app.ShortenAPIBatch(w, req)

			assert.Equal(t, tt.status, w.Code)
			assert.Equal(t, tt.want, w.Body.Bytes())
		})
	}
}

func TestGetURLByID(t *testing.T) {
	tc := New(t)
	tests := []struct {
		name   string
		status int
		id     uuid.UUID
		want   string
	}{
		{
			name:   "id not found",
			status: http.StatusGone,
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
			req := httptest.NewRequest(http.MethodGet, tc.app.URL+tt.id.String(), nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id.String())
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			if tt.status == http.StatusTemporaryRedirect {
				u, _ := url.Parse(tt.want)
				_, _ = tc.app.service.Repository.Save(req.Context(), tt.id, u)
				tc.app.GetURLByID(w, req)

				assert.Equal(t, tt.status, w.Code)
				assert.Equal(t, tt.want, w.Header().Get("Location"))
			} else {
				tc.app.GetURLByID(w, req)

				assert.Equal(t, tt.status, w.Code)
			}
		})
	}
}

func TestGetAllURLByUserID(t *testing.T) {
	tc := New(t)
	tests := []struct {
		name   string
		status int
		userID uuid.UUID
		want   []byte
	}{
		{
			name:   "user id not found",
			status: http.StatusUnauthorized,
			userID: uuid.New(),
			want:   []byte("User unauthorized!"),
		},
		{
			name:   "user not found urls",
			status: http.StatusNoContent,
			userID: uuid.New(),
			want:   []byte("Urls by user id not found!"),
		},
		{
			name:   "ok",
			status: http.StatusOK,
			userID: uuid.New(),
			want:   []byte("[{\"short_url\":\"http://localhost:8080/eefbcef4-3940-5a38-b2f0-877152a6d470\",\"original_url\":\"https://ya.ru/\"}]\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.app.URL+"api/user/urls", nil)

			rctx := chi.NewRouteContext()
			ctx := customContext.SetContextUserID(req.Context(), tt.userID)
			w := httptest.NewRecorder()

			if tt.status == http.StatusUnauthorized {
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			}

			if tt.status == http.StatusNoContent {
				req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			}

			if tt.status == http.StatusOK {
				req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
				u, _ := url.Parse("https://ya.ru/")
				_, _ = tc.app.service.Repository.SaveUser(req.Context(), tt.userID, uuid.NewSHA1(uuid.NameSpaceURL, []byte("https://ya.ru/")), u)
			}

			tc.app.GetAllURLByUserID(w, req)

			assert.Equal(t, tt.status, w.Code)
			assert.Equal(t, tt.want, w.Body.Bytes())
		})
	}
}

func TestDeleteBatchByUserID(t *testing.T) {
	tc := New(t)
	tests := []struct {
		name    string
		payload []byte
		status  int
		want    []byte
	}{
		{
			name:    "body is invalidate",
			payload: []byte("invalid json"),
			status:  http.StatusBadRequest,
			want:    []byte("Body is invalidate!"),
		},
		{
			name:    "ok",
			payload: []byte("[\"eefbcef4-3940-5a38-b2f0-877152a6d470\"]"),
			status:  http.StatusAccepted,
			want:    []byte(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBuffer(tt.payload)
			req := httptest.NewRequest(http.MethodDelete, tc.app.URL+"api/user/urls", b)

			rctx := chi.NewRouteContext()
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pgMock := mock.NewMockRepository(ctrl)
			w := httptest.NewRecorder()

			tc.app.service.Repository = pgMock
			tc.app.DeleteBatchByUserID(w, req)

			assert.Equal(t, tt.status, w.Code)
			if len(tt.want) > 0 {
				assert.Equal(t, tt.want, w.Body.Bytes())
			}
		})
	}
}

func TestPing(t *testing.T) {
	tc := New(t)
	tests := []struct {
		name   string
		status int
	}{
		{
			name:   "database is not active",
			status: http.StatusInternalServerError,
		},
		{
			name:   "ok",
			status: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.app.URL+"ping", nil)

			rctx := chi.NewRouteContext()
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pgMock := mock.NewMockRepository(ctrl)
			w := httptest.NewRecorder()

			if tt.status == http.StatusOK {
				pgMock.EXPECT().
					Ping(gomock.Any()).
					Return(nil).
					AnyTimes()
				tc.app.service.Repository = pgMock

				tc.app.Ping(w, req)

				assert.Equal(t, tt.status, w.Code)
			} else {
				pgMock.EXPECT().
					Ping(gomock.Any()).
					Return(errors.New("database is not active")).
					AnyTimes()
				tc.app.service.Repository = pgMock

				tc.app.Ping(w, req)

				assert.Equal(t, tt.status, w.Code)
			}
		})
	}
}
