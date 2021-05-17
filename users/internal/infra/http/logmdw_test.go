package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/dalmarcogd/mobstore/users/internal/infra/logger"
	"github.com/dalmarcogd/mobstore/users/internal/infra/validator"
	"github.com/dalmarcogd/mobstore/users/internal/services"
)

func TestLoggerMiddlewareSuccessfully(t *testing.T) {
	serviceImpl := New()
	sm := services.New().WithLogger(logger.New()).WithUsersHandler(&mockHandlerSuccessfully{}).WithValidator(validator.New()).WithHttpServer(serviceImpl)

	if err := sm.Init(); err != nil {
		t.Error("unexpected address")
	}

	logMiddleware := LogMiddleware(sm.Logger(), []string{})

	handle := func(ctx echo.Context) error {
		return nil
	}
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"org_id":  null"cid": "{% uuid 'v4' %}" }`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	context := echo.New().NewContext(req, rec)
	context.SetPath("/liveness")
	if err := logMiddleware(handle)(context); err != nil {
		t.Error(err)
	}

	if err := serviceImpl.ServiceManager().Close(); err != nil {
		t.Error(err)
	}
}

func TestLoggerMiddlewareSuccessfullyIgnore(t *testing.T) {
	serviceImpl := New()
	sm := services.New().WithLogger(logger.New()).WithUsersHandler(&mockHandlerSuccessfully{}).WithValidator(validator.New()).WithHttpServer(serviceImpl)

	if err := sm.Init(); err != nil {
		t.Error("unexpected address")
	}

	logMiddleware := LogMiddleware(sm.Logger(), []string{"/liveness"})

	handle := func(ctx echo.Context) error {
		return nil
	}
	req := httptest.NewRequest(http.MethodPost, "/liveness", strings.NewReader(`{"org_id":  null"cid": "{% uuid 'v4' %}" }`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	context := echo.New().NewContext(req, rec)
	context.SetPath("/liveness")
	if err := logMiddleware(handle)(context); err != nil {
		t.Error(err)
	}

	if err := serviceImpl.ServiceManager().Close(); err != nil {
		t.Error(err)
	}
}

func TestLoggerMiddlewareErrorSuccessfully(t *testing.T) {
	serviceImpl := New()
	sm := services.New().WithLogger(logger.New()).WithUsersHandler(&mockHandlerSuccessfully{}).WithValidator(validator.New()).WithHttpServer(serviceImpl)

	if err := sm.Init(); err != nil {
		t.Error("unexpected address")
	}

	logMiddleware := LogMiddleware(sm.Logger(), []string{})

	handle1 := func(ctx echo.Context) error {
		return errors.New("some error")
	}
	handle2 := func(ctx echo.Context) error {
		err := errors.New("some error")
		return echo.NewHTTPError(http.StatusBadRequest, err).SetInternal(err)
	}
	handle3 := func(ctx echo.Context) error {
		err := echo.NewHTTPError(http.StatusInternalServerError, errors.New("some error"))
		return echo.NewHTTPError(http.StatusBadRequest, err).SetInternal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"org_id":  null"cid": "{% uuid 'v4' %}" }`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	if err := logMiddleware(handle1)(echo.New().NewContext(req, rec)); err == nil {
		t.Error("expected error from handle1")
	}
	if err := logMiddleware(handle2)(echo.New().NewContext(req, rec)); err == nil {
		t.Error("expected error from handle2")
	}
	if err := logMiddleware(handle3)(echo.New().NewContext(req, rec)); err == nil {
		t.Error("expected error from handle3")
	}

	if err := serviceImpl.ServiceManager().Close(); err != nil {
		t.Error(err)
	}
}
