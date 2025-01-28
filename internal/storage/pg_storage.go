package storage

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	customErr "github.com/IvanKondrashkov/go-shortener/internal/errors"
	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
	"github.com/IvanKondrashkov/go-shortener/internal/service"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PgRepositoryImpl struct {
	service.Repository
	Logger *logger.ZapLogger
	conn   *pgx.Conn
}

func NewPgRepositoryImpl(ctx context.Context, zl *logger.ZapLogger, dns string) (*PgRepositoryImpl, error) {
	conn, err := pgx.Connect(ctx, dns)
	if err != nil {
		return nil, fmt.Errorf("open database connection error: %w", err)
	}

	m, err := migrate.New("file://migration", dns)
	if err != nil {
		return nil, fmt.Errorf("database migration error: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("database migration error: %w", err)
	}

	return &PgRepositoryImpl{
		Logger: zl,
		conn:   conn,
	}, nil
}

func (pg *PgRepositoryImpl) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	return pg.conn.Ping(ctx)
}

func (pg *PgRepositoryImpl) Save(ctx context.Context, id uuid.UUID, u *url.URL) (uuid.UUID, error) {
	tx, err := pg.conn.Begin(ctx)
	if err != nil {
		return id, fmt.Errorf("open transactional error: %w", err)
	}

	query := `
	INSERT INTO urls(short_url, original_url)
	VALUES ($1, $2)
	ON CONFLICT (short_url) DO UPDATE
	SET
	short_url = EXCLUDED.short_url,
	original_url = EXCLUDED.original_url;
	`

	_, err = tx.Exec(ctx, query, id, u.String())
	if err != nil {
		_ = tx.Rollback(ctx)
		return id, fmt.Errorf("save in pg storage error: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return id, fmt.Errorf("commit transactional error: %w", err)
	}
	return id, nil
}

func (pg *PgRepositoryImpl) SaveBatch(ctx context.Context, batch []*models.RequestShortenAPIBatch) error {
	if len(batch) == 0 {
		return fmt.Errorf("save batch in pg storage error: %w", customErr.ErrBatchIsEmpty)
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

	err := pg.conn.SendBatch(ctx, b).Close()
	if err != nil {
		return fmt.Errorf("save batch in pg storage error: %w", err)
	}
	return nil
}

func (pg *PgRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*url.URL, error) {
	tx, err := pg.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("open transactional error: %w", err)
	}

	query := `
	SELECT original_url
	FROM urls
	WHERE short_url = $1;
	`

	var row string
	err = tx.QueryRow(ctx, query, id).Scan(&row)
	if err != nil {
		_ = tx.Rollback(ctx)
		return nil, fmt.Errorf("get in pg storage error: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit transactional error: %w", err)
	}

	u, err := url.Parse(row)
	if err != nil {
		return nil, fmt.Errorf("get in pg storage error: %w", customErr.ErrURLNotValid)
	}
	return u, nil
}

func (pg *PgRepositoryImpl) Close(ctx context.Context) {
	_ = pg.conn.Close(ctx)
}
