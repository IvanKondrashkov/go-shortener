package storage

import (
	"encoding/json"
	"io"
	"net/url"
	"os"

	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
)

const (
	Perm = uint32(0666)
)

type FileRepositoryImpl struct {
	Logger        *logger.ZapLogger
	memRepository *MemRepositoryImpl
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
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func NewConsumer(filePath string) (*Consumer, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, os.FileMode(Perm))

	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func NewFileRepositoryImpl(zl *logger.ZapLogger, memRepository *MemRepositoryImpl, filePath string) (*FileRepositoryImpl, error) {
	p, err := NewProducer(filePath)
	if err != nil {
		return nil, err
	}

	c, err := NewConsumer(filePath)
	if err != nil {
		return nil, err
	}

	return &FileRepositoryImpl{
		Logger:        zl,
		memRepository: memRepository,
		producer:      p,
		consumer:      c,
	}, nil
}

func (f *FileRepositoryImpl) WriteFile(event *models.Event) (err error) {
	var encoder = f.producer.encoder
	return encoder.Encode(&event)
}

func (f *FileRepositoryImpl) WriteFileBatch(events []*models.Event) (err error) {
	if len(events) == 0 {
		return err
	}

	var encoder = f.producer.encoder
	for _, event := range events {
		err := encoder.Encode(&event)
		if err != nil {
			return err
		}
	}
	return err
}

func (f *FileRepositoryImpl) ReadFile() (err error) {
	var decoder = f.consumer.decoder
	for decoder.More() {
		event := &models.Event{}
		if err := decoder.Decode(&event); err != nil {
			return err
		}

		u, err := url.Parse(event.OriginalURL)
		if err != nil {
			return err
		}

		_, err = f.memRepository.Save(event.ID, u)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FileRepositoryImpl) Load() (err error) {
	err = f.ReadFile()
	if err != io.EOF && err != nil {
		return err
	}
	return nil
}
