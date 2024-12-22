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

type Event struct {
	ID          uuid.UUID `json:"uuid"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}
