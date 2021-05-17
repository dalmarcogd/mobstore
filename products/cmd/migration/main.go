package main

import (
	"github.com/dalmarcogd/mobstore/products/internal/infra/database"
	"github.com/dalmarcogd/mobstore/products/internal/infra/environment"
	"github.com/dalmarcogd/mobstore/products/internal/infra/logger"
	"github.com/dalmarcogd/mobstore/products/internal/repositories/products"
	"github.com/dalmarcogd/mobstore/products/internal/services"
)

func main() {
	ss := services.New().
		WithEnvironment(environment.New()).
		WithLogger(logger.New()).
		WithProductDatabase(database.New().WithProductDatabase()).
		WithProductsRepository(products.New())

	if err := ss.Init(); err != nil {
		ss.Logger().Fatal(ss.Context(), err.Error())
	}

	err := ss.ProductsRepository().Migrate(ss.Context())
	if err != nil {
		ss.Logger().Fatal(ss.Context(), err.Error())
	}

	if err := ss.Close(); err != nil {
		ss.Logger().Fatal(ss.Context(), err.Error())
	}

	ss.Logger().Info(ss.Context(), "All services closed")
}
