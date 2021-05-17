package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/openzipkin/zipkin-go"

	"github.com/dalmarcogd/mobstore/products/internal/domains"
	"github.com/dalmarcogd/mobstore/products/internal/infra/ctxs"
)

type (
	Generic interface {
		ServiceManager() ServiceManager
		Init(ctx context.Context) error
		Close() error
		Readiness(ctx context.Context) error
		Liveness(ctx context.Context) error
	}
	Database interface {
		Generic
		WithServiceManager(c ServiceManager) Database
		WithProductDatabase() Database
		WithMasterClient(*sql.DB) Database
		WithReplicaClient(*sql.DB) Database
		OpenTransactionMaster(ctx context.Context) (context.Context, error)
		TransactionMaster(ctx context.Context, f func(tx DatabaseTransaction) error) error
		OpenTransactionReplica(ctx context.Context) (context.Context, error)
		TransactionReplica(ctx context.Context, f func(tx DatabaseTransaction) error) error
		CloseTransaction(ctx context.Context, err error) error
	}
	DatabaseTransaction interface {
		Insert(query string, args ...interface{}) (sql.Result, error)
		Update(query string, args ...interface{}) (sql.Result, error)
		Get(query string, args ...interface{}) (*sql.Rows, error)
		Migrate(script string) (sql.Result, error)
	}
	Events interface {
		Generic
		WithServiceManager(c ServiceManager) Events
		Send(ctx context.Context, req domains.EventRequest) error
		SendBulk(ctx context.Context, req []domains.EventRequest) error
	}
	Logger interface {
		Generic
		WithServiceManager(c ServiceManager) Logger
		WithCallerCount(c int) Logger
		Info(ctx context.Context, message string, fields ...map[string]interface{})
		Warn(ctx context.Context, message string, fields ...map[string]interface{})
		Error(ctx context.Context, message string, fields ...map[string]interface{})
		Fatal(ctx context.Context, message string, fields ...map[string]interface{})
	}
	HttpServer interface {
		Generic
		WithServiceManager(c ServiceManager) HttpServer
		Run() error
	}
	Environment interface {
		Generic
		WithServiceManager(c ServiceManager) Environment
		Env() string
		Service() string
		Version() string
		DebugPprof() bool
		AwsRegion() string
		AwsEndpoint() string
		AwsAccessKey() string
		AwsSecretKey() string
		DiscountsAddress() string
		SpanUrl() string
		ProductDatabaseDsn() string
		ProductReplicaDatabaseDsn() string
		ProductTopicAddress() string
		Configs() map[string]interface{}
	}
	Discounts interface {
		Generic
		WithServiceManager(c ServiceManager) Discounts
		Get(ctx context.Context, req domains.DiscountRequest) (*domains.DiscountResponse, error)
	}
	Validator interface {
		Generic
		WithServiceManager(c ServiceManager) Validator
		Validate(ctx context.Context, obj interface{}) error
		ValidateSlice(ctx context.Context, objs interface{}) error
	}
	Spans interface {
		Generic
		WithServiceManager(c ServiceManager) Spans
		New(ctx context.Context, spanConfigs ...domains.SpanConfig) (context.Context, *domains.Span)
		Tracer() *zipkin.Tracer
	}
	ProductsRepository interface {
		Generic
		WithServiceManager(ServiceManager) ProductsRepository
		Search(ctx context.Context, search domains.ProductSearch) ([]domains.Product, error)
		Create(ctx context.Context, product *domains.Product) error
		Update(ctx context.Context, product *domains.Product) error
		Delete(ctx context.Context, product *domains.Product) error
		Migrate(ctx context.Context) error
	}
	ProductsHandler interface {
		Generic
		WithServiceManager(c ServiceManager) ProductsHandler
		Create(ctx context.Context, c domains.ProductCreate) (*domains.Product, error)
		Update(ctx context.Context, c domains.ProductUpdate) (*domains.Product, error)
		Get(ctx context.Context, productId string) (*domains.Product, error)
		List(ctx context.Context) ([]domains.Product, error)
		Delete(ctx context.Context, productId string) (*domains.Product, error)
	}

	ServiceManager interface {
		Init() error
		Close() error
		Readiness(ctx context.Context) error
		Liveness(ctx context.Context) error
		Context() context.Context
		WithLogger(d Logger) ServiceManager
		Logger() Logger
		WithHttpServer(d HttpServer) ServiceManager
		HttpServer() HttpServer
		WithEnvironment(d Environment) ServiceManager
		Environment() Environment
		WithDiscounts(d Discounts) ServiceManager
		Discounts() Discounts
		WithValidator(d Validator) ServiceManager
		Validator() Validator
		WithSpans(d Spans) ServiceManager
		Spans() Spans
		WithProductDatabase(d Database) ServiceManager
		ProductDatabase() Database
		WithProductsRepository(d ProductsRepository) ServiceManager
		ProductsRepository() ProductsRepository
		WithProductsHandler(d ProductsHandler) ServiceManager
		ProductsHandler() ProductsHandler
		WithProductsEvents(d Events) ServiceManager
		ProductsEvents() Events
	}

	serviceManager struct {
		ctx                context.Context
		cancel             context.CancelFunc
		log                Logger
		httpServer         HttpServer
		environment        Environment
		discounts          Discounts
		validator          Validator
		spans              Spans
		productDatabase    Database
		productsRepository ProductsRepository
		productsHandler    ProductsHandler
		productsEvents     Events
	}
)

func New() *serviceManager {
	ctx, cancel := context.WithCancel(ctxs.ContextWithCid(context.Background(), "main"))
	return &serviceManager{
		ctx:                ctx,
		cancel:             cancel,
		log:                NewNoopLogger(),
		httpServer:         NewNoopHttpServer(),
		environment:        NewNoopEnvironment(),
		discounts:          NewNoopDiscounts(),
		validator:          NewNoopValidator(),
		spans:              NewNoopSpans(),
		productDatabase:    NewNoopDatabase(),
		productsRepository: NewNoopProductsRepository(),
		productsHandler:    NewNoopCardsHandler(),
		productsEvents:     NewNoopEvents(),
	}
}

func (s *serviceManager) Init() error {
	now := time.Now()
	if err := s.Environment().Init(s.ctx); err != nil {
		return err
	}
	if err := s.Logger().Init(s.ctx); err != nil {
		return err
	}
	if err := s.Spans().Init(s.ctx); err != nil {
		return err
	}
	if err := s.Validator().Init(s.ctx); err != nil {
		return err
	}
	if err := s.Discounts().Init(s.ctx); err != nil {
		return err
	}
	if err := s.HttpServer().Init(s.ctx); err != nil {
		return err
	}
	if err := s.ProductDatabase().Init(s.ctx); err != nil {
		return err
	}
	if err := s.ProductsRepository().Init(s.ctx); err != nil {
		return err
	}
	if err := s.ProductsHandler().Init(s.ctx); err != nil {
		return err
	}
	if err := s.ProductsEvents().Init(s.ctx); err != nil {
		return err
	}

	s.Logger().Info(
		s.ctx,
		fmt.Sprintf("All services initialized in %v", time.Since(now).String()),
		map[string]interface{}{"envs": s.Environment().Configs()},
	)

	return nil
}

func (s *serviceManager) Close() error {
	now := time.Now()
	var err error
	if errC := s.HttpServer().Close(); errC != nil {
		err = fmt.Errorf("%v - %v", err, errC)
	}
	if errC := s.Environment().Close(); errC != nil {
		err = fmt.Errorf("%v - %v", err, errC)
	}
	if errC := s.Logger().Close(); errC != nil {
		err = fmt.Errorf("%v - %v", err, errC)
	}
	if errC := s.Discounts().Close(); errC != nil {
		err = fmt.Errorf("%v - %v", err, errC)
	}
	if errC := s.Validator().Close(); errC != nil {
		err = fmt.Errorf("%v - %v", err, errC)
	}
	if errC := s.Spans().Close(); errC != nil {
		err = fmt.Errorf("%v - %v", err, errC)
	}
	if errC := s.ProductDatabase().Close(); errC != nil {
		err = fmt.Errorf("%v - %v", err, errC)
	}
	if errC := s.ProductsRepository().Close(); errC != nil {
		err = fmt.Errorf("%v - %v", err, errC)
	}
	if errC := s.ProductsHandler().Close(); errC != nil {
		err = fmt.Errorf("%v - %v", err, errC)
	}
	if errC := s.ProductsEvents().Close(); errC != nil {
		err = fmt.Errorf("%v - %v", err, errC)
	}
	s.cancel()
	time.Sleep(2 * time.Second)
	s.Logger().Info(s.ctx, fmt.Sprintf("All services closed in %v with %v", time.Since(now).String(), err))
	return err
}

func (s *serviceManager) Readiness(ctx context.Context) error {
	if err := s.Environment().Readiness(ctx); err != nil {
		return err
	}
	if err := s.Logger().Readiness(ctx); err != nil {
		return err
	}
	if err := s.Spans().Readiness(ctx); err != nil {
		return err
	}
	if err := s.Validator().Readiness(ctx); err != nil {
		return err
	}
	if err := s.Discounts().Readiness(ctx); err != nil {
		return err
	}
	if err := s.HttpServer().Readiness(ctx); err != nil {
		return err
	}
	if err := s.ProductDatabase().Readiness(ctx); err != nil {
		return err
	}
	if err := s.ProductsRepository().Readiness(ctx); err != nil {
		return err
	}
	if err := s.ProductsHandler().Readiness(ctx); err != nil {
		return err
	}
	if err := s.ProductsEvents().Readiness(ctx); err != nil {
		return err
	}
	return nil
}

func (s *serviceManager) Liveness(ctx context.Context) error {
	if err := s.Environment().Liveness(ctx); err != nil {
		return err
	}
	if err := s.Logger().Liveness(ctx); err != nil {
		return err
	}
	if err := s.Spans().Liveness(ctx); err != nil {
		return err
	}
	if err := s.Validator().Liveness(ctx); err != nil {
		return err
	}
	if err := s.Discounts().Liveness(ctx); err != nil {
		return err
	}
	if err := s.HttpServer().Liveness(ctx); err != nil {
		return err
	}
	if err := s.ProductDatabase().Liveness(ctx); err != nil {
		return err
	}
	if err := s.ProductsRepository().Liveness(ctx); err != nil {
		return err
	}
	if err := s.ProductsHandler().Liveness(ctx); err != nil {
		return err
	}
	if err := s.ProductsEvents().Liveness(ctx); err != nil {
		return err
	}
	return nil
}

func (s *serviceManager) Context() context.Context {
	return s.ctx
}

func (s *serviceManager) WithLogger(d Logger) ServiceManager {
	s.log = d.WithServiceManager(s)
	return s
}

func (s *serviceManager) Logger() Logger {
	return s.log
}

func (s *serviceManager) WithHttpServer(d HttpServer) ServiceManager {
	s.httpServer = d.WithServiceManager(s)
	return s
}

func (s *serviceManager) HttpServer() HttpServer {
	return s.httpServer
}

func (s *serviceManager) WithEnvironment(d Environment) ServiceManager {
	s.environment = d.WithServiceManager(s)
	return s
}

func (s *serviceManager) Environment() Environment {
	return s.environment
}

func (s *serviceManager) WithDiscounts(d Discounts) ServiceManager {
	s.discounts = d.WithServiceManager(s)
	return s
}

func (s *serviceManager) Discounts() Discounts {
	return s.discounts
}

func (s *serviceManager) WithValidator(d Validator) ServiceManager {
	s.validator = d.WithServiceManager(s)
	return s
}

func (s *serviceManager) Validator() Validator {
	return s.validator
}

func (s *serviceManager) WithSpans(d Spans) ServiceManager {
	s.spans = d.WithServiceManager(s)
	return s
}

func (s *serviceManager) Spans() Spans {
	return s.spans
}

func (s *serviceManager) WithProductDatabase(d Database) ServiceManager {
	s.productDatabase = d.WithServiceManager(s)
	return s
}

func (s *serviceManager) ProductDatabase() Database {
	return s.productDatabase
}

func (s *serviceManager) WithProductsRepository(d ProductsRepository) ServiceManager {
	s.productsRepository = d.WithServiceManager(s)
	return s
}

func (s *serviceManager) ProductsRepository() ProductsRepository {
	return s.productsRepository
}

func (s *serviceManager) WithProductsHandler(d ProductsHandler) ServiceManager {
	s.productsHandler = d.WithServiceManager(s)
	return s
}

func (s *serviceManager) ProductsHandler() ProductsHandler {
	return s.productsHandler
}

func (s *serviceManager) WithProductsEvents(d Events) ServiceManager {
	s.productsEvents = d.WithServiceManager(s)
	return s
}

func (s *serviceManager) ProductsEvents() Events {
	return s.productsEvents
}
