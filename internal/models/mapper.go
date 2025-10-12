package models

import (
	"github.com/IvanKondrashkov/go-shortener/internal/config"
	pb "github.com/IvanKondrashkov/go-shortener/internal/proto"
	"github.com/google/uuid"
)

// RequestBatchToEvents маппер для преобразования []*RequestShortenAPIBatch в []*Event
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

// RequestBatchUserToEvents маппер для преобразования []*RequestShortenAPIBatch в []*Event
func RequestBatchUserToEvents(userID uuid.UUID, batch []*RequestShortenAPIBatch) ([]*Event, error) {
	res := make([]*Event, 0, len(batch))
	for _, b := range batch {
		event := &Event{
			ID:          userID,
			ShortURL:    uuid.NewSHA1(uuid.NameSpaceURL, []byte(b.OriginalURL)).String(),
			OriginalURL: b.OriginalURL,
		}
		res = append(res, event)
	}
	return res, nil
}

// RequestBatchToResponseBatch маппер для преобразования []*RequestShortenAPIBatch в []*ResponseShortenAPIBatch
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

// RequestBatchGrpcToRequestShortenAPIBatch маппер для преобразования []*pb.BatchRequest в []*RequestShortenAPIBatch
func RequestBatchGrpcToRequestShortenAPIBatch(batch []*pb.BatchRequest) ([]*RequestShortenAPIBatch, error) {
	res := make([]*RequestShortenAPIBatch, 0, len(batch))
	for _, b := range batch {
		resp := &RequestShortenAPIBatch{
			CorrelationID: uuid.MustParse(b.GetCorrelationId()),
			OriginalURL:   b.GetOriginalUrl(),
		}
		res = append(res, resp)
	}
	return res, nil
}

// RequestBatchGrpcToResponseBatchGrpc маппер для преобразования []*pb.BatchRequest в []*pb.BatchResponse
func RequestBatchGrpcToResponseBatchGrpc(batch []*pb.BatchRequest) ([]*pb.BatchResponse, error) {
	res := make([]*pb.BatchResponse, 0, len(batch))
	for _, b := range batch {
		shortURL := config.URL + uuid.NewSHA1(uuid.NameSpaceURL, []byte(b.GetOriginalUrl())).String()

		var resp = &pb.BatchResponse{
			CorrelationId: b.CorrelationId,
			ShortUrl:      &shortURL,
		}
		res = append(res, resp)
	}
	return res, nil
}

// ResponseShortenAPIUserToURLGrpc маппер для преобразования []*ResponseShortenAPIUser в []*pb.URL
func ResponseShortenAPIUserToURLGrpc(batch []*ResponseShortenAPIUser) ([]*pb.URL, error) {
	res := make([]*pb.URL, 0, len(batch))
	for _, b := range batch {
		shortURL := config.URL + uuid.NewSHA1(uuid.NameSpaceURL, []byte(b.OriginalURL)).String()
		originalURL := b.OriginalURL

		var resp = &pb.URL{
			ShortUrl:    &shortURL,
			OriginalUrl: &originalURL,
		}
		res = append(res, resp)
	}
	return res, nil
}
