package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/IvanKondrashkov/go-shortener/internal/handlers/mock"
	customContext "github.com/IvanKondrashkov/go-shortener/internal/service/middleware/auth"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllURLByUserID(t *testing.T) {
	tc := NewSuite(t)
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
				_, _ = tc.app.service.Repository.SaveUser(req.Context(), nil, tt.userID, uuid.NewSHA1(uuid.NameSpaceURL, []byte("https://ya.ru/")), u)
			}

			tc.app.GetAllURLByUserID(w, req)

			assert.Equal(t, tt.status, w.Code)
			assert.Equal(t, tt.want, w.Body.Bytes())
		})
	}
}

func TestDeleteBatchByUserID(t *testing.T) {
	tc := NewSuite(t)
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
