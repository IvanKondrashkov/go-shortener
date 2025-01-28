package models

import (
	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/google/uuid"
)

func RequestBatchToEvents(batch []*RequestShortenAPIBatch) ([]*Event, error) {
	res := make([]*Event, 0, len(batch))
	for _, b := range batch {
		event := &Event{
			ID:          b.CorrelationID,
			ShortURL:    uuid.NewSHA1(uuid.NameSpaceURL, []byte(b.OriginalURL)).String(),
			OriginalURL: b.OriginalURL,
		}
		res = append(res, event)
	}
	return res, nil
}

func RequestBatchToResponseBatch(batch []*RequestShortenAPIBatch) ([]*ResponseShortenAPIBatch, error) {
	res := make([]*ResponseShortenAPIBatch, 0, len(batch))
	for _, b := range batch {
		resp := &ResponseShortenAPIBatch{
			CorrelationID: b.CorrelationID,
			ShortURL:      config.URL + uuid.NewSHA1(uuid.NameSpaceURL, []byte(b.OriginalURL)).String(),
		}
		res = append(res, resp)
	}
	return res, nil
}
