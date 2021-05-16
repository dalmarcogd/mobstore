package userhandler

import (
	"context"
	"fmt"

	"github.com/dalmarcogd/mobstore/users/internal/domains"
	errors2 "github.com/dalmarcogd/mobstore/users/internal/infra/errors"
	"github.com/dalmarcogd/mobstore/users/internal/infra/ptrs"
)

func (s *productsHandler) List(ctx context.Context) ([]domains.User, error) {
	ctx, span := s.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	products, err := s.ServiceManager().UsersRepository().Search(ctx, domains.UserSearch{
		Projection: domains.UserProjection{Id: true, FirstName: true, LastName: true, BirthDate: true},
	})
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during get User, err=%v", err))
		span.Error(err)
		return nil, err
	}

	return products, nil
}

func (s *productsHandler) Get(ctx context.Context, productId string) (*domains.User, error) {
	ctx, span := s.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	products, err := s.ServiceManager().UsersRepository().Search(ctx, domains.UserSearch{
		Filter:     domains.UserFilter{Id: ptrs.String(productId)},
		Projection: domains.UserProjection{Id: true, FirstName: true, LastName: true, BirthDate: true},
	})
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during get User, err=%v", err))
		span.Error(err)
		return nil, err
	}

	if len(products) == 0 {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during get User, err=%v", errors2.UserNotFoundError))
		span.Error(errors2.UserNotFoundError)
		return nil, errors2.UserNotFoundError
	} else if len(products) > 1 {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during get User, err=%v", errors2.UserReturnWrongLenError))
		span.Error(errors2.UserReturnWrongLenError)
		return nil, errors2.UserReturnWrongLenError
	}

	return &products[0], nil
}
