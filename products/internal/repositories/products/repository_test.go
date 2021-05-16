package products

import (
	"context"
	"log"
	"reflect"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/dalmarcogd/mobstore/products/internal/domains"
	"github.com/dalmarcogd/mobstore/products/internal/infra/database"
	"github.com/dalmarcogd/mobstore/products/internal/infra/ptrs"
	"github.com/dalmarcogd/mobstore/products/internal/infra/validator"
	"github.com/dalmarcogd/mobstore/products/internal/services"
)

func newMock() (services.ServiceManager, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' occured on database mock connection", err)
	}

	mysql := database.New()
	serviceImpl := services.New().WithValidator(validator.New()).WithProductDatabase(mysql)
	mysql.WithServiceManager(serviceImpl)
	mysql.WithMasterClient(db).WithReplicaClient(db)

	return serviceImpl, mock
}

func TestCardRepository_Search(t *testing.T) {
	serviceManger, mock := newMock()

	repo := New()
	_ = repo.Init(context.Background())
	serviceManger.WithProductsRepository(repo)
	err := serviceManger.Init()
	if err != nil {
		t.Error(err)
	}

	_, err = repo.Search(context.Background(), domains.ProductSearch{})
	if err == nil {
		t.Error("expected error with empty filter")
	}

	mock.ExpectBegin()
	mock.ExpectQuery("select id from products where title = ? order by id").RowsWillBeClosed().WithArgs("my-title").WillReturnRows(mock.NewRows([]string{"id"}).AddRow("my-id-1").AddRow("my-id-2"))
	mock.ExpectRollback()
	mock.ExpectClose()
	products, err := repo.Search(context.Background(), domains.ProductSearch{Filter: domains.ProductFilter{Title: ptrs.String("my-title")}, Projection: domains.ProductProjection{Id: true}})
	if err != nil {
		t.Error(err)
	}
	if len(products) != 2 {
		t.Errorf("expected one record from .Search got=%v", len(products))
	} else {
		if !reflect.DeepEqual(products[0], domains.Product{Id: ptrs.String("my-id-1")}) {
			t.Errorf(".Search got=%v, expected=%v", products[0], domains.Product{Id: ptrs.String("my-id-1")})
		} else if !reflect.DeepEqual(products[1], domains.Product{Id: ptrs.String("my-id-2")}) {
			t.Errorf(".Search got=%v, expected=%v", products[1], domains.Product{Id: ptrs.String("my-id-2")})
		}
	}
}

func BenchmarkCardRepository_Search(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		serviceManger, mock := newMock()

		repo := New()
		_ = repo.Init(context.Background())
		serviceManger.WithProductsRepository(repo)
		err := serviceManger.Init()
		if err != nil {
			b.Error(err)
		}
		mock.ExpectBegin()
		mock.ExpectQuery("select id from products where title = ? order by id").RowsWillBeClosed().WithArgs("my-title").WillReturnRows(mock.NewRows([]string{"id"}).AddRow("my-id-1").AddRow("my-id-2"))
		mock.ExpectRollback()
		mock.ExpectClose()
		products, err := repo.Search(context.Background(), domains.ProductSearch{Filter: domains.ProductFilter{Title: ptrs.String("my-title")}, Projection: domains.ProductProjection{Id: true}})
		if err != nil {
			b.Error(err)
		}
		if len(products) != 2 {
			b.Errorf("expected one record from .Search got=%v", len(products))
		} else {
			var foundA, foundB bool
			for _, product := range products {
				if reflect.DeepEqual(product, domains.Product{Id: ptrs.String(strconv.Itoa(i + 1))}) {
					foundA = true
				}
				if reflect.DeepEqual(product, domains.Product{Id: ptrs.String(strconv.Itoa(i + 2))}) {
					foundB = true
				}
			}
			if !foundA {
				b.Errorf("product %v not found", i+1)
			}
			if !foundB {
				b.Errorf("product %v not found", i+2)
			}
		}
	}
}
