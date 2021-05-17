package http

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/dalmarcogd/mobstore/products/internal/infra/ctxs"
	"github.com/dalmarcogd/mobstore/products/internal/infra/reqgetter"
)

func CidContextMiddleware() func(h echo.HandlerFunc) echo.HandlerFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctxEcho echo.Context) error {
			request := ctxEcho.Request()
			ctx := request.Context()
			cid := reqgetter.GetCid(request)
			if cid == nil {
				newUUIDStr := uuid.New().String()
				cid = &newUUIDStr
			}
			ctx = ctxs.ContextWithCid(ctx, *cid)
			ctxEcho.SetRequest(request.WithContext(ctx))
			err := h(ctxEcho)
			return err
		}
	}
}
