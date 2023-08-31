package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	DB      *pgxpool.Pool
	Builder squirrel.StatementBuilderType
}

func New(cfg *Config) (*Postgres, error) {
	const fn = "db.postgres.New"

	pg := &Postgres{}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	pgxConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	pgxConfig.MaxConns = cfg.MaxConns

	pg.DB, err = pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	err = pg.DB.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return pg, nil
}

func (p *Postgres) IsEqualErrors(err error, pgErrCode string) bool {
	var pgErr *pgconn.PgError
	if ok := errors.As(err, &pgErr); ok {
		if pgErr.Code == pgErrCode {
			return true
		}
	}

	return false
}

func (p *Postgres) Close() {
	if p.DB != nil {
		p.DB.Close()
	}
}
