package models

import (
	"github.com/IvanKondrashkov/go-shortener/internal/config"
	pb "github.com/IvanKondrashkov/go-shortener/internal/proto"
	"github.com/google/uuid"
)

// RequestBatchToEvents маппер для преобразования RequestShortenAPIBatch в Event.
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

// RequestBatchUserToEvents маппер для преобразования RequestShortenAPIBatch в Event, пользователя.
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

// RequestBatchToResponseBatch маппер для преобразования RequestShortenAPIBatch в ResponseShortenAPIBatch.
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

// RequestBatchGrpcToRequestShortenAPIBatch маппер для преобразования BatchRequest в RequestShortenAPIBatch.
func RequestBatchGrpcToRequestShortenAPIBatch(batch []*pb.BatchRequest) ([]*RequestShortenAPIBatch, error) {
	res := make([]*RequestShortenAPIBatch, 0, len(batch))
	for _, b := range batch {
		resp := &RequestShortenAPIBatch{
			CorrelationID: uuid.MustParse(b.CorrelationId),
			OriginalURL:   b.OriginalUrl,
		}
		res = append(res, resp)
	}
	return res, nil
}

// RequestBatchGrpcToResponseBatchGrpc маппер для преобразования BatchRequest в BatchResponse.
func RequestBatchGrpcToResponseBatchGrpc(batch []*pb.BatchRequest) ([]*pb.BatchResponse, error) {
	res := make([]*pb.BatchResponse, 0, len(batch))
	for _, b := range batch {
		var resp = &pb.BatchResponse{
			CorrelationId: b.CorrelationId,
			ShortUrl:      config.URL + uuid.NewSHA1(uuid.NameSpaceURL, []byte(b.OriginalUrl)).String(),
		}
		res = append(res, resp)
	}
	return res, nil
}

// ResponseShortenAPIUserToURLGrpc маппер для преобразования ResponseShortenAPIUser в URL.
func ResponseShortenAPIUserToURLGrpc(batch []*ResponseShortenAPIUser) ([]*pb.URL, error) {
	res := make([]*pb.URL, 0, len(batch))
	for _, b := range batch {
		var resp = &pb.URL{
			ShortUrl:    config.URL + uuid.NewSHA1(uuid.NameSpaceURL, []byte(b.OriginalURL)).String(),
			OriginalUrl: b.OriginalURL,
		}
		res = append(res, resp)
	}
	return res, nil
}
