package userhandler

import (
	"context"
	"fmt"

	"github.com/dalmarcogd/mobstore/users/internal/domains"
	errors2 "github.com/dalmarcogd/mobstore/users/internal/infra/errors"
	"github.com/dalmarcogd/mobstore/users/internal/infra/ptrs"
)

func (s *productsHandler) Delete(ctx context.Context, userId string) (*domains.User, error) {
	ctx, span := s.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	users, err := s.ServiceManager().UsersRepository().Search(ctx, domains.UserSearch{
		Filter:     domains.UserFilter{Id: ptrs.String(userId)},
		Projection: domains.UserProjection{Id: true, FirstName: true, LastName: true, BirthDate: true},
	})
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during get User, err=%v", err))
		span.Error(err)
		return nil, err
	}

	if len(users) == 0 {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during get User, err=%v", errors2.UserNotFoundError))
		span.Error(errors2.UserNotFoundError)
		return nil, errors2.UserNotFoundError
	} else if len(users) > 1 {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during get User, err=%v", errors2.UserReturnWrongLenError))
		span.Error(errors2.UserReturnWrongLenError)
		return nil, errors2.UserReturnWrongLenError
	}

	user := &users[0]

	err = s.ServiceManager().UsersRepository().Delete(ctx, user)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during delete User, err=%v", err))
		span.Error(err)
		return nil, err
	}

	err = s.ServiceManager().UsersEvents().Send(ctx, domains.NewUserDeletedEventRequest(ctx, ptrs.StringValue(user.Id), ptrs.StringValue(user.FirstName), ptrs.StringValue(user.LastName), ptrs.TimeValue(user.BirthDate)))
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during send event for userId=%v, err=%v", ptrs.StringValue(user.Id), err))
		span.Error(err)
		return nil, err
	}

	return user, nil
}
