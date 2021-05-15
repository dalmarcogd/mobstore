package ctxs

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
)

func TestContextWithPipeline(t *testing.T) {
	type args struct {
		ctx  context.Context
		pipe redis.Pipeliner
	}
	pipe1, pipe2, pipe3 := &redis.Pipeline{}, &redis.Pipeline{}, &redis.Pipeline{}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{name: "with-pipe-value-1", args: args{ctx: context.Background(), pipe: pipe1}, want: context.WithValue(context.Background(), xPipelineKey, pipe1)},
		{name: "with-pipe-value-2", args: args{ctx: context.Background()}, want: context.WithValue(context.Background(), "other-key", pipe2)},
		{name: "with-pipe-value-3", args: args{ctx: context.Background(), pipe: pipe3}, want: context.WithValue(context.Background(), xPipelineKey, pipe3)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ContextWithPipeline(tt.args.ctx, tt.args.pipe)
			v1, v2 := GetPipelineFromContext(ctx), GetPipelineFromContext(tt.want)
			var s1, s2 redis.Pipeliner
			if v1 != nil {
				s1 = v1
			}
			if v2 != nil {
				s2 = v2
			}
			if s1 != s2 {
				t.Errorf("ContextWithPipeline() = %v, want %v", s1, tt.want)
			}
		})
	}
}

func TestGetPipelineFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	pipe1, pipe2, pipe3 := &redis.Pipeline{}, &redis.Pipeline{}, &redis.Pipeline{}
	tests := []struct {
		name string
		args args
		want redis.Pipeliner
	}{
		{name: "get-pipeline-value-1", args: args{ctx: context.WithValue(context.Background(), xPipelineKey, pipe1)}, want: pipe1},
		{name: "get-pipeline-value-2", args: args{ctx: context.WithValue(context.Background(), "other-key", pipe2)}, want: nil},
		{name: "get-pipeline-value-3", args: args{ctx: context.WithValue(context.Background(), xPipelineKey, pipe3)}, want: pipe3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPipelineFromContext(tt.args.ctx)
			var v redis.Pipeliner
			if got != nil {
				v = got
			}
			if v != tt.want {
				t.Errorf("GetPipelineFromContext() = %v, want %v", got, &tt.want)
			}
		})
	}
}
