package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	zipkinsql "github.com/openzipkin-contrib/zipkin-go-sql"

	"github.com/dalmarcogd/mobstore/users/internal/infra/ctxs"
	"github.com/dalmarcogd/mobstore/users/internal/infra/errors"
	"github.com/dalmarcogd/mobstore/users/internal/services"
)

const (
	userDatabase DB = iota
)

const mysqlDuplicateEntry = 1062

type (
	DB int

	mySQL struct {
		dbMaster       *sql.DB
		dbReplica      *sql.DB
		ctx            context.Context
		serviceManager services.ServiceManager
		database       DB
	}
	mySQLTransaction struct {
		*mySQL
		ctx context.Context
	}
)

func New() *mySQL {
	return &mySQL{}
}

func (m *mySQL) ServiceManager() services.ServiceManager {
	return m.serviceManager
}

func (m *mySQL) WithServiceManager(c services.ServiceManager) services.Database {
	m.serviceManager = c
	return m
}

func (m *mySQL) WithMasterClient(db *sql.DB) services.Database {
	m.dbMaster = db
	return m
}

func (m *mySQL) WithReplicaClient(db *sql.DB) services.Database {
	m.dbReplica = db
	return m
}

func (m *mySQL) WithUserDatabase() services.Database {
	m.database = userDatabase
	return m
}

func (m *mySQL) Init(ctx context.Context) error {
	logger := m.ServiceManager().Logger()

	m.ctx = ctx

	var databaseMasterDsn, databaseReplicaDsn string
	switch m.database {
	case userDatabase:
		databaseMasterDsn = m.ServiceManager().Environment().UserDatabaseDsn()
		databaseReplicaDsn = m.ServiceManager().Environment().UserReplicaDatabaseDsn()
	default:
		return errors.DatabaseDsnIsRequiredError
	}

	driverName, err := zipkinsql.Register("mysql", m.ServiceManager().Spans().Tracer(), zipkinsql.WithAllTraceOptions())
	if err != nil {
		return err
	}

	if m.dbMaster == nil && databaseMasterDsn != "" {
		dbMaster, err := sql.Open(driverName, databaseMasterDsn)
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("Error on trying to connect with database=%v, err=%v", m.database, err))
			return err
		}

		dbMaster.SetConnMaxLifetime(time.Minute * 3)
		dbMaster.SetMaxOpenConns(10)
		dbMaster.SetMaxIdleConns(4)

		m.dbMaster = dbMaster
	}

	if m.dbReplica == nil && databaseReplicaDsn != "" {
		dbReplica, err := sql.Open(driverName, databaseReplicaDsn)
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("Error on trying to connect with database=%v, err=%v", m.database, err))
			return err
		}

		dbReplica.SetConnMaxLifetime(time.Minute * 3)
		dbReplica.SetMaxOpenConns(10)
		dbReplica.SetMaxIdleConns(4)

		m.dbReplica = dbReplica
	}

	return nil
}

func (m *mySQL) Close() error {
	var err error
	if errC := m.dbMaster.Close(); errC != nil {
		err = fmt.Errorf("%v - %v", err, errC)
	}
	if errC := m.dbReplica.Close(); errC != nil {
		err = fmt.Errorf("%v - %v", err, errC)
	}
	return err
}

func (m *mySQL) Readiness(ctx context.Context) error {
	if m.dbMaster != nil {
		if err := m.dbMaster.PingContext(ctx); err != nil {
			return err
		}
	}
	if m.dbReplica != nil {
		if err := m.dbReplica.PingContext(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (m *mySQL) Liveness(ctx context.Context) error {
	return m.Readiness(ctx)
}

func (m *mySQL) OpenTransactionMaster(ctx context.Context) (context.Context, error) {
	ctx, span := m.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	tx := ctxs.GetTransactionFromContext(ctx)
	if tx == nil {
		tx, err := m.dbMaster.BeginTx(ctx, &sql.TxOptions{})
		return ctxs.ContextWithTransaction(ctx, tx), err
	}
	return ctx, nil
}

func (m *mySQL) TransactionMaster(ctx context.Context, f func(tx services.DatabaseTransaction) error) error {
	ctx, span := m.ServiceManager().Spans().New(ctx)
	defer span.Finish()
	logger := m.ServiceManager().Logger()

	ctx, err := m.OpenTransactionMaster(ctx)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("Error opening transaciton err=%v", err))
		span.Error(err)
		return err
	}

	defer func(c context.Context, e *error) {
		if err := m.CloseTransaction(c, *e); err != nil {
			logger.Error(c, fmt.Sprintf("Error closing transaciton err=%v", err))
			span.Error(err)
		}
	}(ctx, &err)

	err = f(&mySQLTransaction{ctx: ctx, mySQL: m})
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("Error executing transaciton function err=%v", err))
		span.Error(err)
		return err
	}

	return nil
}

func (m *mySQL) OpenTransactionReplica(ctx context.Context) (context.Context, error) {
	ctx, span := m.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	tx := ctxs.GetTransactionFromContext(ctx)
	if tx == nil {
		tx, err := m.dbReplica.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
		return ctxs.ContextWithTransaction(ctx, tx), err
	}
	return ctx, nil
}

func (m *mySQL) TransactionReplica(ctx context.Context, f func(tx services.DatabaseTransaction) error) error {
	ctx, span := m.ServiceManager().Spans().New(ctx)
	defer span.Finish()
	logger := m.ServiceManager().Logger()

	ctx, err := m.OpenTransactionReplica(ctx)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("Error opening transaciton err=%v", err))
		span.Error(err)
		return err
	}

	defer func(c context.Context) {
		if err := m.rollback(c); err != nil {
			logger.Error(c, fmt.Sprintf("Error closing transaciton err=%v", err))
			span.Error(err)
		}
	}(ctx)

	err = f(&mySQLTransaction{ctx: ctx, mySQL: m})
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("Error executing transaciton function err=%v", err))
		span.Error(err)
		return err
	}

	return nil
}

func (m *mySQL) CloseTransaction(ctx context.Context, err error) error {
	ctx, span := m.ServiceManager().Spans().New(ctx)
	defer span.Finish()
	if err == nil {
		if errCommit := m.commit(ctx); errCommit != nil {
			span.Error(errCommit)
			return errCommit
		}
	} else if errRollback := m.rollback(ctx); errRollback != nil {
		span.Error(errRollback)
		return errRollback
	}
	return nil
}

func (m *mySQL) commit(ctx context.Context) error {
	ctx, span := m.ServiceManager().Spans().New(ctx)
	defer span.Finish()
	tx := ctxs.GetTransactionFromContext(ctx)
	if tx == nil {
		return errors.DatabaseTransactionAtContextNotFoundError
	}
	if errCommit := tx.Commit(); errCommit != nil {
		span.Error(errCommit)
		return errCommit
	}
	return nil
}

func (m *mySQL) rollback(ctx context.Context) error {
	ctx, span := m.ServiceManager().Spans().New(ctx)
	defer span.Finish()
	tx := ctxs.GetTransactionFromContext(ctx)
	if tx == nil {
		return errors.DatabaseTransactionAtContextNotFoundError
	}
	if errRollback := tx.Rollback(); errRollback != nil {
		span.Error(errRollback)
		return errRollback
	}
	return nil
}

func (m *mySQLTransaction) Insert(query string, args ...interface{}) (sql.Result, error) {
	ctx, span := m.ServiceManager().Spans().New(m.ctx)
	defer span.Finish()
	logger := m.ServiceManager().Logger()

	result, err := ctxs.GetTransactionFromContext(ctx).ExecContext(ctx, query, args...)
	if err != nil {
		if mysqlError, ok := err.(*mysql.MySQLError); ok {
			switch mysqlError.Number {
			case mysqlDuplicateEntry:
				err = errors.DatabaseUniqueEntryViolationError
			}
		}
		span.Error(err)
		logger.Error(ctx, fmt.Sprintf("Error executing insert err=%v", err.Error()))
		return nil, err
	}

	return result, nil
}

func (m *mySQLTransaction) Update(query string, args ...interface{}) (sql.Result, error) {
	ctx, span := m.ServiceManager().Spans().New(m.ctx)
	defer span.Finish()
	logger := m.ServiceManager().Logger()

	result, err := ctxs.GetTransactionFromContext(ctx).ExecContext(ctx, query, args...)
	if err != nil {
		span.Error(err)
		logger.Error(ctx, fmt.Sprintf("Error executing update err=%v", err.Error()))
		return nil, err
	}

	return result, nil
}

func (m *mySQLTransaction) Get(query string, args ...interface{}) (*sql.Rows, error) {
	ctx, span := m.ServiceManager().Spans().New(m.ctx)
	defer span.Finish()
	logger := m.ServiceManager().Logger()

	rows, err := ctxs.GetTransactionFromContext(ctx).QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("Error executing query err=%v", err))
		span.Error(err)
		return nil, err
	}

	return rows, nil
}

func (m *mySQLTransaction) Migrate(script string) (sql.Result, error) {
	ctx, span := m.ServiceManager().Spans().New(m.ctx)
	defer span.Finish()
	logger := m.ServiceManager().Logger()

	result, err := ctxs.GetTransactionFromContext(ctx).ExecContext(ctx, script)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("Error executing script err=%v", err))
		span.Error(err)
		return nil, err
	}

	return result, nil
}
