package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/IvanKondrashkov/go-shortener/internal/models"
	customError "github.com/IvanKondrashkov/go-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (app *App) ShortenURL(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain")

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

	id, err := app.service.Save(req.Context(), uuid.NewSHA1(uuid.NameSpaceURL, []byte(u.String())), u)
	if err != nil && errors.Is(err, customError.ErrConflict) {
		res.WriteHeader(http.StatusConflict)
		_, _ = res.Write([]byte(app.URL + id.String()))
		return
	}

	res.WriteHeader(http.StatusCreated)
	_, _ = res.Write([]byte(app.URL + id.String()))
}

func (app *App) ShortenAPI(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

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

	id, err := app.service.Save(req.Context(), uuid.NewSHA1(uuid.NameSpaceURL, []byte(u.String())), u)
	respDto := models.ResponseShortenAPI{
		Result: app.URL + id.String(),
	}
	if err != nil && errors.Is(err, customError.ErrConflict) {
		res.WriteHeader(http.StatusConflict)
		enc := json.NewEncoder(res)
		err = enc.Encode(&respDto)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, _ = res.Write([]byte("Response is invalidate!"))
			return
		}
		return
	}

	res.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(res)
	err = enc.Encode(&respDto)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte("Response is invalidate!"))
		return
	}
}

func (app *App) ShortenAPIBatch(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	var reqDto []*models.RequestShortenAPIBatch
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&reqDto)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte("Body is invalidate!"))
		return
	}
	defer req.Body.Close()

	err = app.service.SaveBatch(req.Context(), reqDto)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte("Save batch error!"))
		return
	}

	respDto, err := models.RequestBatchToResponseBatch(reqDto)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte("Entity mapping is incorrect!"))
		return
	}

	res.WriteHeader(http.StatusCreated)
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

	u, err := app.service.GetByID(req.Context(), uuid.MustParse(id))
	if err != nil && errors.Is(err, customError.ErrNotFound) {
		res.WriteHeader(http.StatusNotFound)
		_, _ = res.Write([]byte("Url by id not found!"))
		return
	}

	if err != nil && errors.Is(err, customError.ErrDeleteAccepted) {
		res.WriteHeader(http.StatusGone)
		_, _ = res.Write([]byte("Delete url accepted!"))
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Location", u.String())
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (app *App) Ping(res http.ResponseWriter, req *http.Request) {
	err := app.service.Ping(req.Context())
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		_, _ = res.Write([]byte("Database is not active!"))
		return
	}

	res.WriteHeader(http.StatusOK)
}
