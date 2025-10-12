// Package models содержит модели для API
package models

import (
	"github.com/google/uuid"
)

// RequestShortenAPI запрос на сокращение URL
// @Description Запрос на создание сокращенного URL
type RequestShortenAPI struct {
	URL string `json:"url" example:"https://mail.ru/"`
}

// ResponseShortenAPI ответ с сокращенным URL
// @Description Сокращенный URL
type ResponseShortenAPI struct {
	Result string `json:"result" example:"https://localhost:8080/614904e8-28e8-5fec-aef9-56c6713e3107"`
}

// RequestShortenAPIBatch элемент пакетного запроса на сокращение
// @Description Элемент пакетного запроса на сокращение URL
type RequestShortenAPIBatch struct {
	CorrelationID uuid.UUID `json:"correlation_id" example:"eefbcef4-3940-5a38-b2f0-877152a6d471"`
	OriginalURL   string    `json:"original_url" example:"https://mail.ru/"`
}

// ResponseShortenAPIBatch элемент пакетного ответа с сокращенным URL
// @Description Элемент пакетного ответа с сокращенным URL
type ResponseShortenAPIBatch struct {
	CorrelationID uuid.UUID `json:"correlation_id" example:"eefbcef4-3940-5a38-b2f0-877152a6d471"`
	ShortURL      string    `json:"short_url" example:"https://localhost:8080/614904e8-28e8-5fec-aef9-56c6713e3107"`
}

// ResponseShortenAPIUser элемент ответа с URL пользователя
// @Description Информация о сокращенном URL пользователя
type ResponseShortenAPIUser struct {
	ShortURL    string `json:"short_url" example:"https://localhost:8080/614904e8-28e8-5fec-aef9-56c6713e3107"`
	OriginalURL string `json:"original_url" example:"https://mail.ru/"`
}

// Event элемент события для записи в файловое хранилище
// @Description Информация о сокращенном URL пользователя
type Event struct {
	ID          uuid.UUID `json:"uuid" example:"eefbcef4-3940-5a38-b2f0-877152a6d471"`
	ShortURL    string    `json:"short_url" example:"https://localhost:8080/614904e8-28e8-5fec-aef9-56c6713e3107"`
	OriginalURL string    `json:"original_url" example:"https://mail.ru/"`
}

// DeleteEvent элемент события для удаления батча URL пользователя
// @Description Информация об удаляемых URL пользователя
type DeleteEvent struct {
	UserID *uuid.UUID
	Batch  []uuid.UUID
}

// Stats представляет статистику сервиса
// @Description Информация о кол-ве пользователей в системе и сокращенных URL
type Stats struct {
	URLs  int `json:"urls" example:"10"`
	Users int `json:"users" example:"5"`
}
