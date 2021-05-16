package database

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	errors2 "github.com/dalmarcogd/mobstore/users/internal/infra/errors"
	"github.com/dalmarcogd/mobstore/users/internal/services"
)

func newMock(t *testing.T) (*mySQL, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' occured on database mock connection", err)
	}
	mysql := New()
	serviceImpl := services.New().WithUserDatabase(mysql)

	if err := mysql.Init(context.Background()); err != nil {
		t.Error(err)
	}

	mysql.dbMaster = db
	mysql.dbReplica = db
	mysql.ctx = context.Background()
	mysql.serviceManager = serviceImpl

	mysql.WithUserDatabase()
	if mysql.database != userDatabase {
		t.Error("expected userDatabase")
	}

	if err := mysql.Init(context.Background()); err != nil {
		t.Error(err)
	}

	if err := mysql.Init(context.Background()); err != nil {
		t.Error(err)
	}

	return mysql, mock
}

func TestMySQL_Get(t *testing.T) {
	ctx := context.Background()
	mysql, mock := newMock(t)

	query := "SELECT Id, Org_ID, Program_ID, Ica, BulkId, Endpoint, Enable FROM CardAutomaticUpdater WHERE id = ?"

	rows := sqlmock.NewRows([]string{"Id", "Org_ID", "Program_ID", "Ica", "BulkId", "Endpoint", "Enable"}).
		AddRow(1, "org-18726-ashgd-1qsd", 222, 1234, 12312, "https://google.com", true).
		AddRow(2, "org-sdfsfdg-ashgd-1qsd", 111, 4321, 2354, "https://yahoo.com", true)

	mock.ExpectBegin()
	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
	mock.ExpectCommit()
	mock.ExpectClose()

	err := mysql.TransactionMaster(ctx, func(tx services.DatabaseTransaction) error {
		if r, err := tx.Get(query, 1); err != nil || r == nil {
			t.Error(err)
			return err
		}
		return nil
	})

	if err != nil {
		t.Error(err)
	}

	if err := mysql.Close(); err != nil {
		t.Error(err)
	}
}

func TestMySQL_Insert(t *testing.T) {
	ctx := context.Background()
	mysql, mock := newMock(t)
	query := "INSERT INTO CardAutomaticUpdater(Id, Org_ID, Program_ID, Ica, BulkId, Endpoint, Enable) VALUES(?,?,?,?,?,?,?)"

	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs(1, "org-18726-ashgd-1qsd", 222, 1234, 12312, "https://google.com", true).
		WillReturnResult(sqlmock.NewResult(2, 2))
	mock.ExpectCommit()

	err := mysql.TransactionMaster(ctx, func(tx services.DatabaseTransaction) error {
		if r, err := tx.Insert(query, 1, "org-18726-ashgd-1qsd", 222, 1234, 12312, "https://google.com", true); err != nil || r == nil {
			t.Error(err)
			return err
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
	mock.ExpectClose()

	if err := mysql.Close(); err != nil {
		t.Error(err)
	}
}

func TestMySQL_Update(t *testing.T) {
	ctx := context.Background()
	mysql, mock := newMock(t)
	query := "UPDATE CardAutomaticUpdater SET Org_ID = ? WHERE Id = ?"

	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs("org-18726-ashgd-1qsd", 1).
		WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()

	err := mysql.TransactionMaster(ctx, func(tx services.DatabaseTransaction) error {
		if r, err := tx.Update(query, "org-18726-ashgd-1qsd", 1); err != nil || r == nil {
			t.Error(err)
			return err
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
	mock.ExpectClose()

	if err := mysql.Close(); err != nil {
		t.Error(err)
	}
}

func TestMySQL_GetError(t *testing.T) {
	ctx := context.Background()
	mysql, mock := newMock(t)

	query := "SELECT Id, Org_ID, Program_ID, Ica, BulkId, Endpoint, Enable FROM CardAutomaticUpdater WHERE id = ?"

	rows := sqlmock.NewRows([]string{"Id", "Org_ID", "Program_ID", "Ica", "BulkId", "Endpoint", "Enable"}).
		AddRow(1, "org-18726-ashgd-1qsd", 222, 1234, 12312, "https://google.com", true).
		AddRow(2, "org-sdfsfdg-ashgd-1qsd", 111, 4321, 2354, "https://yahoo.com", true)

	mock.ExpectBegin()
	mock.ExpectQuery(query).WithArgs(1, 2, 3).WillReturnRows(rows)
	mock.ExpectClose()
	t.Run("Wrong numbers of arguments", func(tt *testing.T) {
		err := mysql.TransactionMaster(ctx, func(tx services.DatabaseTransaction) error {
			if r, err := tx.Get(query, 1); err == nil || r != nil {
				tt.Error(err)
				return err
			}
			return nil
		})
		if err != nil {
			t.Error(err)
		}
	})

	mysql, mock = newMock(t)
	mock.ExpectBegin()
	mock.ExpectQuery(query).WithArgs(1).WillReturnError(errors.New("error test"))
	mock.ExpectClose()

	t.Run("Generic error", func(tt *testing.T) {
		err := mysql.TransactionMaster(ctx, func(tx services.DatabaseTransaction) error {
			if _, err := tx.Get(query, 1); err == nil || err.Error() != "error test" {
				tt.Error(err)
				return err
			}
			return nil
		})
		if err != nil {
			t.Error(err)
		}
	})
	mock.ExpectClose()

	if err := mysql.Close(); err != nil {
		t.Error(err)
	}
}

func TestMySQL_InsertError(t *testing.T) {
	ctx := context.Background()
	mysql, mock := newMock(t)
	query := "INSERT INTO CardAutomaticUpdater(Id, Org_ID, Program_ID, Ica, BulkId, Endpoint, Enable) VALUES(?,?,?,?,?,?,?)"

	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs(1, "org-18726-ashgd-1qsd", 222, 1234, 12312, "https://google.com", true).
		WillReturnError(errors.New("some error"))
	mock.ExpectRollback()
	mock.ExpectClose()

	err := mysql.TransactionMaster(ctx, func(tx services.DatabaseTransaction) error {
		if _, err := tx.Insert(query, 1, "org-18726-ashgd-1qsd", 222, 1234, 12312, "https://google.com", true); err != nil {
			return err
		}
		t.Error("expected error")
		return nil
	})
	if err == nil {
		t.Error("expected error")
	}

	if err := mysql.Close(); err != nil {
		t.Error(err)
	}

	mysql, mock = newMock(t)
	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs(1, "org-18726-ashgd-1qsd", 222, 1234, 12312, "https://google.com", true).
		WillReturnError(errors.New("error test"))
	mock.ExpectRollback()
	mock.ExpectClose()
	err = mysql.TransactionMaster(ctx, func(tx services.DatabaseTransaction) error {
		if _, err := tx.Insert(query, 1, "org-18726-ashgd-1qsd", 222, 1234, 12312, "https://google.com", true); err != nil {
			return err
		}
		t.Error("expected error")
		return nil
	})
	if err == nil {
		t.Error("expected error")
	}

	if err := mysql.Close(); err != nil {
		t.Error(err)
	}
}

func TestMySQL_UpdateError(t *testing.T) {
	ctx := context.Background()
	mysql, mock := newMock(t)
	query := "UPDATE CardAutomaticUpdater SET Org_ID = ? WHERE Id = ?"

	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs("org-18726-ashgd-1qsd", 1).WillReturnError(errors.New("some error"))
	mock.ExpectRollback()
	mock.ExpectClose()

	err := mysql.TransactionMaster(ctx, func(tx services.DatabaseTransaction) error {
		if _, err := tx.Update(query, "org-18726-ashgd-1qsd", 1); err != nil {
			return err
		}
		t.Error("expected error")
		return nil
	})
	if err == nil {
		t.Error("expected error")
	}

	if err := mysql.Close(); err != nil {
		t.Error(err)
	}

	mysql, mock = newMock(t)
	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs("org-18726-ashgd-1qsd", 1).
		WillReturnError(errors.New("error test"))
	mock.ExpectRollback()
	mock.ExpectClose()

	err = mysql.TransactionMaster(ctx, func(tx services.DatabaseTransaction) error {
		if _, err := tx.Update(query, "org-18726-ashgd-1qsd", 1); err != nil {
			return err
		}
		t.Error("expected error")
		return nil
	})
	if err == nil {
		t.Error("expected error")
	}

	mysql, mock = newMock(t)
	mock.ExpectBegin().WillReturnError(errors.New("some error"))
	mock.ExpectClose()

	err = mysql.TransactionMaster(ctx, func(tx services.DatabaseTransaction) error {
		if _, err := tx.Update(query, "org-18726-ashgd-1qsd", 1); err != nil {
			return err
		}
		t.Error("expected error")
		return nil
	})
	if err == nil {
		t.Error("expected error")
	}

	if err := mysql.Close(); err != nil {
		t.Error(err)
	}
}

func TestMySQLTransaction_Get(t *testing.T) {
	ctx := context.Background()
	mysql, mock := newMock(t)
	query := "select * from CardAutomaticUpdater"

	mock.ExpectBegin()
	mock.ExpectQuery(query).WillReturnRows(
		sqlmock.NewRows([]string{"1", "2", "3"}))
	mock.ExpectRollback()

	err := mysql.TransactionReplica(ctx, func(tx services.DatabaseTransaction) error {
		rows, err := tx.Get(query)
		if err != nil {
			return err
		}
		defer func() {
			if err := rows.Close(); err != nil {
				t.Error(err)
			}
		}()
		return nil
	})
	if err != nil {
		t.Error(err)
	}

	mock.ExpectBegin().WillReturnError(errors.New("some error"))
	mock.ExpectClose()

	err = mysql.TransactionReplica(ctx, func(tx services.DatabaseTransaction) error {
		rows, err := tx.Get(query)
		if err != nil {
			return err
		}
		defer func() {
			if err := rows.Close(); err != nil {
				t.Error(err)
			}
		}()
		return nil
	})
	if err == nil {
		t.Error("expected error")
	}

	if err := mysql.Close(); err != nil {
		t.Error(err)
	}

	mysql, mock = newMock(t)

	mock.ExpectBegin()
	mock.ExpectQuery(query).WillReturnError(errors.New("some error"))
	mock.ExpectRollback()
	mock.ExpectClose()

	err = mysql.TransactionReplica(ctx, func(tx services.DatabaseTransaction) error {
		rows, err := tx.Get(query)
		if err != nil {
			return err
		}
		t.Error("expected error")
		defer func() {
			if err := rows.Close(); err != nil {
				t.Error(err)
			}
		}()
		return nil
	})
	if err == nil {
		t.Error("expected error")
	}

	if err := mysql.Close(); err != nil {
		t.Error(err)
	}
}

func TestMySQL_CloseTransaction(t *testing.T) {
	mysql, mock := newMock(t)
	err := mysql.CloseTransaction(context.Background(), nil)
	if !errors.Is(err, errors2.DatabaseTransactionAtContextNotFoundError) {
		t.Error(err)
	}
	err = mysql.CloseTransaction(context.Background(), errors.New("some error"))
	if !errors.Is(err, errors2.DatabaseTransactionAtContextNotFoundError) {
		t.Error(err)
	}

	mock.ExpectBegin()
	mock.ExpectCommit().WillReturnError(errors.New("some error"))
	ctx, err := mysql.OpenTransactionReplica(context.Background())
	if err != nil {
		t.Error(err)
	}
	err = mysql.CloseTransaction(ctx, nil)
	if err == nil {
		t.Error("expected error")
	}

	mock.ExpectBegin()
	mock.ExpectRollback().WillReturnError(errors.New("some error"))
	ctx, err = mysql.OpenTransactionReplica(context.Background())
	if err != nil {
		t.Error(err)
	}

	err = mysql.CloseTransaction(ctx, errors.New("some error2"))
	if err == nil {
		t.Error("expected error")
	}
}

func TestMySQL_OpenTransaction(t *testing.T) {
	mysql, mock := newMock(t)
	mock.ExpectBegin().WillReturnError(errors.New("some error"))
	_, err := mysql.OpenTransactionMaster(context.Background())
	if err == nil {
		t.Error("expected error")
	}

	mock.ExpectBegin().WillReturnError(errors.New("some error"))
	_, err = mysql.OpenTransactionReplica(context.Background())
	if err == nil {
		t.Error("expected error")
	}
}

func TestMySQL_Close(t *testing.T) {
	mysql, mock := newMock(t)
	mock.ExpectClose().WillReturnError(errors.New("some error"))
	if err := mysql.Close(); err == nil {
		t.Error("expected close")
	}
}
