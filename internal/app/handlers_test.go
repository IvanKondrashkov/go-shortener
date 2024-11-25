package app

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/IvanKondrashkov/go-shortener/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL(t *testing.T) {
	baseURL := "http://localhost:8080/"
	app := &App{
		memRepository: storage.NewMemRepositoryImpl(),
	}

	tests := []struct {
		name     string
		status   int
		url      string
		response []byte
	}{
		{
			name:     "is invalidate url",
			status:   http.StatusBadRequest,
			url:      "://ya.ru/",
			response: []byte("Url is invalidate!"),
		},
		{
			name:     "ok",
			status:   http.StatusCreated,
			url:      "https://ya.ru/",
			response: []byte(baseURL + uuid.NewSHA1(uuid.NameSpaceURL, []byte("https://ya.ru/")).String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewBuffer([]byte(tt.url))
			request := httptest.NewRequest(http.MethodPost, baseURL, body)
			w := httptest.NewRecorder()
			app.ShortenURL(w, request)

			assert.Equal(t, tt.status, w.Code)
			assert.Equal(t, tt.response, w.Body.Bytes())
		})
	}
}

func TestGetURLByID(t *testing.T) {
	baseURL := "http://localhost:8080/"
	app := &App{
		memRepository: storage.NewMemRepositoryImpl(),
	}

	tests := []struct {
		name     string
		status   int
		id       uuid.UUID
		response string
	}{
		{
			name:     "id not found",
			status:   http.StatusNotFound,
			id:       uuid.New(),
			response: "Url by id not found!",
		},
		{
			name:     "ok",
			id:       uuid.NewSHA1(uuid.NameSpaceURL, []byte("https://ya.ru/")),
			status:   http.StatusTemporaryRedirect,
			response: "https://ya.ru/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, baseURL+tt.id.String(), nil)
			request.SetPathValue("id", tt.id.String())
			w := httptest.NewRecorder()

			if tt.status == http.StatusTemporaryRedirect {
				u, _ := url.Parse(tt.response)
				app.memRepository.Save(tt.id, u)
				app.GetURLByID(w, request)

				assert.Equal(t, tt.status, w.Code)
				assert.Equal(t, tt.response, w.Header().Get("Location"))
			} else {
				app.GetURLByID(w, request)

				assert.Equal(t, tt.status, w.Code)
			}

		})
	}
}
