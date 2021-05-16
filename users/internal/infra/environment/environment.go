package environment

import (
	"context"
	"fmt"

	"github.com/dalmarcogd/mobstore/users/internal/services"

	"github.com/gosidekick/goconfig"
)

//Env this object keep the all variables environment
type (
	environment struct {
		Environment         string `cfg:"ENVIRONMENT" cfgDefault:"local" cfgRequired:"true"`
		Service             string `cfg:"SYSTEM" cfgDefault:"api-products" cfgRequired:"true"`
		Version             string `cfg:"SYSTEM_VERSION" cfgDefault:"v1.0.0" cfgRequired:"true"`
		DebugPprof          bool   `cfg:"DEBUG_PPROF" cfgDefault:"true" `
		AwsRegion           string `cfg:"AWS_REGION" cfgDefault:"sa-east-1" `
		AwsEndpoint         string `cfg:"AWS_ENDPOINT" cfgDefault:"http://localhost:4566" `
		AwsAccessKey        string `cfg:"AWS_SECRET_KEY" cfgDefault:"fake_access_key" `
		AwsSecretKey        string `cfg:"AWS_ACCESS_KEY" cfgDefault:"fake_secret_key" `
		ZipkinUrl           string `cfg:"ZIPKIN_URL_V2" cfgDefault:"http://localhost:9411/api/v2/spans" cfgRequired:"true"`
		UserTopicAddress    string `cfg:"USER_TOPIC_ADDRESS" cfgDefault:"arn:aws:sns:sa-east-1:000000000000:UsersEvents.fifo" cfgRequired:"true"`
		UserMysqlUser       string `cfg:"USER_MYSQL_USER" cfgDefault:"users" `
		UserMysqlPassword   string `cfg:"USER_MYSQL_PASS" cfgDefault:"my-password" `
		UserMysqlDatabase   string `cfg:"USER_MYSQL_DATABASE" cfgDefault:"users" `
		UserMysqlEndpoint   string `cfg:"USER_MYSQL_ENDPOINT" cfgDefault:"localhost:3306" `
		UserRRMysqlEndpoint string `cfg:"USER_RR_MYSQL_ENDPOINT" cfgDefault:"localhost:3306" `
	}

	envService struct {
		services.NoopHealth
		serviceManager services.ServiceManager
		ctx            context.Context
		environment    *environment
	}
)

func New() *envService {
	return &envService{}
}

func (s *envService) Init(ctx context.Context) error {
	s.ctx = ctx
	s.environment = &environment{}
	if err := goconfig.Parse(s.environment); err != nil {
		return err
	}
	return nil
}

func (s *envService) Close() error {
	return nil
}

func (s *envService) WithServiceManager(c services.ServiceManager) services.Environment {
	s.serviceManager = c
	return s
}

func (s *envService) ServiceManager() services.ServiceManager {
	return s.serviceManager
}

func (s *envService) Env() string {
	return s.environment.Environment
}

func (s *envService) Service() string {
	return s.environment.Service
}

func (s *envService) Version() string {
	return s.environment.Version
}

func (s *envService) DebugPprof() bool {
	return s.environment.DebugPprof
}

func (s *envService) AwsRegion() string {
	return s.environment.AwsRegion
}

func (s *envService) AwsEndpoint() string {
	return s.environment.AwsEndpoint
}

func (s *envService) AwsAccessKey() string {
	return s.environment.AwsAccessKey
}

func (s *envService) AwsSecretKey() string {
	return s.environment.AwsSecretKey
}

func (s *envService) SpanUrl() string {
	return s.environment.ZipkinUrl
}

func (s *envService) UserTopicAddress() string {
	return s.environment.UserTopicAddress
}

func (s *envService) UserReplicaDatabaseDsn() string {
	return fmt.Sprintf(
		"%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local&tls=skip-verify",
		s.environment.UserMysqlUser,
		s.environment.UserMysqlPassword,
		s.environment.UserRRMysqlEndpoint,
		s.environment.UserMysqlDatabase,
	)
}

func (s *envService) UserDatabaseDsn() string {
	return fmt.Sprintf(
		"%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local&tls=skip-verify",
		s.environment.UserMysqlUser,
		s.environment.UserMysqlPassword,
		s.environment.UserMysqlEndpoint,
		s.environment.UserMysqlDatabase,
	)
}

func (s *envService) Configs() map[string]interface{} {
	return map[string]interface{}{
		"ENVIRONMENT":            s.environment.Environment,
		"SYSTEM":                 s.environment.Service,
		"SYSTEM_VERSION":         s.environment.Version,
		"DEBUG_PPROF":            s.environment.DebugPprof,
		"AWS_REGION":             s.environment.AwsRegion,
		"AWS_ENDPOINT":           s.environment.AwsEndpoint,
		"USER_TOPIC_ADDRESS":     s.environment.UserTopicAddress,
		"ZIPKIN_URL_V2":          s.environment.ZipkinUrl,
		"USER_MYSQL_USER":        s.environment.UserMysqlUser,
		"USER_MYSQL_PASS":        "********",
		"USER_MYSQL_DATABASE":    s.environment.UserMysqlDatabase,
		"USER_MYSQL_ENDPOINT":    s.environment.UserMysqlEndpoint,
		"USER_RR_MYSQL_ENDPOINT": s.environment.UserRRMysqlEndpoint,
	}
}
