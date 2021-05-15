package prodthandler

import (
	"context"
	"fmt"

	"github.com/dalmarcogd/mobstore/products/internal/domains"
	"github.com/dalmarcogd/mobstore/products/internal/infra/ptrs"
)

//Create this handlers will create the Product
func (s *productsHandler) Create(ctx context.Context, productCreate domains.ProductCreate) (*domains.Product, error) {
	ctx, span := s.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	if err := s.ServiceManager().Validator().Validate(ctx, productCreate); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during validate Product, err=%v", err))
		span.Error(err)
		return nil, err
	}

	product := &domains.Product{
		PriceInCents: ptrs.Int64(productCreate.PriceInCents),
		Title:        ptrs.String(productCreate.Title),
		Description:  ptrs.String(productCreate.Description),
	}

	err := s.ServiceManager().ProductsRepository().Create(ctx, product)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during create Product, err=%v", err))
		span.Error(err)
		return nil, err
	}

	err = s.ServiceManager().ProductsEvents().Send(ctx, domains.NewProductCreatedEventRequest(ctx, ptrs.StringValue(product.Id), ptrs.StringValue(product.Title), ptrs.StringValue(product.Description), ptrs.Int64Value(product.PriceInCents)))
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during send event for productId=%v, err=%v", ptrs.StringValue(product.Id), err))
		span.Error(err)
		return nil, err
	}

	return product, nil
}
