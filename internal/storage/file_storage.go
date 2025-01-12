package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	customErr "github.com/IvanKondrashkov/go-shortener/internal/errors"
	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
	"github.com/IvanKondrashkov/go-shortener/internal/service"
	"github.com/google/uuid"
)

const (
	Perm = uint32(0666)
)

type FileRepositoryImpl struct {
	service.Repository
	Logger        *logger.ZapLogger
	memRepository service.Repository
	producer      *Producer
	consumer      *Consumer
}

type Producer struct {
	file    io.Writer
	encoder *json.Encoder
}

type Consumer struct {
	file    io.Reader
	decoder *json.Decoder
}

func NewProducer(filePath string) (*Producer, error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.FileMode(Perm))

	if err != nil {
		return nil, fmt.Errorf("open file error: %w", err)
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func NewConsumer(filePath string) (*Consumer, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, os.FileMode(Perm))

	if err != nil {
		return nil, fmt.Errorf("open file error: %w", err)
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func NewFileRepositoryImpl(zl *logger.ZapLogger, memRepository service.Repository, filePath string) (service.Repository, error) {
	p, err := NewProducer(filePath)
	if err != nil {
		return nil, fmt.Errorf("file producer error: %w", err)
	}

	c, err := NewConsumer(filePath)
	if err != nil {
		return nil, fmt.Errorf("file consumer error: %w", err)
	}

	return &FileRepositoryImpl{
		Logger:        zl,
		memRepository: memRepository,
		producer:      p,
		consumer:      c,
	}, nil
}

func (f *FileRepositoryImpl) Save(ctx context.Context, id uuid.UUID, u *url.URL) (uuid.UUID, error) {
	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	var encoder = f.producer.encoder
	event := &models.Event{
		ID:          id,
		ShortURL:    id.String(),
		OriginalURL: u.String(),
	}

	err := encoder.Encode(&event)
	if err != nil {
		return id, fmt.Errorf("serialize error: %w", err)
	}

	u, err = url.Parse(event.OriginalURL)
	if err != nil {
		return id, fmt.Errorf("save in mem storage error: %w", customErr.ErrURLNotValid)
	}

	_, err = f.memRepository.Save(ctx, event.ID, u)
	if err != nil {
		return id, fmt.Errorf("save in mem storage error: %w", err)
	}
	return id, nil
}

func (f *FileRepositoryImpl) SaveBatch(ctx context.Context, batch []*models.RequestShortenAPIBatch) error {
	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	if len(batch) == 0 {
		return fmt.Errorf("save batch in file storage error: %w", customErr.ErrBatchIsEmpty)
	}

	var encoder = f.producer.encoder
	events, _ := models.RequestBatchToEvents(batch)
	for _, event := range events {
		err := encoder.Encode(&event)
		if err != nil {
			return fmt.Errorf("serialize error: %w", err)
		}

		u, err := url.Parse(event.OriginalURL)
		if err != nil {
			return fmt.Errorf("save in mem storage error: %w", customErr.ErrURLNotValid)
		}

		_, err = f.memRepository.Save(ctx, event.ID, u)
		if err != nil && !errors.Is(err, customErr.ErrConflict) {
			return fmt.Errorf("save in mem storage error: %w", err)
		}
	}
	return nil
}

func (f *FileRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*url.URL, error) {
	return f.memRepository.GetByID(ctx, id)
}

func (f *FileRepositoryImpl) ReadFile(ctx context.Context) error {
	var decoder = f.consumer.decoder
	for decoder.More() {
		event := &models.Event{}
		if err := decoder.Decode(&event); err != nil {
			return fmt.Errorf("deserialize error: %w", err)
		}

		u, err := url.Parse(event.OriginalURL)
		if err != nil {
			return fmt.Errorf("save in mem storage error: %w", customErr.ErrURLNotValid)
		}

		_, err = f.memRepository.Save(ctx, event.ID, u)
		if err != nil && !errors.Is(err, customErr.ErrConflict) {
			return fmt.Errorf("save in mem storage error: %w", err)
		}
	}
	return nil
}

func (f *FileRepositoryImpl) Load(ctx context.Context) error {
	err := f.ReadFile(ctx)
	if err != io.EOF && err != nil {
		return fmt.Errorf("read file in file storage error: %w", err)
	}
	return nil
}
