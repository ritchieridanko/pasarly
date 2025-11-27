package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritchieridanko/pasarly/backend/shared/ce"
	"go.opentelemetry.io/otel"
)

type Transactor struct {
	pool *pgxpool.Pool
}

func NewTransactor(p *pgxpool.Pool) *Transactor {
	return &Transactor{pool: p}
}

func (t *Transactor) WithTx(ctx context.Context, fn func(context.Context) *ce.Error) *ce.Error {
	ctx, span := otel.Tracer("database.transactor").Start(ctx, "WithTx")
	defer span.End()

	tx := txFromCtx(ctx)
	isNewTx := false

	var err error
	if tx == nil {
		tx, err = t.pool.Begin(ctx)
		if err != nil {
			e := fmt.Errorf("failed to begin database transaction: %w", err)
			return ce.NewError(span, ce.CodeDBTx, ce.MsgInternalServer, e)
		}
		ctx = txToCtx(ctx, tx)
		isNewTx = true
	}

	if err := fn(ctx); err != nil {
		if isNewTx {
			_ = tx.Rollback(ctx)
		}
		return err
	}

	if isNewTx {
		if err := tx.Commit(ctx); err != nil {
			e := fmt.Errorf("failed to commit database transaction: %w", err)
			return ce.NewError(span, ce.CodeDBTx, ce.MsgInternalServer, e)
		}
	}

	return nil
}
