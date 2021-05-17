package userhandler

import (
	"testing"

	"github.com/dalmarcogd/mobstore/users/internal/infra/validator"
	"github.com/dalmarcogd/mobstore/users/internal/services"
)

func TestHandleSuccess(t *testing.T) {
	serviceImpl := New()
	sm := services.New().WithValidator(validator.New()).WithUsersHandler(serviceImpl)

	if err := sm.Init(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := sm.UsersHandler().ServiceManager().UsersHandler().Close(); err != nil {
		t.Error(err)
	}
}
