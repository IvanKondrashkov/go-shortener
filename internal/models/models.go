package models

import (
	"github.com/google/uuid"
)

type RequestShortenAPI struct {
	URL string `json:"url"`
}

type ResponseShortenAPI struct {
	Result string `json:"result"`
}

type RequestShortenAPIBatch struct {
	CorrelationID uuid.UUID `json:"correlation_id"`
	OriginalURL   string    `json:"original_url"`
}

type ResponseShortenAPIBatch struct {
	CorrelationID uuid.UUID `json:"correlation_id"`
	ShortURL      string    `json:"short_url"`
}

type Event struct {
	ID          uuid.UUID `json:"uuid"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}
