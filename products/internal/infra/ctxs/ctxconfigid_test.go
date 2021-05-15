package ctxs

import (
	"context"
	"reflect"
	"testing"
)

func TestContextWithConfigId(t *testing.T) {
	type args struct {
		ctx      context.Context
		ConfigId string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{name: "with-ConfigId-value-1", args: args{ctx: context.Background(), ConfigId: "test-1"}, want: context.WithValue(context.Background(), configIdKey, "test-1")},
		{name: "with-ConfigId-value-2", args: args{ctx: context.Background()}, want: context.WithValue(context.Background(), "other-key", "test-1")},
		{name: "with-ConfigId-value-3", args: args{ctx: context.Background(), ConfigId: "test-3"}, want: context.WithValue(context.Background(), configIdKey, "test-3")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ContextWithConfigId(tt.args.ctx, tt.args.ConfigId)
			v1, v2 := GetConfigIdFromContext(ctx), GetConfigIdFromContext(tt.want)
			var s1, s2 string
			if v1 != nil {
				s1 = *v1
			}
			if v2 != nil {
				s2 = *v2
			}
			if s1 != s2 {
				t.Errorf("ContextWithConfigId() = %v, want %v", ctx, tt.want)
			}
		})
	}
}

func TestGetConfigIdFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "get-ConfigId-value-1", args: args{ctx: context.WithValue(context.Background(), configIdKey, "test-1")}, want: "test-1"},
		{name: "get-ConfigId-value-2", args: args{ctx: context.WithValue(context.Background(), "other-key", "test-1")}, want: ""},
		{name: "get-ConfigId-value-3", args: args{ctx: context.WithValue(context.Background(), configIdKey, "test-2")}, want: "test-2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetConfigIdFromContext(tt.args.ctx)
			var v string
			if got != nil {
				v = *got
			}
			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("GetConfigIdFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
