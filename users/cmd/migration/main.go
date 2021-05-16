package main

import (
	"github.com/dalmarcogd/mobstore/users/internal/infra/database"
	"github.com/dalmarcogd/mobstore/users/internal/infra/environment"
	"github.com/dalmarcogd/mobstore/users/internal/infra/logger"
	"github.com/dalmarcogd/mobstore/users/internal/repositories/users"
	"github.com/dalmarcogd/mobstore/users/internal/services"
)

func main() {
	ss := services.New().
		WithEnvironment(environment.New()).
		WithLogger(logger.New()).
		WithUserDatabase(database.New().WithUserDatabase()).
		WithUsersRepository(users.New())

	if err := ss.Init(); err != nil {
		ss.Logger().Fatal(ss.Context(), err.Error())
	}

	err := ss.UsersRepository().Migrate(ss.Context())
	if err != nil {
		ss.Logger().Fatal(ss.Context(), err.Error())
	}

	if err := ss.Close(); err != nil {
		ss.Logger().Fatal(ss.Context(), err.Error())
	}

	ss.Logger().Info(ss.Context(), "All services closed")
}
