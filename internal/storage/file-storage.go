package storage

import (
	"encoding/json"
	"io"
	"net/url"
	"os"

	"github.com/IvanKondrashkov/go-shortener/internal/models"
)

const (
	basePerm = uint32(0666)
)

type FileRepositoryImpl struct {
	memRepository *MemRepositoryImpl
	producer      *Producer
	consumer      *Consumer
}

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewProducer(filePath string) (*Producer, error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.FileMode(basePerm))

	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func NewConsumer(filePath string) (*Consumer, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, os.FileMode(basePerm))

	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func NewFileRepositoryImpl(memRepository *MemRepositoryImpl, filePath string) (*FileRepositoryImpl, error) {
	p, err := NewProducer(filePath)
	if err != nil {
		return nil, err
	}

	c, err := NewConsumer(filePath)
	if err != nil {
		return nil, err
	}

	return &FileRepositoryImpl{
		memRepository: memRepository,
		producer:      p,
		consumer:      c,
	}, nil
}

func (f *FileRepositoryImpl) WriteFile(event *models.Event) (err error) {
	var encoder = f.producer.encoder
	return encoder.Encode(&event)
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

func (f *FileRepositoryImpl) Load() error {
	err := f.ReadFile()
	if err != io.EOF && err != nil {
		return err
	}
	return nil
}
