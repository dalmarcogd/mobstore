package http

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/propagation/b3"

	"github.com/dalmarcogd/mobstore/products/internal/services"
)

func SpanMiddleware(spans services.Spans, serviceName string, ignorePaths []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()

			path := c.Path()
			shouldIgnore := false
			for _, ignorePath := range ignorePaths {
				if strings.EqualFold(path, ignorePath) {
					shouldIgnore = true
					break
				}
			}
			if !shouldIgnore {
				tracer := spans.Tracer()
				parentSpanContext := tracer.Extract(b3.ExtractHTTP(r))

				name := r.Method + " " + r.URL.Path
				span, ctx := tracer.StartSpanFromContext(r.Context(), name, zipkin.Parent(parentSpanContext))
				defer span.Finish()

				span.Tag("http.method", r.Method)
				span.Tag("http.url", r.URL.String())
				span.Tag("http.request.size", strconv.FormatInt(c.Request().ContentLength, 10))
				span.Tag("component", serviceName)
				span.Tag("serviceName", serviceName)
				span.Tag("service", serviceName)

				c.SetRequest(r.WithContext(ctx))

				if err := next(c); err != nil {
					span.Tag("error", fmt.Sprintf("%v", true))
					c.Error(err)
				}

				span.Tag("http.status_code", strconv.Itoa(c.Response().Status))
				span.Tag("http.response.size", strconv.FormatInt(c.Response().Size, 10))
			} else {
				if err := next(c); err != nil {
					c.Error(err)
				}
			}
			return nil
		}
	}
}
