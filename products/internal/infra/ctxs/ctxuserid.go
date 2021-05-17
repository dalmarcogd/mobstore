package ctxs

import "context"

const (
	userIdKey = "userId"
)

func GetUserIdFromContext(ctx context.Context) *string {
	value := ctx.Value(userIdKey)
	if xcid, ok := value.(string); ok {
		return &xcid
	}
	return nil
}

func ContextWithUserId(ctx context.Context, cid string) context.Context {
	return context.WithValue(ctx, userIdKey, cid)
}
