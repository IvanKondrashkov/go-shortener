package mem

import (
	"context"
	"fmt"
	"net/url"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
	customError "github.com/IvanKondrashkov/go-shortener/internal/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// BeginTx начинает новую транзакцию (заглушка для in-memory хранилища).
// Возвращает nil транзакцию и nil ошибку, так как in-memory хранилище не поддерживает транзакции.
func (m *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return nil, nil
}

// Save сохраняет URL в in-memory хранилище с указанным UUID в качестве ключа.
// Возвращает UUID и ErrConflict, если ключ уже существует.
func (m *Repository) Save(ctx context.Context, tx pgx.Tx, id uuid.UUID, u *url.URL) (uuid.UUID, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	_, ok := m.memRepository[id]
	if ok {
		m.memRepository[id] = u
		return id, fmt.Errorf("save in mem storage error: %w", customError.ErrConflict)
	}

	m.memRepository[id] = u
	return id, nil
}

// SaveUser сохраняет URL в in-memory хранилище, ассоциированный с конкретным пользователем.
// Возвращает UUID и ErrConflict, если ключ уже существует для этого пользователя.
func (m *Repository) SaveUser(ctx context.Context, tx pgx.Tx, userID, id uuid.UUID, u *url.URL) (uuid.UUID, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	_, ok := m.userRepository[userID]
	if !ok {
		m.userRepository[userID] = make(map[uuid.UUID]*url.URL)
	}

	_, ok = m.userRepository[userID][id]
	if ok {
		m.memRepository[id] = u
		m.userRepository[userID][id] = u
		return id, fmt.Errorf("save in mem storage error: %w", customError.ErrConflict)
	}

	m.memRepository[id] = u
	m.userRepository[userID][id] = u
	return id, nil
}

// SaveBatch сохраняет несколько URL в in-memory хранилище одной операцией.
// Возвращает ErrBatchIsEmpty если batch пуст или ErrURLNotValid если какой-то URL невалиден.
func (m *Repository) SaveBatch(ctx context.Context, batch []*models.RequestShortenAPIBatch) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	if len(batch) == 0 {
		return fmt.Errorf("save batch in mem storage error: %w", customError.ErrBatchIsEmpty)
	}

	for _, b := range batch {
		u, err := url.Parse(b.OriginalURL)
		if err != nil {
			return fmt.Errorf("save batch in mem storage error: %w", customError.ErrURLNotValid)
		}
		m.memRepository[uuid.NewSHA1(uuid.NameSpaceURL, []byte(u.String()))] = u
	}
	return nil
}

// SaveBatchUser сохраняет несколько URL в in-memory хранилище, ассоциированных с пользователем.
// Возвращает ErrBatchIsEmpty если batch пуст или ErrURLNotValid если какой-то URL невалиден.
func (m *Repository) SaveBatchUser(ctx context.Context, userID uuid.UUID, batch []*models.RequestShortenAPIBatch) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	if len(batch) == 0 {
		return fmt.Errorf("save batch in mem storage error: %w", customError.ErrBatchIsEmpty)
	}

	_, ok := m.userRepository[userID]
	if !ok {
		m.userRepository[userID] = make(map[uuid.UUID]*url.URL)
	}

	for _, b := range batch {
		u, err := url.Parse(b.OriginalURL)
		if err != nil {
			return fmt.Errorf("save batch in mem storage error: %w", customError.ErrURLNotValid)
		}
		m.memRepository[uuid.NewSHA1(uuid.NameSpaceURL, []byte(u.String()))] = u
		m.userRepository[userID][uuid.NewSHA1(uuid.NameSpaceURL, []byte(u.String()))] = u
	}
	return nil
}

// GetByID получает URL из in-memory хранилища по его UUID ключу.
// Возвращает ErrNotFound если ключ не существует или ErrDeleteAccepted если URL был удален.
func (m *Repository) GetByID(ctx context.Context, id uuid.UUID) (*url.URL, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	u, ok := m.memRepository[id]
	if !ok && u != nil {
		return nil, fmt.Errorf("get in mem storage error: %w", customError.ErrNotFound)
	}

	if u == nil {
		return nil, fmt.Errorf("get in mem storage error: %w", customError.ErrDeleteAccepted)
	}
	return u, nil
}

// GetAllByUserID получает все URL, ассоциированные с конкретным пользователем.
// Возвращает ErrNotFound если у пользователя нет сохраненных URL.
func (m *Repository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*models.ResponseShortenAPIUser, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	urls, ok := m.userRepository[userID]
	if !ok {
		return nil, fmt.Errorf("get all in mem storage error: %w", customError.ErrNotFound)
	}

	res := make([]*models.ResponseShortenAPIUser, 0, len(urls))
	for k, v := range urls {
		u := models.ResponseShortenAPIUser{
			ShortURL:    config.URL + k.String(),
			OriginalURL: v.String(),
		}
		res = append(res, &u)
	}
	return res, nil
}

// DeleteBatchByUserID помечает несколько URL как удаленные для конкретного пользователя.
// Возвращает ErrNotFound если у пользователя нет сохраненных URL.
func (m *Repository) DeleteBatchByUserID(ctx context.Context, userID uuid.UUID, batch []uuid.UUID) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	urls, ok := m.userRepository[userID]
	if !ok {
		return fmt.Errorf("delete batch in mem storage error: %w", customError.ErrNotFound)
	}

	for _, b := range batch {
		urls[b] = nil
		m.memRepository[b] = nil
	}
	return nil
}

// GetStats получить статистику сервиса
// Принимает:
// - ctx: контекст
// Возвращает:
// - статистику сервиса *models.Stats
// - ошибку, если запрос не удался
func (m *Repository) GetStats(ctx context.Context) (*models.Stats, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	urlsCount := len(m.memRepository)
	usersCount := len(m.userRepository)

	return &models.Stats{
		URLs:  urlsCount,
		Users: usersCount,
	}, nil
}
