package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/openzipkin/zipkin-go"

	"github.com/dalmarcogd/mobstore/users/internal/domains"
	"github.com/dalmarcogd/mobstore/users/internal/infra/ctxs"
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
		WithUserDatabase() Database
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
		SpanUrl() string
		UserDatabaseDsn() string
		UserReplicaDatabaseDsn() string
		UserTopicAddress() string
		Configs() map[string]interface{}
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
	UsersRepository interface {
		Generic
		WithServiceManager(ServiceManager) UsersRepository
		Search(ctx context.Context, search domains.UserSearch) ([]domains.User, error)
		Create(ctx context.Context, product *domains.User) error
		Update(ctx context.Context, product *domains.User) error
		Delete(ctx context.Context, product *domains.User) error
		Migrate(ctx context.Context) error
	}
	UsersHandler interface {
		Generic
		WithServiceManager(c ServiceManager) UsersHandler
		Create(ctx context.Context, c domains.UserCreate) (*domains.User, error)
		Update(ctx context.Context, c domains.UserUpdate) (*domains.User, error)
		Get(ctx context.Context, productId string) (*domains.User, error)
		List(ctx context.Context) ([]domains.User, error)
		Delete(ctx context.Context, productId string) (*domains.User, error)
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
		WithValidator(d Validator) ServiceManager
		Validator() Validator
		WithSpans(d Spans) ServiceManager
		Spans() Spans
		WithUserDatabase(d Database) ServiceManager
		UserDatabase() Database
		WithUsersRepository(d UsersRepository) ServiceManager
		UsersRepository() UsersRepository
		WithUsersHandler(d UsersHandler) ServiceManager
		UsersHandler() UsersHandler
		WithUsersEvents(d Events) ServiceManager
		UsersEvents() Events
	}

	serviceManager struct {
		ctx             context.Context
		cancel          context.CancelFunc
		log             Logger
		httpServer      HttpServer
		environment     Environment
		validator       Validator
		spans           Spans
		userDatabase    Database
		usersRepository UsersRepository
		usersHandler    UsersHandler
		usersEvents     Events
	}
)

func New() *serviceManager {
	ctx, cancel := context.WithCancel(ctxs.ContextWithCid(context.Background(), "main"))
	return &serviceManager{
		ctx:             ctx,
		cancel:          cancel,
		log:             NewNoopLogger(),
		httpServer:      NewNoopHttpServer(),
		environment:     NewNoopEnvironment(),
		validator:       NewNoopValidator(),
		spans:           NewNoopSpans(),
		userDatabase:    NewNoopDatabase(),
		usersRepository: NewNoopUsersRepository(),
		usersHandler:    NewNoopUserHandler(),
		usersEvents:     NewNoopEvents(),
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
	if err := s.HttpServer().Init(s.ctx); err != nil {
		return err
	}
	if err := s.UserDatabase().Init(s.ctx); err != nil {
		return err
	}
	if err := s.UsersRepository().Init(s.ctx); err != nil {
		return err
	}
	if err := s.UsersHandler().Init(s.ctx); err != nil {
		return err
	}
	if err := s.UsersEvents().Init(s.ctx); err != nil {
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
	if errC := s.Validator().Close(); errC != nil {
		err = fmt.Errorf("%v - %v", err, errC)
	}
	if errC := s.Spans().Close(); errC != nil {
		err = fmt.Errorf("%v - %v", err, errC)
	}
	if errC := s.UserDatabase().Close(); errC != nil {
		err = fmt.Errorf("%v - %v", err, errC)
	}
	if errC := s.UsersRepository().Close(); errC != nil {
		err = fmt.Errorf("%v - %v", err, errC)
	}
	if errC := s.UsersHandler().Close(); errC != nil {
		err = fmt.Errorf("%v - %v", err, errC)
	}
	if errC := s.UsersEvents().Close(); errC != nil {
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
	if err := s.HttpServer().Readiness(ctx); err != nil {
		return err
	}
	if err := s.UserDatabase().Readiness(ctx); err != nil {
		return err
	}
	if err := s.UsersRepository().Readiness(ctx); err != nil {
		return err
	}
	if err := s.UsersHandler().Readiness(ctx); err != nil {
		return err
	}
	if err := s.UsersEvents().Readiness(ctx); err != nil {
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
	if err := s.HttpServer().Liveness(ctx); err != nil {
		return err
	}
	if err := s.UserDatabase().Liveness(ctx); err != nil {
		return err
	}
	if err := s.UsersRepository().Liveness(ctx); err != nil {
		return err
	}
	if err := s.UsersHandler().Liveness(ctx); err != nil {
		return err
	}
	if err := s.UsersEvents().Liveness(ctx); err != nil {
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

func (s *serviceManager) WithUserDatabase(d Database) ServiceManager {
	s.userDatabase = d.WithServiceManager(s)
	return s
}

func (s *serviceManager) UserDatabase() Database {
	return s.userDatabase
}

func (s *serviceManager) WithUsersRepository(d UsersRepository) ServiceManager {
	s.usersRepository = d.WithServiceManager(s)
	return s
}

func (s *serviceManager) UsersRepository() UsersRepository {
	return s.usersRepository
}

func (s *serviceManager) WithUsersHandler(d UsersHandler) ServiceManager {
	s.usersHandler = d.WithServiceManager(s)
	return s
}

func (s *serviceManager) UsersHandler() UsersHandler {
	return s.usersHandler
}

func (s *serviceManager) WithUsersEvents(d Events) ServiceManager {
	s.usersEvents = d.WithServiceManager(s)
	return s
}

func (s *serviceManager) UsersEvents() Events {
	return s.usersEvents
}
