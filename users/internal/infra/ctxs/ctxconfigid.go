package ctxs

import "context"

const (
	configIdKey = "configId"
)

func GetConfigIdFromContext(ctx context.Context) *string {
	value := ctx.Value(configIdKey)
	if xcid, ok := value.(string); ok {
		return &xcid
	}
	return nil
}

func ContextWithConfigId(ctx context.Context, cid string) context.Context {
	return context.WithValue(ctx, configIdKey, cid)
}
