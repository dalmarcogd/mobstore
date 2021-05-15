package ctxs

import (
	"context"
	"reflect"
	"testing"
)

func TestContextWithCid(t *testing.T) {
	type args struct {
		ctx context.Context
		cid string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{name: "with-cid-value-1", args: args{ctx: context.Background(), cid: "test-1"}, want: context.WithValue(context.Background(), xcidKey, "test-1")},
		{name: "with-cid-value-2", args: args{ctx: context.Background()}, want: context.WithValue(context.Background(), "other-key", "test-1")},
		{name: "with-cid-value-3", args: args{ctx: context.Background(), cid: "test-3"}, want: context.WithValue(context.Background(), xcidKey, "test-3")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ContextWithCid(tt.args.ctx, tt.args.cid)
			v1, v2 := GetCidFromContext(ctx), GetCidFromContext(tt.want)
			var s1, s2 string
			if v1 != nil {
				s1 = *v1
			}
			if v2 != nil {
				s2 = *v2
			}
			if s1 != s2 {
				t.Errorf("ContextWithCid() = %v, want %v", ctx, tt.want)
			}
		})
	}
}

func TestGetCidFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "get-cid-value-1", args: args{ctx: context.WithValue(context.Background(), xcidKey, "test-1")}, want: "test-1"},
		{name: "get-cid-value-2", args: args{ctx: context.WithValue(context.Background(), "other-key", "test-1")}, want: ""},
		{name: "get-cid-value-3", args: args{ctx: context.WithValue(context.Background(), xcidKey, "test-2")}, want: "test-2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCidFromContext(tt.args.ctx)
			var v string
			if got != nil {
				v = *got
			}
			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("GetCidFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
