package logger

import (
	"context"
	"testing"
	"time"

	"github.com/dalmarcogd/mobstore/products/internal/infra/ctxs"
	"github.com/dalmarcogd/mobstore/products/internal/infra/environment"
	"github.com/dalmarcogd/mobstore/products/internal/services"
)

func TestLoggerSuccessfully(t *testing.T) {
	ctx := ctxs.ContextWithCid(context.Background(), "some-id")

	l := New().WithDefaultFields(func(ctx context.Context, fields *map[string]interface{}) {})
	sm := services.New().WithLogger(l)
	if err := sm.Init(); err != nil {
		t.Error(err)
	}

	l.Info(ctx, "just an info log")
	l.Info(ctx, "just an info log", map[string]interface{}{
		"f": time.Now().String(),
	})
	l.Warn(ctx, "just a warn log")
	l.WithCallerCount(1).Warn(ctx, "just a warn log")
	l.Error(ctx, "just an error log")
	l.WithCallerCount(1).Error(ctx, "just an error log")

	if err := l.ServiceManager().Close(); err != nil {
		t.Error(err)
	}
}

func TestLoggerSuccessfully2(t *testing.T) {
	ctx := ctxs.ContextWithCid(context.Background(), "some-id")

	l := New()
	sm := services.New().WithEnvironment(environment.New()).WithLogger(l)
	if err := sm.Init(); err != nil {
		t.Error(err)
	}

	l.Info(ctx, "just an info log")
	l.Info(ctx, "just an info log", map[string]interface{}{
		"f": time.Now().String(),
	})
	l.Warn(ctx, "just a warn log")
	l.WithCallerCount(1).Warn(ctx, "just a warn log")
	l.Error(ctx, "just an error log")
	l.WithCallerCount(1).Error(ctx, "just an error log")

	if err := l.ServiceManager().Close(); err != nil {
		t.Error(err)
	}
}
