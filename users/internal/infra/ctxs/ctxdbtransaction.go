package ctxs

import (
	"context"
	"database/sql"
)

const (
	xTransactionKey = "xTransactionKey"
)

func ContextWithTransaction(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, xTransactionKey, tx)
}

func GetTransactionFromContext(ctx context.Context) *sql.Tx {
	value := ctx.Value(xTransactionKey)
	if tx, ok := value.(*sql.Tx); ok {
		return tx
	}
	return nil
}
