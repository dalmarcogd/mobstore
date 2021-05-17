package discounts

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"

	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"

	"github.com/dalmarcogd/mobstore/products/internal/discounts/discountsgrpc"
	"github.com/dalmarcogd/mobstore/products/internal/domains"
	"github.com/dalmarcogd/mobstore/products/internal/domains/domainsgrpc"
	"github.com/dalmarcogd/mobstore/products/internal/infra/ctxs"
	"github.com/dalmarcogd/mobstore/products/internal/infra/ptrs"
	"github.com/dalmarcogd/mobstore/products/internal/services"
)

type (
	discountsService struct {
		services.NoopHealth
		serviceManager  services.ServiceManager
		ctx             context.Context
		address         string
		clientConn      grpc.ClientConnInterface
		discountsClient discountsgrpc.DiscountsClient
	}
)

func New() *discountsService {
	return &discountsService{}
}

func (s *discountsService) WithAddress(address string) services.Discounts {
	s.address = address
	return s
}

func (s *discountsService) Init(ctx context.Context) error {
	s.ctx = ctx
	if s.address == "" {
		s.address = s.ServiceManager().Environment().DiscountsAddress()
	}
	if s.clientConn == nil && s.discountsClient == nil {
		dial, err := grpc.DialContext(
			s.ctx,
			s.address,
			grpc.WithInsecure(),
			grpc.WithStatsHandler(zipkingrpc.NewClientHandler(s.ServiceManager().Spans().Tracer())),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                20,
				Timeout:             10,
				PermitWithoutStream: true,
			}),
		)
		if err != nil {
			return err
		}
		s.clientConn = dial

		s.discountsClient = discountsgrpc.NewDiscountsClient(s.clientConn)
	}

	return nil
}

func (s *discountsService) Close() error {
	s.clientConn = nil
	return nil
}

func (s *discountsService) WithServiceManager(c services.ServiceManager) services.Discounts {
	s.serviceManager = c
	return s
}

func (s *discountsService) ServiceManager() services.ServiceManager {
	return s.serviceManager
}

func (s *discountsService) Get(ctx context.Context, req domains.DiscountRequest) (*domains.DiscountResponse, error) {
	ctx, span := s.ServiceManager().Spans().New(ctx)
	defer span.Finish()

	if err := s.ServiceManager().Validator().Validate(ctx, req); err != nil {
		span.Error(err)
		return nil, err
	}

	cid := ctxs.GetCidFromContext(ctx)
	ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{"x-cid": ptrs.StringValue(cid)}))

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	discountResponse, err := s.discountsClient.Get(ctx, &domainsgrpc.DiscountRequest{ProductId: req.ProductId, UserId: req.UserId})
	if err != nil {
		span.Error(err)
		return nil, err
	}

	resp := &domains.DiscountResponse{
		UserId:       discountResponse.UserId,
		ProductId:    discountResponse.ProductId,
		Percentage:   discountResponse.Percentage,
		ValueInCents: discountResponse.ValueInCents,
	}

	if err := s.ServiceManager().Validator().Validate(ctx, resp); err != nil {
		span.Error(err)
		return nil, err
	}

	return resp, nil
}
