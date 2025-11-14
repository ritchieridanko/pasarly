package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritchieridanko/pasarly/auth-service/internal/shared/ce"
)

type Database struct {
	pool *pgxpool.Pool
}

func NewDatabase(p *pgxpool.Pool) *Database {
	return &Database{pool: p}
}

func (d *Database) Execute(ctx context.Context, query string, args ...any) error {
	e := d.executor(ctx)
	res, err := e.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if ra := res.RowsAffected(); ra == 0 {
		return ce.ErrDBAffectNoRows
	}

	return nil
}

func (d *Database) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	e := d.executor(ctx)
	return e.QueryRow(ctx, query, args...)
}

func (d *Database) InTx(ctx context.Context) bool {
	return txFromCtx(ctx) != nil
}
