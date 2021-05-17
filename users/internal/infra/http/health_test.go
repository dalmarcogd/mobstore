package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/dalmarcogd/mobstore/users/internal/infra/validator"
	"github.com/dalmarcogd/mobstore/users/internal/services"
)

func TestHttpService_handleReadiness(t *testing.T) {
	serviceImpl := New()
	sm := services.New().
		WithValidator(validator.New()).
		WithHttpServer(serviceImpl)

	if err := sm.Init(); err != nil {
		t.Error("unexpected address")
	}

	req := httptest.NewRequest(http.MethodPost, "/readiness", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("x-cid", "some-id")
	rec := httptest.NewRecorder()

	if err := serviceImpl.handleReadiness(echo.New().NewContext(req, rec)); err != nil {
		t.Error(err)
	}

	if err := serviceImpl.ServiceManager().Close(); err != nil {
		t.Error(err)
	}
}

func TestHttpService_handleLiveness(t *testing.T) {
	serviceImpl := New()
	sm := services.New().
		WithValidator(validator.New()).
		WithHttpServer(serviceImpl)

	if err := sm.Init(); err != nil {
		t.Error("unexpected address")
	}

	req := httptest.NewRequest(http.MethodPost, "/liveness", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("x-cid", "some-id")
	rec := httptest.NewRecorder()

	if err := serviceImpl.handleLiveness(echo.New().NewContext(req, rec)); err != nil {
		t.Error(err)
	}

	if err := serviceImpl.ServiceManager().Close(); err != nil {
		t.Error(err)
	}
}
