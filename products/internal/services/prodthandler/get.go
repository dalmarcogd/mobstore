package prodthandler

import (
	"context"
	"fmt"

	"github.com/dalmarcogd/mobstore/products/internal/domains"
	errors2 "github.com/dalmarcogd/mobstore/products/internal/infra/errors"
	"github.com/dalmarcogd/mobstore/products/internal/infra/ptrs"
)

func (s *productsHandler) List(ctx context.Context) ([]domains.Product, error) {
	ctx, span := s.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	products, err := s.ServiceManager().ProductsRepository().Search(ctx, domains.ProductSearch{
		Projection: domains.ProductProjection{Id: true, PriceInCents: true, Title: true, Description: true},
	})
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during get Product, err=%v", err))
		span.Error(err)
		return nil, err
	}

	return products, nil
}

func (s *productsHandler) Get(ctx context.Context, productId string) (*domains.Product, error) {
	ctx, span := s.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	products, err := s.ServiceManager().ProductsRepository().Search(ctx, domains.ProductSearch{
		Filter:     domains.ProductFilter{Id: ptrs.String(productId)},
		Projection: domains.ProductProjection{Id: true, PriceInCents: true, Title: true, Description: true},
	})
	if err != nil {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during get Product, err=%v", err))
		span.Error(err)
		return nil, err
	}

	if len(products) == 0 {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during get Product, err=%v", errors2.ProductNotFoundError))
		span.Error(errors2.ProductNotFoundError)
		return nil, errors2.ProductNotFoundError
	} else if len(products) > 1 {
		s.ServiceManager().Logger().Error(ctx, fmt.Sprintf("Error during get Product, err=%v", errors2.ProductReturnWrongLenError))
		span.Error(errors2.ProductReturnWrongLenError)
		return nil, errors2.ProductReturnWrongLenError
	}

	return &products[0], nil
}
