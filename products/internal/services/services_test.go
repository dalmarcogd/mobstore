package services

import (
	"context"
	"testing"
)

func TestInitCloseNoopServices(t *testing.T) {
	impl := New()
	if err := impl.Init(); err != nil {
		t.Error(err)
	}
	if err := impl.Readiness(context.Background()); err != nil {
		t.Error(err)
	}
	if err := impl.Liveness(context.Background()); err != nil {
		t.Error(err)
	}
	if err := impl.Close(); err != nil {
		t.Error(err)
	}
}
