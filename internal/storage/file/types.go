package file

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/service"
)

const (
	Perm = uint32(0666)
)

type Repository struct {
	service.Runner
	service.Repository
	Logger     *logger.ZapLogger
	producer   *Producer
	consumer   *Consumer
	repository service.Repository
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

func NewRepository(zl *logger.ZapLogger, r service.Repository, filePath string) (*Repository, error) {
	p, err := NewProducer(filePath)
	if err != nil {
		return nil, fmt.Errorf("file producer error: %w", err)
	}

	c, err := NewConsumer(filePath)
	if err != nil {
		return nil, fmt.Errorf("file consumer error: %w", err)
	}

	return &Repository{
		Logger:     zl,
		producer:   p,
		consumer:   c,
		repository: r,
	}, nil
}
