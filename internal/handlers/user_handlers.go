// Package handlers содержит HTTP-хендлеры для API
package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/IvanKondrashkov/go-shortener/internal/models"
	"github.com/IvanKondrashkov/go-shortener/internal/service"
	customContext "github.com/IvanKondrashkov/go-shortener/internal/service/middleware/auth"
	"github.com/google/uuid"
)

// GetAllURLByUserID возвращает все URL пользователя
// @Summary Получить URL пользователя
// @Description Возвращает все сокращенные URL, созданные текущим пользователем
// @Tags Пользователь
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.ResponseShortenAPIUser
// @Success 204 "Нет сохраненных URL"
// @Failure 401 {string} string "Пользователь не авторизован"
// @Router /api/user/urls [get]
func (app *App) GetAllURLByUserID(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	respDto, err := app.service.GetAllByUserID(req.Context())
	if err != nil && errors.Is(err, service.ErrUserUnauthorized) {
		res.WriteHeader(http.StatusUnauthorized)
		_, _ = res.Write([]byte("User unauthorized!"))
		return
	}

	if len(respDto) == 0 {
		res.WriteHeader(http.StatusNoContent)
		_, _ = res.Write([]byte("Urls by user id not found!"))
		return
	}

	writer := writerPool.Get().(*bufio.Writer)
	writer.Reset(res)
	defer func() {
		writer.Flush()
		writerPool.Put(writer)
	}()

	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(respDto); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte("Response is invalidate!"))
		return
	}
}

// DeleteBatchByUserID удаляет список URL пользователя
// @Summary Удалить URL пользователя
// @Description Помечает указанные URL как удаленные (асинхронная операция)
// @Tags Пользователь
// @Security ApiKeyAuth
// @Accept json
// @Param input body []uuid.UUID true "Список ID URL для удаления"
// @Success 202 "Запрос на удаление принят"
// @Failure 400 {string} string "Неверный формат запроса"
// @Router /api/user/urls [delete]
func (app *App) DeleteBatchByUserID(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	reader := readerPool.Get().(*bufio.Reader)
	reader.Reset(req.Body)
	defer readerPool.Put(reader)

	var reqDto []uuid.UUID
	if err := json.NewDecoder(reader).Decode(&reqDto); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, _ = res.Write([]byte("Body is invalidate!"))
		return
	}

	event := models.DeleteEvent{
		Batch:  reqDto,
		UserID: customContext.GetContextUserID(req.Context()),
	}

	go app.worker.SendDeleteBatchRequest(context.Background(), event)
	res.WriteHeader(http.StatusAccepted)
}
