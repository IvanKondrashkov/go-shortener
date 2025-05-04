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

func (m *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return nil, nil
}

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
