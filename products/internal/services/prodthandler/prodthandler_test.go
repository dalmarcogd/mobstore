package prodthandler

import (
	"testing"

	"github.com/dalmarcogd/mobstore/products/internal/infra/validator"
	"github.com/dalmarcogd/mobstore/products/internal/services"
)

func TestHandleSuccess(t *testing.T) {
	serviceImpl := New()
	sm := services.New().WithValidator(validator.New()).WithProductsHandler(serviceImpl)

	if err := sm.Init(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := sm.ProductsHandler().ServiceManager().ProductsHandler().Close(); err != nil {
		t.Error(err)
	}
}
