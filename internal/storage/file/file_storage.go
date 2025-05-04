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

func (f *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return nil, nil
}

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

func (f *Repository) GetByID(ctx context.Context, id uuid.UUID) (*url.URL, error) {
	return f.repository.GetByID(ctx, id)
}

func (f *Repository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*models.ResponseShortenAPIUser, error) {
	return f.repository.GetAllByUserID(ctx, userID)
}

func (f *Repository) DeleteBatchByUserID(ctx context.Context, userID uuid.UUID, batch []uuid.UUID) error {
	return f.repository.DeleteBatchByUserID(ctx, userID, batch)
}

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

func (f *Repository) Load(ctx context.Context) error {
	err := f.ReadFile(ctx)
	if err != io.EOF && err != nil {
		return fmt.Errorf("read file in file storage error: %w", err)
	}
	return nil
}
