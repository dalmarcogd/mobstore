package userhandler

import (
	"context"
	"fmt"

	"github.com/dalmarcogd/mobstore/users/internal/domains"
	"github.com/dalmarcogd/mobstore/users/internal/infra/ptrs"
)

func (s *productsHandler) Update(ctx context.Context, productUpdate domains.UserUpdate) (*domains.User, error) {
	ctx, span := s.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	if err := s.ServiceManager().Validator().Validate(ctx, productUpdate); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during validate User, err=%v", err))
		span.Error(err)
		return nil, err
	}

	user := &domains.User{
		Id:        ptrs.String(productUpdate.Id),
		FirstName: productUpdate.FirstName,
		LastName:  productUpdate.LastName,
		BirthDate: productUpdate.BirthDate,
	}

	err := s.ServiceManager().UsersRepository().Update(ctx, user)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during update User, err=%v", err))
		span.Error(err)
		return nil, err
	}

	user, err = s.Get(ctx, ptrs.StringValue(user.Id))
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during get User, err=%v", err))
		span.Error(err)
		return nil, err
	}

	err = s.ServiceManager().UsersEvents().Send(ctx, domains.NewUserUpdatedEventRequest(ctx, ptrs.StringValue(user.Id), ptrs.StringValue(user.FirstName), ptrs.StringValue(user.LastName), ptrs.TimeValue(user.BirthDate)))
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during send event for userId=%v, err=%v", ptrs.StringValue(user.Id), err))
		span.Error(err)
		return nil, err
	}

	return user, nil
}
