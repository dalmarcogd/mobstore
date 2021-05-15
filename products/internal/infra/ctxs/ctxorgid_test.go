package ctxs

import (
	"context"
	"reflect"
	"testing"
)

func TestContextWithOrgId(t *testing.T) {
	type args struct {
		ctx   context.Context
		OrgId string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{name: "with-OrgId-value-1", args: args{ctx: context.Background(), OrgId: "test-1"}, want: context.WithValue(context.Background(), orgIdKey, "test-1")},
		{name: "with-OrgId-value-2", args: args{ctx: context.Background()}, want: context.WithValue(context.Background(), "other-key", "test-1")},
		{name: "with-OrgId-value-3", args: args{ctx: context.Background(), OrgId: "test-3"}, want: context.WithValue(context.Background(), orgIdKey, "test-3")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ContextWithOrgId(tt.args.ctx, tt.args.OrgId)
			v1, v2 := GetOrgIdFromContext(ctx), GetOrgIdFromContext(tt.want)
			var s1, s2 string
			if v1 != nil {
				s1 = *v1
			}
			if v2 != nil {
				s2 = *v2
			}
			if s1 != s2 {
				t.Errorf("ContextWithOrgId() = %v, want %v", ctx, tt.want)
			}
		})
	}
}

func TestGetOrgIdFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "get-OrgId-value-1", args: args{ctx: context.WithValue(context.Background(), orgIdKey, "test-1")}, want: "test-1"},
		{name: "get-OrgId-value-2", args: args{ctx: context.WithValue(context.Background(), "other-key", "test-1")}, want: ""},
		{name: "get-OrgId-value-3", args: args{ctx: context.WithValue(context.Background(), orgIdKey, "test-2")}, want: "test-2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetOrgIdFromContext(tt.args.ctx)
			var v string
			if got != nil {
				v = *got
			}
			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("GetOrgIdFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
