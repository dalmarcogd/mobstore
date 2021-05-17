package userhandler

import (
	"context"

	"github.com/dalmarcogd/mobstore/users/internal/services"
)

type (
	productsHandler struct {
		services.NoopHealth
		serviceManager services.ServiceManager
		ctx            context.Context
	}
)

func New() *productsHandler {
	return &productsHandler{}
}

func (s *productsHandler) Init(ctx context.Context) error {
	s.ctx = ctx
	return nil
}

func (s *productsHandler) Close() error {
	return nil
}

func (s *productsHandler) WithServiceManager(c services.ServiceManager) services.UsersHandler {
	s.serviceManager = c
	return s
}

func (s *productsHandler) ServiceManager() services.ServiceManager {
	return s.serviceManager
}
