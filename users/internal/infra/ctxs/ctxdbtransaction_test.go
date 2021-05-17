package ctxs

import (
	"context"
	"database/sql"
	"testing"
)

func TestContextWithTransaction(t *testing.T) {
	type args struct {
		ctx context.Context
		tx  *sql.Tx
	}
	tx1, tx2, tx3 := &sql.Tx{}, &sql.Tx{}, &sql.Tx{}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{name: "with-pipe-value-1", args: args{ctx: context.Background(), tx: tx1}, want: context.WithValue(context.Background(), xTransactionKey, tx1)},
		{name: "with-pipe-value-2", args: args{ctx: context.Background()}, want: context.WithValue(context.Background(), "other-key", tx2)},
		{name: "with-pipe-value-3", args: args{ctx: context.Background(), tx: tx3}, want: context.WithValue(context.Background(), xTransactionKey, tx3)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ContextWithTransaction(tt.args.ctx, tt.args.tx)
			v1, v2 := GetTransactionFromContext(ctx), GetTransactionFromContext(tt.want)
			var s1, s2 *sql.Tx
			if v1 != nil {
				s1 = v1
			}
			if v2 != nil {
				s2 = v2
			}
			if s1 != s2 {
				t.Errorf("ContextWithTransaction() = %v, want %v", s1, tt.want)
			}
		})
	}
}

func TestGetTransactionFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tx1, tx2, tx3 := &sql.Tx{}, &sql.Tx{}, &sql.Tx{}
	tests := []struct {
		name string
		args args
		want *sql.Tx
	}{
		{name: "get-transaction-value-1", args: args{ctx: context.WithValue(context.Background(), xTransactionKey, tx1)}, want: tx1},
		{name: "get-transaction-value-2", args: args{ctx: context.WithValue(context.Background(), "other-key", tx2)}, want: nil},
		{name: "get-transaction-value-3", args: args{ctx: context.WithValue(context.Background(), xTransactionKey, tx3)}, want: tx3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTransactionFromContext(tt.args.ctx)
			var v *sql.Tx
			if got != nil {
				v = got
			}
			if v != tt.want {
				t.Errorf("GetTransactionFromContext() = %v, want %v", got, &tt.want)
			}
		})
	}
}
