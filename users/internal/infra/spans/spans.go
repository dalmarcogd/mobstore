package spans

import (
	"context"
	"runtime"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
	"github.com/openzipkin/zipkin-go/reporter/http"

	"github.com/dalmarcogd/mobstore/users/internal/domains"
	"github.com/dalmarcogd/mobstore/users/internal/infra/ctxs"
	"github.com/dalmarcogd/mobstore/users/internal/services"
)

type (
	spansService struct {
		services.NoopHealth
		serviceManager services.ServiceManager
		ctx            context.Context
		tracer         *zipkin.Tracer
		reporter       reporter.Reporter
		host           string
		serviceName    string
		version        string
	}
)

func New() *spansService {
	return &spansService{}
}

func (s *spansService) Init(ctx context.Context) error {
	s.ctx = ctx
	if s.tracer == nil {
		if s.host == "" {
			s.host = s.ServiceManager().Environment().SpanUrl()
		}
		s.reporter = http.NewReporter(s.host, http.BatchInterval(time.Second*3))

		if s.serviceName == "" {
			s.serviceName = s.ServiceManager().Environment().Service()
		}
		// create our local spansService endpoint
		endpoint, err := zipkin.NewEndpoint(s.serviceName, "0.0.0.0:8080")
		if err != nil {
			return err
		}

		s.tracer, err = zipkin.NewTracer(
			s.reporter,
			zipkin.WithLocalEndpoint(endpoint),
			zipkin.WithTraceID128Bit(true),
		)
		if err != nil {
			return err
		}
	}

	// optionally set as Global OpenTracing tracer instance
	opentracing.SetGlobalTracer(zipkintracer.Wrap(s.tracer))
	return nil
}

func (s *spansService) Close() error {
	if err := s.reporter.Close(); err != nil {
		return err
	}
	return nil
}

func (s *spansService) WithServiceManager(c services.ServiceManager) services.Spans {
	s.serviceManager = c
	return s
}

func (s *spansService) ServiceManager() services.ServiceManager {
	return s.serviceManager
}

func (s *spansService) New(ctx context.Context, spanConfigs ...domains.SpanConfig) (context.Context, *domains.Span) {
	var cid string
	if c := ctxs.GetCidFromContext(ctx); c != nil {
		cid = *c
	}

	funcName := ""
	line := 0
	fileName := ""
	if pc, f, l, ok := runtime.Caller(1); ok {
		funcName = runtime.FuncForPC(pc).Name()
		lastDot := strings.LastIndexByte(funcName, '.')
		if lastDot < 0 {
			lastDot = 0
		}
		funcName = funcName[lastDot+1:]
		fileName = f
		line = l
	}

	sp := &domains.Span{
		Name:         funcName,
		Cid:          cid,
		Line:         line,
		FileName:     fileName,
		FuncName:     funcName,
		Version:      s.version,
		Resource:     "cpu",
		Custom:       map[string]interface{}{},
		InternalSpan: nil,
	}

	for _, config := range spanConfigs {
		config.Apply(ctx, sp)
	}

	sp.InternalSpan, ctx = s.tracer.StartSpanFromContext(ctx, sp.Name)

	return ctx, sp
}

func (s *spansService) Tracer() *zipkin.Tracer {
	return s.tracer
}
