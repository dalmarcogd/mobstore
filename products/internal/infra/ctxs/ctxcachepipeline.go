package ctxs

import (
	"context"

	"github.com/go-redis/redis/v8"
)

const (
	xPipelineKey = "xPipelineKey"
)

func ContextWithPipeline(ctx context.Context, tx redis.Pipeliner) context.Context {
	return context.WithValue(ctx, xPipelineKey, tx)
}

func GetPipelineFromContext(ctx context.Context) redis.Pipeliner {
	value := ctx.Value(xPipelineKey)
	if pipe, ok := value.(redis.Pipeliner); ok {
		return pipe
	}
	return nil
}
