package models

import (
	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/google/uuid"
)

func RequestBatchToEvents(batch []*RequestShortenAPIBatch) (res []*Event, err error) {
	for _, b := range batch {
		event := &Event{
			ID:          b.CorrelationID,
			ShortURL:    uuid.NewSHA1(uuid.NameSpaceURL, []byte(b.OriginalURL)).String(),
			OriginalURL: b.OriginalURL,
		}
		res = append(res, event)
	}
	return res, err
}

func RequestBatchToResponseBatch(batch []*RequestShortenAPIBatch) (res []*ResponseShortenAPIBatch, err error) {
	for _, b := range batch {
		resp := &ResponseShortenAPIBatch{
			CorrelationID: b.CorrelationID,
			ShortURL:      config.URL + uuid.NewSHA1(uuid.NameSpaceURL, []byte(b.OriginalURL)).String(),
		}
		res = append(res, resp)
	}
	return res, err
}
