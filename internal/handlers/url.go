// Package handlers содержит HTTP-хендлеры для API
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

// ShortenURL обрабатывает запрос на сокращение URL
// @Summary Сократить URL
// @Description Создает короткую версию переданного URL
// @Tags URL
// @Accept plain
// @Produce plain
// @Param url body string true "Оригинальный URL для сокращения"
// @Success 201 {string} string "Сокращенный URL"
// @Success 409 {string} string "URL уже был сокращен ранее"
// @Failure 400 {string} string "Неверный формат URL"
// @Router / [post]
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

// ShortenAPI обрабатывает JSON запрос на сокращение URL
// @Summary Сократить URL (JSON)
// @Description Создает короткую версию переданного URL (JSON формат)
// @Tags URL
// @Accept json
// @Produce json
// @Param input body models.RequestShortenAPI true "Запрос на сокращение URL"
// @Success 201 {object} models.ResponseShortenAPI
// @Success 409 {object} models.ResponseShortenAPI
// @Failure 400 {string} string "Неверный формат запроса"
// @Router /api/shorten [post]
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

// ShortenAPIBatch обрабатывает пакетный запрос на сокращение URL
// @Summary Пакетное сокращение URL
// @Description Создает короткие версии для списка URL
// @Tags URL
// @Accept json
// @Produce json
// @Param input body []models.RequestShortenAPIBatch true "Список URL для сокращения"
// @Success 201 {object} []models.ResponseShortenAPIBatch
// @Failure 400 {string} string "Неверный формат запроса"
// @Router /api/shorten/batch [post]
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

// GetURLByID возвращает оригинальный URL по его ID
// @Summary Получить оригинальный URL
// @Description Перенаправляет на оригинальный URL по его сокращенному ID
// @Tags URL
// @Param id path string true "ID сокращенного URL"
// @Success 307 "Перенаправление на оригинальный URL"
// @Failure 404 {string} string "URL не найден"
// @Failure 410 {string} string "URL был удален"
// @Router /{id} [get]
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

// Ping проверяет доступность базы данных
// @Summary Проверка состояния
// @Description Проверяет соединение с базой данных
// @Tags Сервис
// @Success 200 "База данных доступна"
// @Failure 500 {string} string "База данных недоступна"
// @Router /ping [get]
func (app *App) Ping(res http.ResponseWriter, req *http.Request) {
	err := app.service.Ping(req.Context())
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		_, _ = res.Write([]byte("Database is not active!"))
		return
	}

	res.WriteHeader(http.StatusOK)
}
