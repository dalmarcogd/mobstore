package discounts

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dalmarcogd/mobstore/products/internal/discounts/discountsgrpc"
	"github.com/dalmarcogd/mobstore/products/internal/domains"
	"github.com/dalmarcogd/mobstore/products/internal/domains/domainsgrpc"
	"github.com/dalmarcogd/mobstore/products/internal/infra/ctxs"
	"github.com/dalmarcogd/mobstore/products/internal/infra/validator"
	"github.com/dalmarcogd/mobstore/products/internal/services"
)

type rpcMsg struct {
	msg proto.Message
}

func (r *rpcMsg) Matches(msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(m, r.msg)
}

func (r *rpcMsg) String() string {
	return fmt.Sprintf("is %s", r.msg)
}

func TestDiscountsService_Get(t *testing.T) {
	ctx := ctxs.ContextWithCid(context.Background(), "some-cid")
	serviceImpl := New()
	serviceImpl.WithAddress("http://localhost:9999")
	mockDiscountsClient := discountsgrpc.NewMockDiscountsClient(gomock.NewController(t))
	serviceImpl.discountsClient = mockDiscountsClient
	sm := services.New().WithValidator(validator.New()).WithDiscounts(serviceImpl)
	if err := sm.Init(); err != nil {
		t.Error("unexpected address")
	}

	mockDiscountsClient.EXPECT().
		Get(gomock.Any(), &rpcMsg{msg: &domainsgrpc.DiscountRequest{ProductId: "321", UserId: "123"}}).
		Return(&domainsgrpc.DiscountResponse{
			ProductId:    "321",
			UserId:       "123",
			Percentage:   10,
			ValueInCents: 300,
		}, nil).Times(1)

	req := domains.DiscountRequest{
		UserId:    "123",
		ProductId: "321",
	}
	res, err := serviceImpl.Get(ctx, req)
	if err != nil {
		t.Error(err)
	}

	if res.UserId != req.UserId {
		t.Errorf("expected %v and got %v", req.UserId, res.UserId)
	}
	if res.ProductId != req.ProductId {
		t.Errorf("expected %v and got %v", req.ProductId, res.ProductId)
	}
	if res.ValueInCents != 300 {
		t.Errorf("expected %v and got %v", 300, res.ValueInCents)
	}
	if res.Percentage != 10 {
		t.Errorf("expected %v and got %v", 10, res.Percentage)
	}

	errDeadlineExceeded := status.Error(codes.DeadlineExceeded, "timeout")
	mockDiscountsClient.EXPECT().
		Get(gomock.Any(), &rpcMsg{msg: &domainsgrpc.DiscountRequest{ProductId: "321", UserId: "123"}}).
		Return(nil, errDeadlineExceeded).Times(1)

	req = domains.DiscountRequest{
		UserId:    "123",
		ProductId: "321",
	}
	res, err = serviceImpl.Get(ctx, req)
	if !errors.Is(err, errDeadlineExceeded) {
		t.Errorf("expected (%v) got (%v)", errDeadlineExceeded, err)
	}

	errInternal := status.Error(codes.Internal, "internal")
	mockDiscountsClient.EXPECT().
		Get(gomock.Any(), &rpcMsg{msg: &domainsgrpc.DiscountRequest{ProductId: "321", UserId: "123"}}).
		Return(nil, errInternal).Times(1)

	req = domains.DiscountRequest{
		UserId:    "123",
		ProductId: "321",
	}
	res, err = serviceImpl.Get(ctx, req)
	if !errors.Is(err, errInternal) {
		t.Errorf("expected (%v) got (%v)", errInternal, err)
	}

	errInvalidArgument := status.Error(codes.InvalidArgument, "internal")
	mockDiscountsClient.EXPECT().
		Get(gomock.Any(), &rpcMsg{msg: &domainsgrpc.DiscountRequest{ProductId: "321", UserId: "123"}}).
		Return(nil, errInvalidArgument).Times(1)

	req = domains.DiscountRequest{
		UserId:    "123",
		ProductId: "321",
	}
	res, err = serviceImpl.Get(ctx, req)
	if !errors.Is(err, errInvalidArgument) {
		t.Errorf("expected (%v) got (%v)", errInvalidArgument, err)
	}
}
