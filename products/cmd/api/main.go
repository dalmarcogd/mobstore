package main

import (
	"errors"
	"fmt"
	http2 "net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dalmarcogd/mobstore/products/internal/discounts"
	"github.com/dalmarcogd/mobstore/products/internal/infra/database"
	"github.com/dalmarcogd/mobstore/products/internal/infra/environment"
	"github.com/dalmarcogd/mobstore/products/internal/infra/events"
	"github.com/dalmarcogd/mobstore/products/internal/infra/http"
	"github.com/dalmarcogd/mobstore/products/internal/infra/logger"
	"github.com/dalmarcogd/mobstore/products/internal/infra/spans"
	"github.com/dalmarcogd/mobstore/products/internal/infra/validator"
	"github.com/dalmarcogd/mobstore/products/internal/repositories/products"
	"github.com/dalmarcogd/mobstore/products/internal/services"
	"github.com/dalmarcogd/mobstore/products/internal/services/prodthandler"
)

func main() {
	ss := services.New().
		WithEnvironment(environment.New()).
		WithLogger(logger.New()).
		WithDiscounts(discounts.New()).
		WithValidator(validator.New()).
		WithSpans(spans.New()).
		WithHttpServer(http.New().WithAddress(":8080")).
		WithProductDatabase(database.New().WithProductDatabase()).
		WithProductsRepository(products.New()).
		WithProductsHandler(prodthandler.New()).
		WithProductsEvents(events.New().WithProductTopic())

	if err := ss.Init(); err != nil {
		ss.Logger().Fatal(ss.Context(), err.Error())
	}

	go func() {
		ss.Logger().Info(ss.Context(), "Http server stated")
		if err := ss.HttpServer().Run(); !errors.Is(err, http2.ErrServerClosed) {
			ss.Logger().Fatal(ss.Context(), err.Error())
		}
		ss.Logger().Info(ss.Context(), "Http server stopped")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-quit

	ss.Logger().Info(ss.Context(), fmt.Sprintf("Shutdown by %v", sig.String()))

	if err := ss.Close(); err != nil {
		ss.Logger().Fatal(ss.Context(), err.Error())
	}

	ss.Logger().Info(ss.Context(), "All services closed")
}
