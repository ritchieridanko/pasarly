package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type executor interface {
	Exec(ctx context.Context, query string, args ...any) (ct pgconn.CommandTag, err error)
	Query(ctx context.Context, query string, args ...any) (rows pgx.Rows, err error)
	QueryRow(ctx context.Context, query string, args ...any) (row pgx.Row)
}

func (d *Database) executor(ctx context.Context) executor {
	if tx := txFromCtx(ctx); tx != nil {
		return tx
	}
	return d.pool
}
