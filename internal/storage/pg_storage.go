package storage

import (
	"context"
	"database/sql"

	"github.com/IvanKondrashkov/go-shortener/internal/logger"
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
	return &PgRepositoryImpl{
		Logger: zl,
		conn:   conn,
	}, nil
}

func (pg *PgRepositoryImpl) Ping(ctx context.Context) (err error) {
	if ctx.Err() != nil {
		pg.Logger.Log.Warn("Context is canceled!")
		return
	}
	return pg.conn.PingContext(ctx)
}

func (pg *PgRepositoryImpl) Close() {
	_ = pg.conn.Close()
}
