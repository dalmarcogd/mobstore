package userhandler

import (
	"context"
	"fmt"

	"github.com/dalmarcogd/mobstore/users/internal/domains"
	"github.com/dalmarcogd/mobstore/users/internal/infra/ptrs"
)

//Create this handlers will create the User
func (s *productsHandler) Create(ctx context.Context, userCreate domains.UserCreate) (*domains.User, error) {
	ctx, span := s.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	if err := s.ServiceManager().Validator().Validate(ctx, userCreate); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during validate User, err=%v", err))
		span.Error(err)
		return nil, err
	}

	user := &domains.User{
		FirstName: ptrs.String(userCreate.FirstName),
		LastName:  ptrs.String(userCreate.LastName),
		BirthDate: ptrs.Time(userCreate.BirthDate),
	}

	err := s.ServiceManager().UsersRepository().Create(ctx, user)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during create User, err=%v", err))
		span.Error(err)
		return nil, err
	}

	err = s.ServiceManager().UsersEvents().Send(ctx, domains.NewUserCreatedEventRequest(ctx, ptrs.StringValue(user.Id), ptrs.StringValue(user.FirstName), ptrs.StringValue(user.LastName), ptrs.TimeValue(user.BirthDate)))
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during send event for userId=%v, err=%v", ptrs.StringValue(user.Id), err))
		span.Error(err)
		return nil, err
	}

	return user, nil
}
