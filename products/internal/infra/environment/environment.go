package environment

import (
	"context"
	"fmt"

	"github.com/dalmarcogd/mobstore/products/internal/services"

	"github.com/gosidekick/goconfig"
)

//Env this object keep the all variables environment
type (
	environment struct {
		Environment            string `cfg:"ENVIRONMENT" cfgDefault:"local" cfgRequired:"true"`
		Service                string `cfg:"SYSTEM" cfgDefault:"api-products" cfgRequired:"true"`
		Version                string `cfg:"SYSTEM_VERSION" cfgDefault:"v1.0.0" cfgRequired:"true"`
		DebugPprof             bool   `cfg:"DEBUG_PPROF" cfgDefault:"true" `
		AwsRegion              string `cfg:"AWS_REGION" cfgDefault:"sa-east-1" `
		AwsEndpoint            string `cfg:"AWS_ENDPOINT" cfgDefault:"http://localhost:4566" `
		AwsAccessKey           string `cfg:"AWS_ACCESS_KEY" cfgDefault:"fake_access_key" `
		AwsSecretKey           string `cfg:"AWS_SECRET_KEY" cfgDefault:"fake_secret_key" `
		DiscountsAddress       string `cfg:"DISCOUNTS_API_URL" cfgDefault:"http://localhost:3466"`
		ZipkinUrl              string `cfg:"ZIPKIN_URL_V2" cfgDefault:"http://localhost:9411/api/v2/spans" cfgRequired:"true"`
		ProductTopicAddress    string `cfg:"PRODUCT_TOPIC_ADDRESS" cfgDefault:"arn:aws:sns:sa-east-1:000000000000:ProductsEvents.fifo" cfgRequired:"true"`
		ProductMysqlUser       string `cfg:"PRODUCT_MYSQL_USER" cfgDefault:"products" `
		ProductMysqlPassword   string `cfg:"PRODUCT_MYSQL_PASS" cfgDefault:"my-password" `
		ProductMysqlDatabase   string `cfg:"PRODUCT_MYSQL_DATABASE" cfgDefault:"products" `
		ProductMysqlEndpoint   string `cfg:"PRODUCT_MYSQL_ENDPOINT" cfgDefault:"localhost:3306" `
		ProductRRMysqlEndpoint string `cfg:"PRODUCT_RR_MYSQL_ENDPOINT" cfgDefault:"localhost:3306" `
		RedisAddress           string `cfg:"REDIS_ADDRESS" cfgDefault:"localhost:6379"`
		RedisReplicaAddress    string `cfg:"REDIS_REPLICA_ADDRESS" cfgDefault:"localhost:6379"`
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

func (s *envService) DiscountsAddress() string {
	return s.environment.DiscountsAddress
}

func (s *envService) SpanUrl() string {
	return s.environment.ZipkinUrl
}

func (s *envService) ProductTopicAddress() string {
	return s.environment.ProductTopicAddress
}

func (s *envService) ProductReplicaDatabaseDsn() string {
	return fmt.Sprintf(
		"%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=UTC&tls=skip-verify",
		s.environment.ProductMysqlUser,
		s.environment.ProductMysqlPassword,
		s.environment.ProductRRMysqlEndpoint,
		s.environment.ProductMysqlDatabase,
	)
}

func (s *envService) ProductDatabaseDsn() string {
	return fmt.Sprintf(
		"%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=UTC&tls=skip-verify",
		s.environment.ProductMysqlUser,
		s.environment.ProductMysqlPassword,
		s.environment.ProductMysqlEndpoint,
		s.environment.ProductMysqlDatabase,
	)
}

func (s *envService) RedisAddress() string {
	return s.environment.RedisAddress
}

func (s *envService) RedisReplicaAddress() string {
	return s.environment.RedisReplicaAddress
}

func (s *envService) Configs() map[string]interface{} {
	return map[string]interface{}{
		"ENVIRONMENT":               s.environment.Environment,
		"SYSTEM":                    s.environment.Service,
		"SYSTEM_VERSION":            s.environment.Version,
		"DEBUG_PPROF":               s.environment.DebugPprof,
		"AWS_REGION":                s.environment.AwsRegion,
		"AWS_ENDPOINT":              s.environment.AwsEndpoint,
		"PRODUCT_TOPIC_ADDRESS":     s.environment.ProductTopicAddress,
		"ZIPKIN_URL_V2":             s.environment.ZipkinUrl,
		"DISCOUNTS_API_URL":         s.environment.DiscountsAddress,
		"PRODUCT_MYSQL_USER":        s.environment.ProductMysqlUser,
		"PRODUCT_MYSQL_PASS":        "********",
		"PRODUCT_MYSQL_DATABASE":    s.environment.ProductMysqlDatabase,
		"PRODUCT_MYSQL_ENDPOINT":    s.environment.ProductMysqlEndpoint,
		"PRODUCT_RR_MYSQL_ENDPOINT": s.environment.ProductRRMysqlEndpoint,
		"REDIS_ADDRESS":             s.environment.RedisAddress,
		"REDIS_REPLICA_ADDRESS":     s.environment.RedisReplicaAddress,
	}
}
