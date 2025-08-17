package handlers

import (
	"bufio"
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

	reader := readerPool.Get().(*bufio.Reader)
	reader.Reset(req.Body)
	defer readerPool.Put(reader)

	body, err := io.ReadAll(reader)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte("Body is invalidate!"))
		return
	}

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

	reader := readerPool.Get().(*bufio.Reader)
	reader.Reset(req.Body)
	defer readerPool.Put(reader)

	var reqDto models.RequestShortenAPI
	if err := json.NewDecoder(reader).Decode(&reqDto); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte("Body is invalidate!"))
		return
	}

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

	writer := writerPool.Get().(*bufio.Writer)
	writer.Reset(res)
	defer func() {
		writer.Flush()
		writerPool.Put(writer)
	}()

	if err != nil && errors.Is(err, customError.ErrConflict) {
		res.WriteHeader(http.StatusConflict)
		if err := json.NewEncoder(writer).Encode(respDto); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, _ = res.Write([]byte("Response is invalidate!"))
			return
		}
		return
	}

	res.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(writer).Encode(respDto); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte("Response is invalidate!"))
		return
	}
}

func (app *App) ShortenAPIBatch(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	reader := readerPool.Get().(*bufio.Reader)
	reader.Reset(req.Body)
	defer readerPool.Put(reader)

	var reqDto []*models.RequestShortenAPIBatch
	if err := json.NewDecoder(reader).Decode(&reqDto); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte("Body is invalidate!"))
		return
	}

	err := app.service.SaveBatch(req.Context(), reqDto)
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

	writer := writerPool.Get().(*bufio.Writer)
	writer.Reset(res)
	defer func() {
		writer.Flush()
		writerPool.Put(writer)
	}()

	res.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(writer).Encode(respDto); err != nil {
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
