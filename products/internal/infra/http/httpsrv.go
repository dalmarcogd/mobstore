package http

import (
	"context"
	"io"
	"net/http/pprof"

	"github.com/labstack/echo/v4"

	"github.com/dalmarcogd/mobstore/products/internal/services"
)

type (
	httpService struct {
		services.NoopHealth
		serviceManager services.ServiceManager
		ctx            context.Context
		echo           *echo.Echo
		address        string
	}
)

func New() *httpService {
	return &httpService{}
}

func (s *httpService) WithAddress(address string) *httpService {
	s.address = address
	return s
}

func (s *httpService) Init(ctx context.Context) error {
	s.ctx = ctx
	if s.echo == nil {
		s.echo = echo.New()
		s.echo.Logger.SetOutput(io.Discard)
		s.echo.Use(
			CidContextMiddleware(),
			SpanMiddleware(s.ServiceManager().Spans(), s.ServiceManager().Environment().Service(), []string{"/readiness", "/liveness"}),
			LogMiddleware(s.ServiceManager().Logger(), []string{"/readiness", "/liveness"}),
		)
	}
	s.RegisterRoutes()
	return nil
}

func (s *httpService) Close() error {
	if err := s.echo.Shutdown(s.ctx); err != nil {
		return err
	}
	return s.echo.Close()
}

func (s *httpService) WithServiceManager(c services.ServiceManager) services.HttpServer {
	s.serviceManager = c
	return s
}

func (s *httpService) ServiceManager() services.ServiceManager {
	return s.serviceManager
}

func (s *httpService) RegisterRoutes() *httpService {
	s.echo.GET("/readiness", s.handleReadiness)
	s.echo.GET("/liveness", s.handleLiveness)

	gV1 := s.echo.Group("/v1")

	gV1.POST("/products", s.handleCreateV1)
	gV1.GET("/products", s.handleGetV1)
	gV1.GET("/products/:productId", s.handleGetByIdV1)
	gV1.PATCH("/products/:productId", s.handleUpdateV1)
	gV1.DELETE("/products/:productId", s.handleDeleteV1)

	//routes of pprof
	if s.ServiceManager().Environment().DebugPprof() {
		gDebug := s.echo.Group("/debug")
		gDebug.GET("/pprof/", func(c echo.Context) error {
			pprof.Index(c.Response().Writer, c.Request())
			return nil
		})
		gDebug.GET("/pprof/allocs", func(c echo.Context) error {
			pprof.Handler("allocs").ServeHTTP(c.Response().Writer, c.Request())
			return nil
		})
		gDebug.GET("/pprof/heap", func(c echo.Context) error {
			pprof.Handler("heap").ServeHTTP(c.Response().Writer, c.Request())
			return nil
		})
		gDebug.GET("/pprof/goroutine", func(c echo.Context) error {
			pprof.Handler("goroutine").ServeHTTP(c.Response().Writer, c.Request())
			return nil
		})
		gDebug.GET("/pprof/block", func(c echo.Context) error {
			pprof.Handler("block").ServeHTTP(c.Response().Writer, c.Request())
			return nil
		})
		gDebug.GET("/pprof/threadcreate", func(c echo.Context) error {
			pprof.Handler("threadcreate").ServeHTTP(c.Response().Writer, c.Request())
			return nil
		})
		gDebug.GET("/pprof/cmdline", func(c echo.Context) error {
			pprof.Cmdline(c.Response().Writer, c.Request())
			return nil
		})
		gDebug.GET("/pprof/profile", func(c echo.Context) error {
			pprof.Profile(c.Response().Writer, c.Request())
			return nil
		})
		gDebug.GET("/pprof/symbol", func(c echo.Context) error {
			pprof.Symbol(c.Response().Writer, c.Request())
			return nil
		})
		gDebug.POST("/pprof/symbol", func(c echo.Context) error {
			pprof.Trace(c.Response().Writer, c.Request())
			return nil
		})
		gDebug.GET("/pprof/trace", func(c echo.Context) error {
			pprof.Handler("block").ServeHTTP(c.Response().Writer, c.Request())
			return nil
		})
		gDebug.GET("/pprof/mutex", func(c echo.Context) error {
			pprof.Handler("mutex").ServeHTTP(c.Response().Writer, c.Request())
			return nil
		})
		s.ServiceManager().Logger().Info(s.ctx, "Registered route of /debug/pprof")
	}

	return s
}

func (s *httpService) Run() error {
	return s.echo.Start(s.address)
}
