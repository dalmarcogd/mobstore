package ctxs

import "context"

const (
	xcidKey = "xcid"
)

func GetCidFromContext(ctx context.Context) *string {
	value := ctx.Value(xcidKey)
	if xcid, ok := value.(string); ok {
		return &xcid
	}
	return nil
}

func ContextWithCid(ctx context.Context, cid string) context.Context {
	return context.WithValue(ctx, xcidKey, cid)
}
