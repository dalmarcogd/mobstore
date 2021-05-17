package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *httpService) handleReadiness(c echo.Context) error {
	ctx := c.Request().Context()
	if err := s.ServiceManager().Readiness(ctx); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("readiness, unhealthy by: %v", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"services": "OK"})
}

func (s *httpService) handleLiveness(c echo.Context) error {
	ctx := c.Request().Context()
	if err := s.ServiceManager().Liveness(ctx); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("readiness, unhealthy by: %v", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"services": "OK"})
}
