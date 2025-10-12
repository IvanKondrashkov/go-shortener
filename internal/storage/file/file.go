package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
	customError "github.com/IvanKondrashkov/go-shortener/internal/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// BeginTx начинает новую транзакцию (заглушка для файлового хранилища)
// Возвращает nil транзакцию и nil ошибку, так как файловое хранилище не поддерживает транзакции
func (f *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return nil, nil
}

// Commit сохраняет изменение в базе данных (заглушка для файлового хранилища)
// Возвращает nil ошибку, так как файловое хранилище не поддерживает транзакции
func (f *Repository) Commit(ctx context.Context, tx pgx.Tx) error {
	return nil
}

// Rollback откатывает транзакцию в базе данных (заглушка для файлового хранилища)
// Возвращает nil ошибку, так как файловое хранилище не поддерживает транзакции
func (f *Repository) Rollback(ctx context.Context, tx pgx.Tx) error {
	return nil
}

// Save сохраняет URL, сохранение в файловое хранилище и in-memory хранилище
// Принимает:
// - ctx: контекст с информацией о пользователе
// - tx: транзакцию
// - id: короткий URL
// - u: оригинальный URL
// Возвращает:
// - id записи или ошибку, если ключ уже существует (ErrConflict)
func (f *Repository) Save(ctx context.Context, tx pgx.Tx, id uuid.UUID, u *url.URL) (uuid.UUID, error) {
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
		return id, fmt.Errorf("save in mem storage error: %w", customError.ErrURLNotValid)
	}

	_, err = f.repository.Save(ctx, tx, event.ID, u)
	if err != nil {
		return id, fmt.Errorf("save in mem storage error: %w", err)
	}
	return id, nil
}

// SaveUser сохраняет URL ассоциированный с пользователем, сохранение в файловое хранилище и in-memory хранилище
// Принимает:
// - ctx: контекст с информацией о пользователе
// - tx: транзакцию
// - userID: идентификатор пользователя
// - id: короткий URL
// - u: оригинальный URL
// Возвращает:
// - id записи или ошибку, если ключ уже существует (ErrConflict)
func (f *Repository) SaveUser(ctx context.Context, tx pgx.Tx, userID, id uuid.UUID, u *url.URL) (uuid.UUID, error) {
	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	var encoder = f.producer.encoder
	event := &models.Event{
		ID:          userID,
		ShortURL:    id.String(),
		OriginalURL: u.String(),
	}

	err := encoder.Encode(&event)
	if err != nil {
		return id, fmt.Errorf("serialize error: %w", err)
	}

	u, err = url.Parse(event.OriginalURL)
	if err != nil {
		return id, fmt.Errorf("save in mem storage error: %w", customError.ErrURLNotValid)
	}

	_, err = f.repository.SaveUser(ctx, tx, event.ID, uuid.MustParse(event.ShortURL), u)
	if err != nil {
		return id, fmt.Errorf("save in mem storage error: %w", err)
	}
	return id, nil
}

// SaveBatch сохраняет массив URL одной операцией, сохранение в файловое хранилище и in-memory хранилище
// Принимает:
// - ctx: контекст с информацией о пользователе
// - batch: массив URL
// Возвращает:
// - ошибку, если batch пуст (ErrBatchIsEmpty) или если URL не валиден (ErrURLNotValid)
func (f *Repository) SaveBatch(ctx context.Context, batch []*models.RequestShortenAPIBatch) error {
	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	if len(batch) == 0 {
		return fmt.Errorf("save batch in file storage error: %w", customError.ErrBatchIsEmpty)
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
			return fmt.Errorf("save in mem storage error: %w", customError.ErrURLNotValid)
		}

		_, err = f.repository.Save(ctx, nil, event.ID, u)
		if err != nil && !errors.Is(err, customError.ErrConflict) {
			return fmt.Errorf("save in mem storage error: %w", err)
		}
	}
	return nil
}

// SaveBatchUser сохраняет массив URL ассоциированный с пользователем одной операцией, сохранение в файловое хранилище и in-memory хранилище
// Принимает:
// - ctx: контекст с информацией о пользователе
// - userID: идентификатор пользователя
// - batch: массив URL
// Возвращает:
// - ошибку, если batch пуст (ErrBatchIsEmpty) или если URL не валиден (ErrURLNotValid)
func (f *Repository) SaveBatchUser(ctx context.Context, userID uuid.UUID, batch []*models.RequestShortenAPIBatch) error {
	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	if len(batch) == 0 {
		return fmt.Errorf("save batch in file storage error: %w", customError.ErrBatchIsEmpty)
	}

	var encoder = f.producer.encoder
	events, _ := models.RequestBatchUserToEvents(userID, batch)
	for _, event := range events {
		err := encoder.Encode(&event)
		if err != nil {
			return fmt.Errorf("serialize error: %w", err)
		}

		u, err := url.Parse(event.OriginalURL)
		if err != nil {
			return fmt.Errorf("save in mem storage error: %w", customError.ErrURLNotValid)
		}

		_, err = f.repository.SaveUser(ctx, nil, event.ID, uuid.MustParse(event.ShortURL), u)
		if err != nil && !errors.Is(err, customError.ErrConflict) {
			return fmt.Errorf("save in mem storage error: %w", err)
		}
	}
	return nil
}

// GetByID получает URL, получение данных из in-memory хранилища
// Принимает:
// - ctx: контекст с информацией о пользователе
// - id: короткий URL
// Возвращает:
// - Возвращает оригинальный URL или ошибку, если URL не найден (ErrNotFound) или если URL был удален (ErrDeleteAccepted)
func (f *Repository) GetByID(ctx context.Context, id uuid.UUID) (*url.URL, error) {
	return f.repository.GetByID(ctx, id)
}

// GetAllByUserID получает все URL ассоциированные с пользователем, получение данных из in-memory хранилища
// Принимает:
// - ctx: контекст с информацией о пользователе
// - userID: идентификатор пользователя
// Возвращает:
// - []*models.ResponseShortenAPIUser или ошибку, если возникли проблемы при получении данных
func (f *Repository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*models.ResponseShortenAPIUser, error) {
	return f.repository.GetAllByUserID(ctx, userID)
}

// DeleteBatchByUserID помечает несколько URL как удаленные для пользователя, удаление данных из in-memory хранилища
// Принимает:
// - ctx: контекст с информацией о пользователе
// - userID: идентификатор пользователя
// - batch: массив коротких URL
// Возвращает:
// - ошибку, если batch пуст (ErrBatchIsEmpty) или возникли проблемы при удалении данных
func (f *Repository) DeleteBatchByUserID(ctx context.Context, userID uuid.UUID, batch []uuid.UUID) error {
	return f.repository.DeleteBatchByUserID(ctx, userID, batch)
}

// Ping проверяет соединение с базой данных
// Принимает:
// - ctx: контекст
// Возвращает:
// - ошибку, если возникли проблемы при получении данных
func (f *Repository) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	return f.repository.Ping(ctx)
}

// GetStats получить статистику сервиса
// Принимает:
// - ctx: контекст
// Возвращает:
// - статистику сервиса *models.Stats
// - ошибку, если запрос не удался
func (f *Repository) GetStats(ctx context.Context) (*models.Stats, error) {
	return f.repository.GetStats(ctx)
}

// ReadFile читает URL из файлового хранилища и загружает их в память
// Принимает:
// - ctx: контекст
// Возвращает:
// - ошибку, если возникли проблемы при получении данных
func (f *Repository) ReadFile(ctx context.Context) error {
	var decoder = f.consumer.decoder
	for decoder.More() {
		event := &models.Event{}
		if err := decoder.Decode(&event); err != nil {
			return fmt.Errorf("deserialize error: %w", err)
		}

		u, err := url.Parse(event.OriginalURL)
		if err != nil {
			return fmt.Errorf("save in mem storage error: %w", customError.ErrURLNotValid)
		}

		_, err = f.repository.Save(ctx, nil, event.ID, u)
		if err != nil && !errors.Is(err, customError.ErrConflict) {
			return fmt.Errorf("save in mem storage error: %w", err)
		}

		_, err = f.repository.SaveUser(ctx, nil, event.ID, uuid.MustParse(event.ShortURL), u)
		if err != nil && !errors.Is(err, customError.ErrConflict) {
			return fmt.Errorf("save in mem storage error: %w", err)
		}
	}
	return nil
}

// Load инициализирует хранилище, читая данные из файлового хранилища.
// Принимает:
// - ctx: контекст
// Возвращает:
// - ошибку, если чтение файла не удалось
func (f *Repository) Load(ctx context.Context) error {
	err := f.ReadFile(ctx)
	if err != io.EOF && err != nil {
		return fmt.Errorf("read file in file storage error: %w", err)
	}
	return nil
}
