package ctxs

import "context"

const (
	orgIdKey = "orgId"
)

func GetOrgIdFromContext(ctx context.Context) *string {
	value := ctx.Value(orgIdKey)
	if xcid, ok := value.(string); ok {
		return &xcid
	}
	return nil
}

func ContextWithOrgId(ctx context.Context, cid string) context.Context {
	return context.WithValue(ctx, orgIdKey, cid)
}
