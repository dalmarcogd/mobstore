package prodthandler

import (
	"context"
	"fmt"

	"github.com/dalmarcogd/mobstore/products/internal/domains"
	"github.com/dalmarcogd/mobstore/products/internal/infra/ptrs"
)

func (s *productsHandler) Update(ctx context.Context, productUpdate domains.ProductUpdate) (*domains.Product, error) {
	ctx, span := s.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	if err := s.ServiceManager().Validator().Validate(ctx, productUpdate); err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during validate Product, err=%v", err))
		span.Error(err)
		return nil, err
	}

	product := &domains.Product{
		Id:           ptrs.String(productUpdate.Id),
		PriceInCents: productUpdate.PriceInCents,
		Title:        productUpdate.Title,
		Description:  productUpdate.Description,
	}

	err := s.ServiceManager().ProductsRepository().Update(ctx, product)
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during update Product, err=%v", err))
		span.Error(err)
		return nil, err
	}

	product, err = s.Get(ctx, ptrs.StringValue(product.Id))
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during get Product, err=%v", err))
		span.Error(err)
		return nil, err
	}

	err = s.ServiceManager().ProductsEvents().Send(ctx, domains.NewProductUpdatedEventRequest(ctx, ptrs.StringValue(product.Id), ptrs.StringValue(product.Title), ptrs.StringValue(product.Description), ptrs.Int64Value(product.PriceInCents)))
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during send event for productId=%v, err=%v", ptrs.StringValue(product.Id), err))
		span.Error(err)
		return nil, err
	}

	return product, nil
}
