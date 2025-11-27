package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type ctxKeyTx struct{}

var key ctxKeyTx = ctxKeyTx{}

func txToCtx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, key, tx)
}

func txFromCtx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(key).(pgx.Tx); ok {
		return tx
	}
	return nil
}
