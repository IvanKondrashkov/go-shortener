package app

import (
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (app *App) ShortenURL(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Body is invalidate!"))
		return
	}

	u, err := url.Parse(string(body))
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Url is invalidate!"))
		return
	}

	id, err := app.memRepository.Save(uuid.NewSHA1(uuid.NameSpaceURL, []byte(u.String())), u)
	if err != nil {
		res.WriteHeader(http.StatusConflict)
		res.Write([]byte("Id already exists!"))
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(app.BaseURL + id.String()))
}

func (app *App) GetURLByID(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	u, err := app.memRepository.GetByID(uuid.MustParse(id))
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("Url by id not found!"))
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Location", u.String())
	res.WriteHeader(http.StatusTemporaryRedirect)
}
