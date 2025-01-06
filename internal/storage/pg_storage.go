package storage

import (
	"context"
	"database/sql"
	"errors"
	"net/url"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/IvanKondrashkov/go-shortener/internal/logger"
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
	ctx, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

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

func (pg *PgRepositoryImpl) Close() {
	_ = pg.conn.Close()
}
