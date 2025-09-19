package handlers

import (
	"bufio"
	"encoding/json"
	"net/http"
)

// GetStats возвращает статистику сервиса
// @Summary Получить статистику
// @Description Возвращает количество URL и пользователей в сервисе
// @Tags Сервис
// @Security ApiKeyAuth
// @Produce json
// @Param X-Real-IP header string true "IP адрес клиента"
// @Success 200 {object} models.Stats
// @Failure 403 {string} string "Доступ запрещен"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/internal/stats [get]
func (app *App) GetStats(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	respDto, err := app.service.GetStats(req.Context())
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		_, _ = res.Write([]byte("Internal server error"))
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
		res.WriteHeader(http.StatusInternalServerError)
		_, _ = res.Write([]byte("Response is invalidate!"))
		return
	}
}
