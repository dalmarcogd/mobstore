package validator

import (
	"context"
	"testing"

	"github.com/dalmarcogd/mobstore/products/internal/services"
)

func TestValidatorSuccess(t *testing.T) {
	serviceImpl := New()
	sm := services.New().WithValidator(serviceImpl)
	if err := sm.Init(); err != nil {
		t.Error(err)
	}

	type args struct {
		Name string `validate:"required"`
		Cid  string `validate:"required"`
	}

	ctx := context.Background()
	if err := serviceImpl.Validate(ctx, args{
		Name: "asdmaosidjs",
		Cid:  "asdoiasod",
	}); err != nil {
		t.Error(err)
	}

	if err := sm.Close(); err != nil {
		t.Error(err)
	}
}

func TestValidatorError(t *testing.T) {
	serviceImpl := New()
	sm := services.New().WithValidator(serviceImpl)
	if err := sm.Init(); err != nil {
		t.Error(err)
	}

	type args struct {
		Name string `validate:"required"`
		Cid  string `validate:"required"`
	}

	ctx := context.Background()
	if err := serviceImpl.Validate(ctx, args{}); err == nil {
		t.Error("expected error from validator")
	}

	if err := serviceImpl.ServiceManager().Close(); err != nil {
		t.Error(err)
	}
}
