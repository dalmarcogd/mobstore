package ctxs

import (
	"context"
	"reflect"
	"testing"
)

func TestContextWithUserId(t *testing.T) {
	type args struct {
		ctx      context.Context
		UserId string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{name: "with-UserId-value-1", args: args{ctx: context.Background(), UserId: "test-1"}, want: context.WithValue(context.Background(), userIdKey, "test-1")},
		{name: "with-UserId-value-2", args: args{ctx: context.Background()}, want: context.WithValue(context.Background(), "other-key", "test-1")},
		{name: "with-UserId-value-3", args: args{ctx: context.Background(), UserId: "test-3"}, want: context.WithValue(context.Background(), userIdKey, "test-3")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ContextWithUserId(tt.args.ctx, tt.args.UserId)
			v1, v2 := GetUserIdFromContext(ctx), GetUserIdFromContext(tt.want)
			var s1, s2 string
			if v1 != nil {
				s1 = *v1
			}
			if v2 != nil {
				s2 = *v2
			}
			if s1 != s2 {
				t.Errorf("ContextWithUserId() = %v, want %v", ctx, tt.want)
			}
		})
	}
}

func TestGetUserIdFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "get-UserId-value-1", args: args{ctx: context.WithValue(context.Background(), userIdKey, "test-1")}, want: "test-1"},
		{name: "get-UserId-value-2", args: args{ctx: context.WithValue(context.Background(), "other-key", "test-1")}, want: ""},
		{name: "get-UserId-value-3", args: args{ctx: context.WithValue(context.Background(), userIdKey, "test-2")}, want: "test-2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetUserIdFromContext(tt.args.ctx)
			var v string
			if got != nil {
				v = *got
			}
			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("GetUserIdFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
