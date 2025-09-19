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

// BeginTx начинает новую транзакцию в базе данных.
// Возвращает pgx.Tx транзакцию или ошибку если транзакция не может быть начата.
func (pg *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return pg.pool.Begin(ctx)
}

// Save сохраняет URL в PostgreSQL базе данных.
// Возвращает UUID сохраненного URL или ошибку если операция не удалась.
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

// SaveUser сохраняет URL в PostgreSQL базе данных, ассоциированный с пользователем.
// Возвращает UUID сохраненного URL или ошибку если операция не удалась.
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

// SaveBatch сохраняет несколько URL в PostgreSQL базе данных одной операцией.
// Возвращает ErrBatchIsEmpty если batch пуст.
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

// SaveBatchUser сохраняет несколько URL в PostgreSQL базе данных, ассоциированных с пользователем.
// Возвращает ErrBatchIsEmpty если batch пуст.
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

// GetByID получает URL из PostgreSQL базы данных по его UUID ключу.
// Возвращает ErrNotFound если ключ не существует или ErrDeleteAccepted если URL был удален.
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

// GetAllByUserID получает все URL, ассоциированные с пользователем, из PostgreSQL базы данных.
// Возвращает срез URL или ошибку если запрос не удался.
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

// DeleteBatchByUserID помечает несколько URL как удаленные для пользователя в PostgreSQL базе данных.
// Возвращает ErrBatchIsEmpty если batch пуст или ошибку если операция не удалась.
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

// Ping проверяет соединение с базой данных.
// Возвращает ошибку если соединение не может быть установлено.
func (pg *Repository) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	return pg.pool.Ping(ctx)
}

// Close освобождает ресурсы соединения с базой данных.
func (pg *Repository) Close() {
	pg.pool.Close()
}

// GetStats получить статистику сервиса
// Принимает:
// - ctx: контекст
// Возвращает:
// - статистику сервиса *models.Stats
// - ошибку, если запрос не удался
func (pg *Repository) GetStats(ctx context.Context) (*models.Stats, error) {
	queryURLs := `SELECT COUNT(DISTINCT original_url) FROM urls WHERE is_deleted = false;`
	queryUsers := `SELECT COUNT(DISTINCT user_id) FROM urls;`

	var urlsCount, usersCount int

	err := pg.pool.QueryRow(ctx, queryURLs).Scan(&urlsCount)
	if err != nil {
		return nil, fmt.Errorf("get urls count error: %w", err)
	}

	err = pg.pool.QueryRow(ctx, queryUsers).Scan(&usersCount)
	if err != nil {
		return nil, fmt.Errorf("get users count error: %w", err)
	}

	return &models.Stats{
		URLs:  urlsCount,
		Users: usersCount,
	}, nil
}
