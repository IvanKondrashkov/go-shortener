package file

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/service"
)

// Perm определяет права доступа по умолчанию (чтение/запись для владельца и группы)
const (
	Perm = uint32(0666)
)

// Repository реализует файловое хранилище для сервиса сокращения URL.
// Использует JSON кодирование для хранения данных и делегирует in-memory хранилищу.
type Repository struct {
	service.Runner
	service.Repository
	Logger     *logger.ZapLogger  // Логгер для записи событий
	producer   *Producer          // Для записи в файл
	consumer   *Consumer          // Для чтения из файла
	repository service.Repository // In-memory хранилище
}

// Producer реализует запись в файловое хранилище.
type Producer struct {
	file    io.Writer     // Файловый дескриптор для записи
	encoder *json.Encoder // JSON энкодер для сериализации
}

// Consumer реализует чтение из файлового хранилище.
type Consumer struct {
	file    io.Reader     // Файловый дескриптор для чтения
	decoder *json.Decoder // JSON декодер для десериализации
}

// NewProducer создает новый Producer для записи в файловое хранилище.
// Принимает путь к файлу и возвращает Producer или ошибку если файл не может быть открыт.
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

// NewConsumer создает новый Consumer для чтения из файлового хранилища.
// Принимает путь к файлу и возвращает Consumer или ошибку если файл не может быть открыт.
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

// NewRepository создает новый экземпляр файлового хранилища.
// Принимает логгер, in-memory хранилище и путь к файлу.
// Возвращает инициализированный Repository или ошибку если создание producer/consumer не удалось.
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
