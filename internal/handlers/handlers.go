package handlers

import (
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type repository interface {
	Save(id uuid.UUID, url *url.URL) (res uuid.UUID, err error)
	GetByID(id uuid.UUID) (res *url.URL, err error)
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

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	_, _ = res.Write([]byte(app.BaseURL + id.String()))
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
