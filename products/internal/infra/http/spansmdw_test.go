package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/openzipkin/zipkin-go"

	"github.com/dalmarcogd/mobstore/products/internal/infra/logger"
	"github.com/dalmarcogd/mobstore/products/internal/infra/validator"
	"github.com/dalmarcogd/mobstore/products/internal/services"
)

func TestSpansMiddlewareSuccessfully(t *testing.T) {
	serviceImpl := New()
	sm := services.New().WithLogger(logger.New()).WithProductsHandler(&mockHandlerSuccessfully{}).WithValidator(validator.New()).WithHttpServer(serviceImpl)

	if err := sm.Init(); err != nil {
		t.Error("unexpected address")
	}

	spanMiddleware := SpanMiddleware(sm.Spans(), sm.Environment().Service(), []string{})

	handle := func(c echo.Context) error {
		ctx := c.Request().Context()
		span := zipkin.SpanFromContext(ctx)
		if span == nil {
			t.Error("span should be present")
		}
		return nil
	}
	handleError := func(c echo.Context) error {
		ctx := c.Request().Context()
		span := zipkin.SpanFromContext(ctx)
		if span == nil {
			t.Error("span should be present")
		}
		return errors.New("some error")
	}

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"org_id":  null"cid": "{% uuid 'v4' %}" }`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	context := echo.New().NewContext(req, rec)
	context.SetPath("/products")
	if err := spanMiddleware(handle)(context); err != nil {
		t.Error(err)
	}
	if err := spanMiddleware(handleError)(context); err != nil {
		t.Error(err)
	}

	if err := serviceImpl.ServiceManager().Close(); err != nil {
		t.Error(err)
	}
}

func TestSpansMiddlewareSuccessfully2(t *testing.T) {
	serviceImpl := New()
	sm := services.New().WithLogger(logger.New()).WithProductsHandler(&mockHandlerSuccessfully{}).WithValidator(validator.New()).WithHttpServer(serviceImpl)

	if err := sm.Init(); err != nil {
		t.Error("unexpected address")
	}

	spanMiddleware := SpanMiddleware(sm.Spans(), sm.Environment().Service(), []string{"/products"})

	handle := func(c echo.Context) error {
		ctx := c.Request().Context()
		span := zipkin.SpanFromContext(ctx)
		if span != nil {
			t.Error("span should not be present")
		}
		return nil
	}
	handleError := func(c echo.Context) error {
		ctx := c.Request().Context()
		span := zipkin.SpanFromContext(ctx)
		if span != nil {
			t.Error("span should not be present")
		}
		return errors.New("some error")
	}

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"org_id":  null"cid": "{% uuid 'v4' %}" }`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	context := echo.New().NewContext(req, rec)
	context.SetPath("/products")
	if err := spanMiddleware(handle)(context); err != nil {
		t.Error(err)
	}
	if err := spanMiddleware(handleError)(context); err != nil {
		t.Error(err)
	}

	if err := serviceImpl.ServiceManager().Close(); err != nil {
		t.Error(err)
	}
}
