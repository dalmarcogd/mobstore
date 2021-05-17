package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/dalmarcogd/mobstore/users/internal/services"
)

func LogMiddleware(log services.Logger, ignorePaths []string) func(h echo.HandlerFunc) echo.HandlerFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			ctx := context.Request().Context()

			path := context.Path()
			shouldIgnore := false
			for _, ignorePath := range ignorePaths {
				if strings.EqualFold(path, ignorePath) {
					shouldIgnore = true
					break
				}
			}
			if !shouldIgnore {
				log.Info(ctx, fmt.Sprintf("Request %v:%v", context.Request().Method, path))
			}
			err := h(context)
			status := context.Response().Status
			if err != nil {
				status = http.StatusInternalServerError
				he, ok := err.(*echo.HTTPError)
				if ok {
					if he.Internal != nil {
						if herr, ok := he.Internal.(*echo.HTTPError); ok {
							he = herr
						}
					}
					status = he.Code
				}
			}
			if !shouldIgnore {
				log.Info(ctx, fmt.Sprintf("Response %v:%v:%v", context.Request().Method, path, status))
			}

			return err
		}
	}
}
