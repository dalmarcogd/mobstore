package users

import (
	"context"
	"log"
	"reflect"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/dalmarcogd/mobstore/users/internal/domains"
	"github.com/dalmarcogd/mobstore/users/internal/infra/database"
	"github.com/dalmarcogd/mobstore/users/internal/infra/ptrs"
	"github.com/dalmarcogd/mobstore/users/internal/infra/validator"
	"github.com/dalmarcogd/mobstore/users/internal/services"
)

func newMock() (services.ServiceManager, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' occured on database mock connection", err)
	}

	mysql := database.New()
	serviceImpl := services.New().WithValidator(validator.New()).WithUserDatabase(mysql)
	mysql.WithServiceManager(serviceImpl)
	mysql.WithMasterClient(db).WithReplicaClient(db)

	return serviceImpl, mock
}

func TestProductRepository_Search(t *testing.T) {
	serviceManger, mock := newMock()

	repo := New()
	_ = repo.Init(context.Background())
	serviceManger.WithUsersRepository(repo)
	err := serviceManger.Init()
	if err != nil {
		t.Error(err)
	}

	_, err = repo.Search(context.Background(), domains.UserSearch{})
	if err == nil {
		t.Error("expected error with empty filter")
	}

	mock.ExpectBegin()
	mock.ExpectQuery("select id from users where deleted_at is null and last_name = ? order by id").RowsWillBeClosed().WithArgs("my-last-name").WillReturnRows(mock.NewRows([]string{"id"}).AddRow("my-id-1").AddRow("my-id-2"))
	mock.ExpectRollback()
	mock.ExpectClose()
	products, err := repo.Search(context.Background(), domains.UserSearch{Filter: domains.UserFilter{LastName: ptrs.String("my-last-name")}, Projection: domains.UserProjection{Id: true}})
	if err != nil {
		t.Error(err)
	}
	if len(products) != 2 {
		t.Errorf("expected one record from .Search got=%v", len(products))
	} else {
		if !reflect.DeepEqual(products[0], domains.User{Id: ptrs.String("my-id-1")}) {
			t.Errorf(".Search got=%v, expected=%v", products[0], domains.User{Id: ptrs.String("my-id-1")})
		} else if !reflect.DeepEqual(products[1], domains.User{Id: ptrs.String("my-id-2")}) {
			t.Errorf(".Search got=%v, expected=%v", products[1], domains.User{Id: ptrs.String("my-id-2")})
		}
	}
}

func BenchmarkProductRepository_Search(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		serviceManger, mock := newMock()

		repo := New()
		_ = repo.Init(context.Background())
		serviceManger.WithUsersRepository(repo)
		err := serviceManger.Init()
		if err != nil {
			b.Error(err)
		}
		mock.ExpectBegin()
		mock.ExpectQuery("select id from users where last_name = ? order by id").RowsWillBeClosed().WithArgs("my-last-name").WillReturnRows(mock.NewRows([]string{"id"}).AddRow("my-id-1").AddRow("my-id-2"))
		mock.ExpectRollback()
		mock.ExpectClose()
		products, err := repo.Search(context.Background(), domains.UserSearch{Filter: domains.UserFilter{LastName: ptrs.String("my-last-name")}, Projection: domains.UserProjection{Id: true}})
		if err != nil {
			b.Error(err)
		}
		if len(products) != 2 {
			b.Errorf("expected one record from .Search got=%v", len(products))
		} else {
			var foundA, foundB bool
			for _, user := range products {
				if reflect.DeepEqual(user, domains.User{Id: ptrs.String(strconv.Itoa(i + 1))}) {
					foundA = true
				}
				if reflect.DeepEqual(user, domains.User{Id: ptrs.String(strconv.Itoa(i + 2))}) {
					foundB = true
				}
			}
			if !foundA {
				b.Errorf("user %v not found", i+1)
			}
			if !foundB {
				b.Errorf("user %v not found", i+2)
			}
		}
	}
}
