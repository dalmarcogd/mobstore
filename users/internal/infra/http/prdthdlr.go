package http

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/dalmarcogd/mobstore/users/internal/domains"
	"github.com/dalmarcogd/mobstore/users/internal/infra/ptrs"
	"github.com/dalmarcogd/mobstore/users/internal/infra/times"
)

//handleCreateV1 creates a Product
func (s *httpService) handleCreateV1(c echo.Context) error {
	ctx, span := s.ServiceManager().Spans().New(c.Request().Context())
	defer span.Finish()

	req := new(domains.UserCreateRequestV1)
	if err := c.Bind(req); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during bind to UserCreateRequestV1, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error()).SetInternal(err)
	}

	if err := s.ServiceManager().Validator().Validate(ctx, req); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during validate UserCreateRequestV1, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error()).SetInternal(err)
	}

	userCreate := req.UserCreate()
	user, err := s.ServiceManager().UsersHandler().Create(ctx, userCreate)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during handle create User, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	res := domains.UserCreateResponseV1{
		Id:        ptrs.StringValue(user.Id),
		FirstName: ptrs.StringValue(user.FirstName),
		LastName:  ptrs.StringValue(user.LastName),
		BirthDate: times.JsonTime(ptrs.TimeValue(user.BirthDate)),
	}
	if err := s.ServiceManager().Validator().Validate(ctx, res); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during validate UserCreateResponseV1, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	return c.JSON(http.StatusCreated, res)
}

//handleUpdateV1 updates a User
func (s *httpService) handleUpdateV1(c echo.Context) error {
	ctx, span := s.ServiceManager().Spans().New(c.Request().Context())
	defer span.Finish()

	paramUserId := c.Param("userId")
	_, err := uuid.Parse(paramUserId)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during parse userId, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error()).SetInternal(err)
	}

	req := new(domains.UserUpdateRequestV1)
	if err := c.Bind(req); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during bind to UserCreateRequestV1, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error()).SetInternal(err)
	}

	if err := s.ServiceManager().Validator().Validate(ctx, req); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during validate UserCreateRequestV1, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error()).SetInternal(err)
	}

	userUpdate := req.UserUpdate(paramUserId)
	user, err := s.ServiceManager().UsersHandler().Update(ctx, userUpdate)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during handle update User, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	res := domains.UserUpdateResponseV1{
		Id:        ptrs.StringValue(user.Id),
		FirstName: ptrs.StringValue(user.FirstName),
		LastName:  ptrs.StringValue(user.LastName),
		BirthDate: times.JsonTime(ptrs.TimeValue(user.BirthDate)),
	}
	if err := s.ServiceManager().Validator().Validate(ctx, res); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during validate UserUpdateResponseV1, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}

//handleDeleteV1 delete a User
func (s *httpService) handleDeleteV1(c echo.Context) error {
	ctx, span := s.ServiceManager().Spans().New(c.Request().Context())
	defer span.Finish()

	paramUserId := c.Param("userId")
	_, err := uuid.Parse(paramUserId)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during parse userId, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error()).SetInternal(err)
	}

	user, err := s.ServiceManager().UsersHandler().Delete(ctx, paramUserId)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during delete user, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	resp := domains.UserGetResponseV1{
		Id:        ptrs.StringValue(user.Id),
		FirstName: ptrs.StringValue(user.FirstName),
		LastName:  ptrs.StringValue(user.LastName),
		BirthDate: times.JsonTime(ptrs.TimeValue(user.BirthDate)),
	}

	return c.JSON(http.StatusOK, resp)
}

//handleGetV1 return list of User
func (s *httpService) handleGetV1(c echo.Context) error {
	ctx, span := s.ServiceManager().Spans().New(c.Request().Context())
	defer span.Finish()

	users, err := s.ServiceManager().UsersHandler().List(ctx)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during get user, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	resp := domains.UserListResponseV1{Users: make([]domains.UserGetResponseV1, len(users))}
	for i, user := range users {
		resp.Users[i] = domains.UserGetResponseV1{
			Id:        ptrs.StringValue(user.Id),
			FirstName: ptrs.StringValue(user.FirstName),
			LastName:  ptrs.StringValue(user.LastName),
			BirthDate: times.JsonTime(ptrs.TimeValue(user.BirthDate)),
		}
	}

	return c.JSON(http.StatusOK, resp)
}

//handleGetV1 return specific User
func (s *httpService) handleGetByIdV1(c echo.Context) error {
	ctx, span := s.ServiceManager().Spans().New(c.Request().Context())
	defer span.Finish()

	paramUserId := c.Param("userId")
	_, err := uuid.Parse(paramUserId)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during parse userId, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error()).SetInternal(err)
	}

	user, err := s.ServiceManager().UsersHandler().Get(ctx, paramUserId)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during get user, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	resp := domains.UserGetResponseV1{
		Id:        ptrs.StringValue(user.Id),
		FirstName: ptrs.StringValue(user.FirstName),
		LastName:  ptrs.StringValue(user.LastName),
		BirthDate: times.JsonTime(ptrs.TimeValue(user.BirthDate)),
	}

	return c.JSON(http.StatusOK, resp)
}
