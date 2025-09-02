// Package models содержит модели для API
package models

import (
	"github.com/google/uuid"
)

// RequestShortenAPI запрос на сокращение URL
// @Description Запрос на создание сокращенного URL
type RequestShortenAPI struct {
	URL string `json:"url"`
}

// ResponseShortenAPI ответ с сокращенным URL
// @Description Сокращенный URL
type ResponseShortenAPI struct {
	Result string `json:"result"`
}

// RequestShortenAPIBatch элемент пакетного запроса на сокращение
// @Description Элемент пакетного запроса на сокращение URL
type RequestShortenAPIBatch struct {
	CorrelationID uuid.UUID `json:"correlation_id"`
	OriginalURL   string    `json:"original_url"`
}

// ResponseShortenAPIBatch элемент пакетного ответа с сокращенным URL
// @Description Элемент пакетного ответа с сокращенным URL
type ResponseShortenAPIBatch struct {
	CorrelationID uuid.UUID `json:"correlation_id"`
	ShortURL      string    `json:"short_url"`
}

// ResponseShortenAPIUser элемент ответа с URL пользователя
// @Description Информация о сокращенном URL пользователя
type ResponseShortenAPIUser struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// Event элемент события для записи в файловое хранилище
// @Description Информация о сокращенном URL пользователя
type Event struct {
	ID          uuid.UUID `json:"uuid"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}

// DeleteEvent элемент события для удаления батча URL пользователя
// @Description Информация об удаляемых URL пользователя
type DeleteEvent struct {
	UserID *uuid.UUID
	Batch  []uuid.UUID
}
