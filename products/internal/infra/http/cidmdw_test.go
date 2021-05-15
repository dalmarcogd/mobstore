package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/dalmarcogd/mobstore/products/internal/infra/ctxs"
	"github.com/dalmarcogd/mobstore/products/internal/services"
)

func TestCidMiddlewareSuccessfully(t *testing.T) {
	serviceImpl := New()
	sm := services.New().WithProductsHandler(&mockHandlerSuccessfully{}).WithHttpServer(serviceImpl)

	if err := sm.Init(); err != nil {
		t.Error("unexpected address")
	}

	cidMiddleware := CidContextMiddleware()

	withCid := false
	handle := func(ctx echo.Context) error {
		cid := ctxs.GetCidFromContext(ctx.Request().Context())
		withCid = cid != nil && *cid == "some-id"
		return nil
	}
	req := httptest.NewRequest(http.MethodPost, "/cards", strings.NewReader(`{"org_id":  null"cid": "{% uuid 'v4' %}" }`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("x-cid", "some-id")
	rec := httptest.NewRecorder()

	if err := cidMiddleware(handle)(echo.New().NewContext(req, rec)); err != nil {
		t.Error(err)
	}

	if !withCid {
		t.Error("expected context from request with cid")
	}

	if err := serviceImpl.ServiceManager().Close(); err != nil {
		t.Error(err)
	}
}

func TestCidMiddlewareErrorSuccessfully(t *testing.T) {
	serviceImpl := New()
	sm := services.New().WithProductsHandler(&mockHandlerSuccessfully{}).WithHttpServer(serviceImpl)

	if err := sm.Init(); err != nil {
		t.Error("unexpected address")
	}

	cidMiddleware := CidContextMiddleware()

	withCid := false
	handle1 := func(ctx echo.Context) error {
		cid := ctxs.GetCidFromContext(ctx.Request().Context())
		withCid = cid != nil && *cid != ""
		return errors.New("some error")
	}
	req := httptest.NewRequest(http.MethodPost, "/cards", strings.NewReader(`{"org_id":  null"cid": "{% uuid 'v4' %}" }`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	if err := cidMiddleware(handle1)(echo.New().NewContext(req, rec)); err == nil {
		t.Error("expected error from handle1")
	}

	if !withCid {
		t.Error("expected context from request with cid")
	}

	if err := serviceImpl.ServiceManager().Close(); err != nil {
		t.Error(err)
	}
}
