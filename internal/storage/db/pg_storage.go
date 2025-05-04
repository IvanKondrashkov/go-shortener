package db

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

func (pg *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return pg.pool.Begin(ctx)
}

func (pg *Repository) Save(ctx context.Context, tx pgx.Tx, id uuid.UUID, u *url.URL) (uuid.UUID, error) {
	query := `
	INSERT INTO urls(short_url, original_url)
	VALUES ($1, $2)
	ON CONFLICT (short_url) DO UPDATE
	SET
	short_url = EXCLUDED.short_url,
	original_url = EXCLUDED.original_url;
	`

	_, err := tx.Exec(ctx, query, id, u.String())
	if err != nil {
		return id, fmt.Errorf("save in pg storage error: %w", err)
	}
	return id, nil
}

func (pg *Repository) SaveUser(ctx context.Context, tx pgx.Tx, userID, id uuid.UUID, u *url.URL) (uuid.UUID, error) {
	query := `
	INSERT INTO urls(short_url, user_id, original_url)
	VALUES ($1, $2, $3)
	ON CONFLICT (short_url) DO UPDATE
	SET
	short_url = EXCLUDED.short_url,
	user_id = EXCLUDED.user_id,
	original_url = EXCLUDED.original_url;
	`

	_, err := tx.Exec(ctx, query, id, userID, u.String())
	if err != nil {
		return id, fmt.Errorf("save in pg storage error: %w", err)
	}
	return id, nil
}

func (pg *Repository) SaveBatch(ctx context.Context, batch []*models.RequestShortenAPIBatch) error {
	if len(batch) == 0 {
		return fmt.Errorf("save batch in pg storage error: %w", customError.ErrBatchIsEmpty)
	}

	valuesShortURL := make([]uuid.UUID, 0, len(batch))
	valuesOriginalURL := make([]string, 0, len(batch))
	for _, b := range batch {
		valuesShortURL = append(valuesShortURL, uuid.NewSHA1(uuid.NameSpaceURL, []byte(b.OriginalURL)))
		valuesOriginalURL = append(valuesOriginalURL, b.OriginalURL)
	}

	query := `
	INSERT INTO urls(short_url, original_url)
	VALUES (UNNEST($1::UUID[]), UNNEST($2::VARCHAR[]))
	ON CONFLICT (short_url) DO NOTHING;
	`

	b := &pgx.Batch{}
	b.Queue(query, valuesShortURL, valuesOriginalURL)

	err := pg.pool.SendBatch(ctx, b).Close()
	if err != nil {
		return fmt.Errorf("save batch in pg storage error: %w", err)
	}
	return nil
}

func (pg *Repository) SaveBatchUser(ctx context.Context, userID uuid.UUID, batch []*models.RequestShortenAPIBatch) error {
	if len(batch) == 0 {
		return fmt.Errorf("save batch in pg storage error: %w", customError.ErrBatchIsEmpty)
	}

	valuesShortURL := make([]uuid.UUID, 0, len(batch))
	valuesOriginalURL := make([]string, 0, len(batch))
	for _, b := range batch {
		valuesShortURL = append(valuesShortURL, uuid.NewSHA1(uuid.NameSpaceURL, []byte(b.OriginalURL)))
		valuesOriginalURL = append(valuesOriginalURL, b.OriginalURL)
	}

	query := `
	INSERT INTO urls(short_url, user_id, original_url)
	VALUES (UNNEST($1::UUID[]), $2, UNNEST($3::VARCHAR[]))
	ON CONFLICT (short_url) DO NOTHING;
	`

	b := &pgx.Batch{}
	b.Queue(query, valuesShortURL, userID, valuesOriginalURL)

	err := pg.pool.SendBatch(ctx, b).Close()
	if err != nil {
		return fmt.Errorf("save batch in pg storage error: %w", err)
	}
	return nil
}

func (pg *Repository) GetByID(ctx context.Context, id uuid.UUID) (*url.URL, error) {
	query := `
	SELECT original_url, is_deleted
	FROM urls
	WHERE short_url = $1;
	`

	var isDeleted *bool
	var originalURL string
	err := pg.pool.QueryRow(ctx, query, id).Scan(&originalURL, &isDeleted)
	if err != nil {
		return nil, fmt.Errorf("get in pg storage error: %w", customError.ErrNotFound)
	}

	u, err := url.Parse(originalURL)
	if err != nil {
		return nil, fmt.Errorf("get in pg storage error: %w", customError.ErrURLNotValid)
	}

	if isDeleted != nil && *isDeleted {
		return nil, fmt.Errorf("get in pg storage error: %w", customError.ErrDeleteAccepted)
	}
	return u, nil
}

func (pg *Repository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*models.ResponseShortenAPIUser, error) {
	query := `
	SELECT short_url, original_url
	FROM urls
	WHERE user_id = $1;
	`

	rows, err := pg.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get all in pg storage error: %w", err)
	}
	defer rows.Close()

	var urls []*models.ResponseShortenAPIUser
	for rows.Next() {
		var u models.ResponseShortenAPIUser
		if err = rows.Scan(&u.ShortURL, &u.OriginalURL); err != nil {
			return urls, fmt.Errorf("get all in pg storage error: %w", err)
		}
		u.ShortURL = config.URL + u.ShortURL
		urls = append(urls, &u)
	}
	return urls, nil
}

func (pg *Repository) DeleteBatchByUserID(ctx context.Context, userID uuid.UUID, batch []uuid.UUID) error {
	if len(batch) == 0 {
		return fmt.Errorf("delete batch in pg storage error: %w", customError.ErrBatchIsEmpty)
	}

	valuesShortURL := make([]uuid.UUID, 0, len(batch))
	valuesShortURL = append(valuesShortURL, batch...)

	query := `
	UPDATE urls SET is_deleted = true WHERE short_url = ANY($1) AND user_id = $2;
	`

	conn, err := pg.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("delete batch in pg storage error: %w", err)
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, query, valuesShortURL, userID)
	if err != nil {
		return fmt.Errorf("delete batch in pg storage error: %w", err)
	}
	return nil
}

func (pg *Repository) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	return pg.pool.Ping(ctx)
}

func (pg *Repository) Close() {
	pg.pool.Close()
}
