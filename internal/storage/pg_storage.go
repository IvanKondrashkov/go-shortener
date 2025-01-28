package storage

import (
	"context"
	"database/sql"
	"errors"
	"net/url"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PgRepositoryImpl struct {
	Logger *logger.ZapLogger
	conn   *sql.DB
}

func NewPgRepositoryImpl(zl *logger.ZapLogger, dns string) (*PgRepositoryImpl, error) {
	conn, err := sql.Open("pgx", dns)
	if err != nil {
		return nil, err
	}

	m, err := migrate.New("file://internal/db/migration", dns)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return &PgRepositoryImpl{
		Logger: zl,
		conn:   conn,
	}, nil
}

func (pg *PgRepositoryImpl) Ping(ctx context.Context) (err error) {
	ctx, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	return pg.conn.PingContext(ctx)
}

func (pg *PgRepositoryImpl) Save(ctx context.Context, id uuid.UUID, u *url.URL) (err error) {
	tx, err := pg.conn.Begin()
	if err != nil {
		return err
	}

	query := `
	INSERT INTO urls(short_url, original_url)
	VALUES ($1, $2)
	ON CONFLICT (short_url) DO UPDATE
	SET
	short_url = EXCLUDED.short_url,
	original_url = EXCLUDED.original_url;
	`

	_, err = tx.ExecContext(ctx, query, id, u.String())
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pg *PgRepositoryImpl) SaveBatch(ctx context.Context, batch []*models.RequestShortenAPIBatch) (err error) {
	if len(batch) == 0 {
		return err
	}

	valuesShortURL := make([]uuid.UUID, 0, len(batch))
	valuesOriginalURL := make([]string, 0, len(batch))

	for _, b := range batch {
		valuesShortURL = append(valuesShortURL, uuid.NewSHA1(uuid.NameSpaceURL, []byte(b.OriginalURL)))
		valuesOriginalURL = append(valuesOriginalURL, b.OriginalURL)
	}

	tx, err := pg.conn.Begin()
	if err != nil {
		return err
	}

	query := `
	INSERT INTO urls(short_url, original_url)
	VALUES (UNNEST($1::UUID[]), UNNEST($2::VARCHAR[]))
	ON CONFLICT (short_url) DO NOTHING;
	`

	_, err = tx.ExecContext(ctx, query, valuesShortURL, valuesOriginalURL)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pg *PgRepositoryImpl) Close() {
	_ = pg.conn.Close()
}
