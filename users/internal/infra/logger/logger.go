package logger

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"

	"github.com/dalmarcogd/mobstore/users/internal/infra/ctxs"
	"github.com/dalmarcogd/mobstore/users/internal/services"
)

type (
	DefaultFields func(ctx context.Context, fields *map[string]interface{})
	loggerService struct {
		services.NoopHealth
		serviceManager services.ServiceManager
		ctx            context.Context
		logger         *logrus.Logger
		defaultFields  DefaultFields
		callerCount    int
	}
)

func New() *loggerService {
	return &loggerService{
		callerCount: 2,
	}
}

func (s *loggerService) WithDefaultFields(f DefaultFields) *loggerService {
	s.defaultFields = f
	return s
}

func (s *loggerService) WithCallerCount(i int) services.Logger {
	sC := *s
	sC.callerCount = i
	return &sC
}

func (s *loggerService) Init(ctx context.Context) error {
	s.ctx = ctx
	s.logger = logrus.New()
	s.logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "time",
			logrus.FieldKeyMsg:  "text",
		},
	})
	s.logger.SetOutput(os.Stdout)
	if s.defaultFields == nil {
		s.WithDefaultFields(func(ctx context.Context, fields *map[string]interface{}) {
			f := *fields
			f["cid"] = ctxs.GetCidFromContext(ctx)
			f["service"] = s.ServiceManager().Environment().Service()
			f["version"] = s.ServiceManager().Environment().Version()
			fields = &f
		})
	}

	return nil
}

func (s *loggerService) Close() error {
	return nil
}

func (s *loggerService) WithServiceManager(c services.ServiceManager) services.Logger {
	s.serviceManager = c
	return s
}

func (s *loggerService) ServiceManager() services.ServiceManager {
	return s.serviceManager
}

func (s *loggerService) output(ctx context.Context, fields ...map[string]interface{}) *logrus.Entry {
	f := make(map[string]interface{}, 0)

	if s.defaultFields != nil {
		s.defaultFields(ctx, &f)
	}

	for _, ff := range fields {
		for k, v := range ff {
			f[k] = v
		}
	}

	stack := ""
	pc, _, line, ok := runtime.Caller(s.callerCount)
	if ok {
		stack = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), line)
	}
	f["caller"] = stack

	return s.logger.WithFields(f)
}

func (s *loggerService) Info(ctx context.Context, message string, fields ...map[string]interface{}) {
	s.output(ctx, fields...).Info(message)
}

func (s *loggerService) Warn(ctx context.Context, message string, fields ...map[string]interface{}) {
	s.output(ctx, fields...).Warn(message)
}

func (s *loggerService) Error(ctx context.Context, message string, fields ...map[string]interface{}) {
	s.output(ctx, fields...).Error(message)
}

func (s *loggerService) Fatal(ctx context.Context, message string, fields ...map[string]interface{}) {
	s.output(ctx, fields...).Fatal(message)
}
