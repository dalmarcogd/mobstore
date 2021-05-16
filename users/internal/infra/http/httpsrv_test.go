package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/dalmarcogd/mobstore/users/internal/domains"
	"github.com/dalmarcogd/mobstore/users/internal/infra/ptrs"
	"github.com/dalmarcogd/mobstore/users/internal/infra/validator"
	"github.com/dalmarcogd/mobstore/users/internal/services"
)

type (
	mockHandlerSuccessfully struct {
		services.NoopProductsHandler
		created bool
	}
)

func (m *mockHandlerSuccessfully) Create(_ context.Context, _ domains.UserCreate) (*domains.User, error) {
	m.created = true
	return &domains.User{
		Id:           ptrs.String(uuid.New().String()),
		PriceInCents: ptrs.Int64(10),
		Title:        ptrs.String("my title"),
		Description:  ptrs.String("my description"),
	}, nil
}

func (m *mockHandlerSuccessfully) WithServiceManager(_ services.ServiceManager) services.UsersHandler {
	return m
}

func TestResponseHttpServerSuccessfully(t *testing.T) {
	serviceImpl := New()
	serviceImpl.WithAddress("0.0.0.0:8088")
	sm := services.New().WithUsersHandler(&mockHandlerSuccessfully{}).WithValidator(validator.New()).WithHttpServer(serviceImpl)

	if err := sm.Init(); err != nil {
		t.Error("unexpected address")
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"title": "my-product","description": "my beautiful product","price_in_cents": 10 }`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := serviceImpl.handleCreateV1(c); err != nil {
		t.Error(err)
	}

	if err := serviceImpl.ServiceManager().Close(); err != nil {
		t.Error(err)
	}
}

func TestResponseHttpServerWithError(t *testing.T) {
	serviceImpl := New()
	serviceImpl.WithAddress("0.0.0.0:8088")
	sm := services.New().WithUsersHandler(&mockHandlerSuccessfully{}).WithValidator(validator.New()).WithHttpServer(serviceImpl)

	if err := sm.Init(); err != nil {
		t.Error("unexpected address")
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"description": "my beautiful product","price_in_cents": 10 }`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := serviceImpl.handleCreateV1(c); err == nil {
		t.Error("unexpected response from handle, expected error")
	}

	closed := make(chan int, 1)
	go func() {
		if err := serviceImpl.Run(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				t.Error(err)
			}
			closed <- 1
		}
	}()

	if err := serviceImpl.ServiceManager().Close(); err != nil {
		t.Error(err)
	}
	<-closed
}
