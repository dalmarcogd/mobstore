package services

import (
	"context"
	"database/sql"

	"github.com/openzipkin/zipkin-go"

	"github.com/dalmarcogd/mobstore/users/internal/domains"
)

type (
	NoopHealth         struct{}
	NoopServiceManager struct{}
	NoopDatabase       struct {
		NoopHealth
		NoopServiceManager
	}
	NoopHttpServer struct {
		NoopHealth
		NoopServiceManager
	}
	NoopLogger struct {
		NoopHealth
		NoopServiceManager
	}
	NoopEnvironment struct {
		NoopHealth
		NoopServiceManager
	}
	NoopValidator struct {
		NoopHealth
		NoopServiceManager
	}
	NoopSpans struct {
		NoopHealth
		NoopServiceManager
	}
	NoopUsersRepository struct {
		NoopHealth
		NoopServiceManager
	}
	NoopUsersHandler struct {
		NoopHealth
		NoopServiceManager
	}
	NoopEvents struct {
		NoopHealth
		NoopServiceManager
	}
)

func (n *NoopServiceManager) ServiceManager() ServiceManager {
	return nil
}

func (n *NoopServiceManager) Init(_ context.Context) error {
	return nil
}

func (n *NoopServiceManager) Close() error {
	return nil
}

func (n *NoopHealth) Readiness(_ context.Context) error {
	return nil
}

func (n *NoopHealth) Liveness(_ context.Context) error {
	return nil
}

func NewNoopDatabase() *NoopDatabase {
	return &NoopDatabase{}
}

func (n *NoopDatabase) WithServiceManager(_ ServiceManager) Database {
	return n
}

func (n *NoopDatabase) WithCardsAutomaticUpdaterDatabase() Database {
	return n
}

func (n *NoopDatabase) WithCardsEmbossingsDatabase() Database {
	return n
}

func (n *NoopDatabase) WithUserDatabase() Database {
	return n
}

func (n *NoopDatabase) WithMasterClient(_ *sql.DB) Database {
	return nil
}

func (n *NoopDatabase) WithReplicaClient(_ *sql.DB) Database {
	return nil
}

func (n *NoopDatabase) Ping(_ context.Context) error {
	return nil
}

func (n *NoopDatabase) OpenTransactionMaster(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (n *NoopDatabase) TransactionMaster(context.Context, func(tx DatabaseTransaction) error) error {
	return nil
}

func (n *NoopDatabase) OpenTransactionReplica(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (n *NoopDatabase) TransactionReplica(context.Context, func(tx DatabaseTransaction) error) error {
	return nil
}

func (n *NoopDatabase) CloseTransaction(_ context.Context, err error) error {
	return err
}

func (n *NoopDatabase) Insert(_ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (n *NoopDatabase) Update(_ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (n *NoopDatabase) Get(_ context.Context, _ string, _ ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func NewNoopHttpServer() *NoopHttpServer {
	return &NoopHttpServer{}
}

func (n *NoopHttpServer) WithServiceManager(_ ServiceManager) HttpServer {
	return n
}

func (n *NoopHttpServer) Run() error {
	return nil
}

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

func (n *NoopLogger) WithServiceManager(_ ServiceManager) Logger {
	return n
}

func (n *NoopLogger) WithCallerCount(_ int) Logger {
	return n
}

func (n *NoopLogger) Info(_ context.Context, _ string, _ ...map[string]interface{}) {}

func (n *NoopLogger) Warn(_ context.Context, _ string, _ ...map[string]interface{}) {}

func (n *NoopLogger) Error(_ context.Context, _ string, _ ...map[string]interface{}) {}

func (n *NoopLogger) Fatal(_ context.Context, _ string, _ ...map[string]interface{}) {}

func NewNoopEnvironment() *NoopEnvironment {
	return &NoopEnvironment{}
}

func (n *NoopEnvironment) WithServiceManager(_ ServiceManager) Environment {
	return n
}

func (n *NoopEnvironment) Env() string {
	return ""
}

func (n *NoopEnvironment) Service() string {
	return ""
}

func (n *NoopEnvironment) Version() string {
	return ""
}

func (n *NoopEnvironment) DebugPprof() bool {
	return false
}

func (n *NoopEnvironment) LockerConfigDynamoTable() string {
	return ""
}

func (n *NoopEnvironment) UserReplicaDatabaseDsn() string {
	return ""
}

func (n *NoopEnvironment) UserDatabaseDsn() string {
	return ""
}

func (n *NoopEnvironment) AwsRegion() string {
	return ""
}

func (n *NoopEnvironment) AwsEndpoint() string {
	return ""
}

func (n *NoopEnvironment) AwsAccessKey() string {
	return ""
}

func (n *NoopEnvironment) AwsSecretKey() string {
	return ""
}

func (n *NoopEnvironment) StatsdServer() string {
	return ""
}

func (n *NoopEnvironment) DiscountsAddress() string {
	return ""
}

func (n *NoopEnvironment) OrgsAddress() string {
	return ""
}

func (n *NoopEnvironment) OrgsV2Address() string {
	return ""
}

func (n *NoopEnvironment) SpanUrl() string {
	return ""
}

func (n *NoopEnvironment) UserTopicAddress() string {
	return ""
}

func (n *NoopEnvironment) Timeline() string {
	return ""
}

func (n *NoopEnvironment) SecretManagerId() string {
	return ""
}

func (n *NoopEnvironment) RedisAddress() string {
	return ""
}

func (n *NoopEnvironment) RedisReplicaAddress() string {
	return ""
}

func (n *NoopEnvironment) Configs() map[string]interface{} {
	return map[string]interface{}{}
}

func NewNoopValidator() *NoopValidator {
	return &NoopValidator{}
}

func (n *NoopValidator) WithServiceManager(_ ServiceManager) Validator {
	return n
}

func (n *NoopValidator) Validate(_ context.Context, _ interface{}) error {
	return nil
}

func (n *NoopValidator) ValidateSlice(_ context.Context, _ interface{}) error {
	return nil
}

func NewNoopSpans() *NoopSpans {
	return &NoopSpans{}
}

func (n *NoopSpans) WithServiceManager(_ ServiceManager) Spans {
	return n
}

func (n *NoopSpans) New(ctx context.Context, _ ...domains.SpanConfig) (context.Context, *domains.Span) {
	return ctx, &domains.Span{
		Name:         "",
		Cid:          "",
		Resource:     "",
		Version:      "",
		OrgId:        "",
		Line:         0,
		FuncName:     "",
		FileName:     "",
		Custom:       map[string]interface{}{},
		InternalSpan: zipkin.SpanOrNoopFromContext(ctx),
	}
}

func (n *NoopSpans) Tracer() *zipkin.Tracer {
	tracer, _ := zipkin.NewTracer(nil, zipkin.WithNoopTracer(true))
	return tracer
}

func NewNoopUsersRepository() *NoopUsersRepository {
	return &NoopUsersRepository{}
}

func (n *NoopUsersRepository) WithServiceManager(_ ServiceManager) UsersRepository {
	return n
}

func (n *NoopUsersRepository) Search(_ context.Context, _ domains.UserSearch) ([]domains.User, error) {
	return nil, nil
}

func (n *NoopUsersRepository) Create(_ context.Context, _ *domains.User) error {
	return nil
}

func (n *NoopUsersRepository) Update(_ context.Context, _ *domains.User) error {
	return nil
}

func (n *NoopUsersRepository) Delete(_ context.Context, _ *domains.User) error {
	return nil
}

func (n *NoopUsersRepository) Migrate(_ context.Context) error {
	return nil
}

func NewNoopUserHandler() *NoopUsersHandler {
	return &NoopUsersHandler{}
}

func (n *NoopUsersHandler) WithServiceManager(_ ServiceManager) UsersHandler {
	return n
}

func (n *NoopUsersHandler) Create(_ context.Context, _ domains.UserCreate) (*domains.User, error) {
	return nil, nil
}

func (n *NoopUsersHandler) Update(_ context.Context, _ domains.UserUpdate) (*domains.User, error) {
	return nil, nil
}

func (n *NoopUsersHandler) List(_ context.Context) ([]domains.User, error) {
	return nil, nil
}

func (n *NoopUsersHandler) Get(_ context.Context, _ string) (*domains.User, error) {
	return nil, nil
}

func (n *NoopUsersHandler) Delete(_ context.Context, _ string) (*domains.User, error) {
	return nil, nil
}

func NewNoopEvents() *NoopEvents {
	return &NoopEvents{}
}

func (n *NoopEvents) WithServiceManager(_ ServiceManager) Events {
	return n
}

func (n *NoopEvents) Send(_ context.Context, _ domains.EventRequest) error {
	return nil
}

func (n *NoopEvents) SendBulk(_ context.Context, _ []domains.EventRequest) error {
	return nil
}
