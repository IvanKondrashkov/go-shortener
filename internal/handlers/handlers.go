package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type repository interface {
	Save(id uuid.UUID, url *url.URL) (res uuid.UUID, err error)
	GetByID(id uuid.UUID) (res *url.URL, err error)
}

type fileRepository interface {
	WriteFile(event *models.Event) (err error)
	ReadFile() (err error)
	Load() (err error)
}

type pgRepository interface {
	Ping(ctx context.Context) (err error)
	Save(ctx context.Context, id uuid.UUID, url *url.URL) (err error)
}

func (app *App) ShortenURL(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte("Body is invalidate!"))
		return
	}
	defer req.Body.Close()

	u, err := url.Parse(string(body))
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte("Url is invalidate!"))
		return
	}

	id, err := app.repository.Save(uuid.NewSHA1(uuid.NameSpaceURL, []byte(u.String())), u)
	if err != nil {
		res.WriteHeader(http.StatusConflict)
		_, _ = res.Write([]byte("Entity conflict!"))
		return
	}

	if config.DatabaseDsn != "" {
		err = app.pgRepository.Save(req.Context(), uuid.NewSHA1(uuid.NameSpaceURL, []byte(u.String())), u)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, _ = res.Write([]byte("Entity save database is incorrect!"))
			return
		}
	}

	event := &models.Event{
		ID:          id,
		ShortURL:    id.String(),
		OriginalURL: u.String(),
	}

	if config.FileStoragePath != "" {
		err = app.fileRepository.WriteFile(event)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, _ = res.Write([]byte("Write file is incorrect!"))
			return
		}
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	_, _ = res.Write([]byte(app.URL + id.String()))
}

func (app *App) ShortenAPI(res http.ResponseWriter, req *http.Request) {
	var reqDto models.RequestShortenAPI
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&reqDto)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte("Body is invalidate!"))
		return
	}
	defer req.Body.Close()

	u, err := url.Parse(reqDto.URL)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte("Url is invalidate!"))
		return
	}

	id, err := app.repository.Save(uuid.NewSHA1(uuid.NameSpaceURL, []byte(u.String())), u)
	if err != nil {
		res.WriteHeader(http.StatusConflict)
		_, _ = res.Write([]byte("Entity conflict!"))
		return
	}

	if config.DatabaseDsn != "" {
		err = app.pgRepository.Save(req.Context(), uuid.NewSHA1(uuid.NameSpaceURL, []byte(u.String())), u)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, _ = res.Write([]byte("Entity save database is incorrect!"))
			return
		}
	}

	event := &models.Event{
		ID:          id,
		ShortURL:    id.String(),
		OriginalURL: u.String(),
	}

	if config.FileStoragePath != "" {
		err = app.fileRepository.WriteFile(event)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, _ = res.Write([]byte("Write file is incorrect!"))
			return
		}
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

	respDto := models.ResponseShortenAPI{
		Result: app.URL + id.String(),
	}
	enc := json.NewEncoder(res)
	err = enc.Encode(&respDto)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte("Response is invalidate!"))
		return
	}
}

func (app *App) GetURLByID(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	u, err := app.repository.GetByID(uuid.MustParse(id))
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		_, _ = res.Write([]byte("Url by id not found!"))
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Location", u.String())
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (app *App) Ping(res http.ResponseWriter, req *http.Request) {
	err := app.pgRepository.Ping(req.Context())
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		_, _ = res.Write([]byte("Database is not active!"))
		return
	}

	res.WriteHeader(http.StatusOK)
}
