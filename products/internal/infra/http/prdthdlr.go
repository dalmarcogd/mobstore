package http

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/dalmarcogd/mobstore/products/internal/domains"
	"github.com/dalmarcogd/mobstore/products/internal/infra/ctxs"
	"github.com/dalmarcogd/mobstore/products/internal/infra/ptrs"
	"github.com/dalmarcogd/mobstore/products/internal/infra/reqgetter"
)

//handleCreateV1 creates a Product
func (s *httpService) handleCreateV1(c echo.Context) error {
	ctx, span := s.ServiceManager().Spans().New(c.Request().Context())
	defer span.Finish()

	req := new(domains.ProductCreateRequestV1)
	if err := c.Bind(req); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during bind to ProductCreateRequestV1, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error()).SetInternal(err)
	}

	if err := s.ServiceManager().Validator().Validate(ctx, req); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during validate ProductCreateRequestV1, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error()).SetInternal(err)
	}

	productCreate := req.ProductCreate()
	product, err := s.ServiceManager().ProductsHandler().Create(ctx, productCreate)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during handle create Product, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	res := domains.ProductCreateResponseV1{
		Id:           ptrs.StringValue(product.Id),
		PriceInCents: ptrs.Int64Value(product.PriceInCents),
		Title:        ptrs.StringValue(product.Title),
		Description:  ptrs.StringValue(product.Description),
	}
	if err := s.ServiceManager().Validator().Validate(ctx, res); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during validate ProductCreateResponseV1, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	return c.JSON(http.StatusCreated, res)
}

//handleUpdateV1 updates a Product
func (s *httpService) handleUpdateV1(c echo.Context) error {
	ctx, span := s.ServiceManager().Spans().New(c.Request().Context())
	defer span.Finish()

	paramProductId := c.Param("productId")
	_, err := uuid.Parse(paramProductId)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during parse productId, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error()).SetInternal(err)
	}

	req := new(domains.ProductUpdateRequestV1)
	if err := c.Bind(req); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during bind to ProductCreateRequestV1, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error()).SetInternal(err)
	}

	if err := s.ServiceManager().Validator().Validate(ctx, req); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during validate ProductCreateRequestV1, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error()).SetInternal(err)
	}

	productUpdate := req.ProductUpdate(paramProductId)
	product, err := s.ServiceManager().ProductsHandler().Update(ctx, productUpdate)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during handle update Product, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	res := domains.ProductUpdateResponseV1{
		Id:           ptrs.StringValue(product.Id),
		PriceInCents: ptrs.Int64Value(product.PriceInCents),
		Title:        ptrs.StringValue(product.Title),
		Description:  ptrs.StringValue(product.Description),
	}
	if err := s.ServiceManager().Validator().Validate(ctx, res); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during validate ProductUpdateResponseV1, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}

//handleDeleteV1 delete a Product
func (s *httpService) handleDeleteV1(c echo.Context) error {
	ctx, span := s.ServiceManager().Spans().New(c.Request().Context())
	defer span.Finish()

	paramProductId := c.Param("productId")
	_, err := uuid.Parse(paramProductId)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during parse productId, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error()).SetInternal(err)
	}

	product, err := s.ServiceManager().ProductsHandler().Delete(ctx, paramProductId)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during delete product, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	resp := domains.ProductGetResponseV1{
		Id:           ptrs.StringValue(product.Id),
		PriceInCents: ptrs.Int64Value(product.PriceInCents),
		Title:        ptrs.StringValue(product.Title),
		Description:  product.Description,
	}

	return c.JSON(http.StatusOK, resp)
}

//handleGetV1 return list of Product
func (s *httpService) handleGetV1(c echo.Context) error {
	ctx, span := s.ServiceManager().Spans().New(c.Request().Context())
	defer span.Finish()

	userId := reqgetter.GetUserId(c.Request())
	if userId != nil {
		ctx = ctxs.ContextWithUserId(ctx, ptrs.StringValue(userId))
	}

	products, err := s.ServiceManager().ProductsHandler().List(ctx)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during get product, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	resp := domains.ProductListResponseV1{Products: make([]domains.ProductGetResponseV1, len(products))}
	for i, product := range products {
		resp.Products[i] = domains.ProductGetResponseV1{
			Id:           ptrs.StringValue(product.Id),
			PriceInCents: ptrs.Int64Value(product.PriceInCents),
			Title:        ptrs.StringValue(product.Title),
			Description:  product.Description,
			Discount:     product.Discount,
		}
	}

	return c.JSON(http.StatusOK, resp)
}

//handleGetV1 return specific Product
func (s *httpService) handleGetByIdV1(c echo.Context) error {
	ctx, span := s.ServiceManager().Spans().New(c.Request().Context())
	defer span.Finish()

	paramProductId := c.Param("productId")
	_, err := uuid.Parse(paramProductId)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during parse productId, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error()).SetInternal(err)
	}

	userId := reqgetter.GetUserId(c.Request())
	if userId != nil {
		ctx = ctxs.ContextWithUserId(ctx, ptrs.StringValue(userId))
	}

	product, err := s.ServiceManager().ProductsHandler().Get(ctx, paramProductId)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during get product, err=%v", err))
		span.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	resp := domains.ProductGetResponseV1{
		Id:           ptrs.StringValue(product.Id),
		PriceInCents: ptrs.Int64Value(product.PriceInCents),
		Title:        ptrs.StringValue(product.Title),
		Description:  product.Description,
		Discount:     product.Discount,
	}

	return c.JSON(http.StatusOK, resp)
}
