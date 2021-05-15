package events

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"

	"github.com/dalmarcogd/mobstore/products/internal/domains"
	"github.com/dalmarcogd/mobstore/products/internal/infra/ctxs"
	"github.com/dalmarcogd/mobstore/products/internal/infra/iaws"
	"github.com/dalmarcogd/mobstore/products/internal/infra/validator"
	"github.com/dalmarcogd/mobstore/products/internal/services"
)

type (
	mockedSnsSuccess struct {
		iaws.SNSClient
	}
	mockedSnsFailure struct {
		iaws.SNSClient
	}
)

func (m mockedSnsSuccess) Publish(context.Context, *sns.PublishInput, ...func(*sns.Options)) (*sns.PublishOutput, error) {
	return &sns.PublishOutput{MessageId: aws.String("fsdfs")}, nil
}

func (m mockedSnsFailure) Publish(context.Context, *sns.PublishInput, ...func(*sns.Options)) (*sns.PublishOutput, error) {
	return nil, errors.New("some error")
}

func TestEventsService_Send(t *testing.T) {
	serviceImpl := New()
	sm := services.New().WithValidator(validator.New()).WithProductsEvents(serviceImpl)

	serviceImpl.client = mockedSnsSuccess{}
	if err := sm.Init(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	ctx := ctxs.ContextWithCid(context.Background(), "my cid")

	req := domains.NewProductCreatedEventRequest(ctx, "1", "123", "123123", 10)

	if err := sm.ProductsEvents().Send(ctx, req); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	serviceImpl.client = mockedSnsFailure{}
	if err := sm.Init(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := sm.ProductsEvents().Send(ctx, req); err == nil {
		t.Error("expected an error")
	}

	if err := sm.ProductsEvents().Send(ctx, nil); err == nil {
		t.Error("expected an error")
	}

	if err := serviceImpl.ServiceManager().Close(); err != nil {
		t.Error(err)
	}
}

func TestEventsService_SendBulk(t *testing.T) {
	serviceImpl := New()
	sm := services.New().WithValidator(validator.New()).WithProductsEvents(serviceImpl)

	serviceImpl.client = mockedSnsSuccess{}
	if err := sm.Init(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	ctx := ctxs.ContextWithCid(context.Background(), "my-cid")

	reqs := make([]domains.EventRequest, 0)

	for i := 0; i < 501; i++ {
		reqs = append(reqs, domains.NewProductCreatedEventRequest(ctx, "1", "123", "123123", 10))
	}

	if err := sm.ProductsEvents().SendBulk(ctx, reqs); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	serviceImpl.client = mockedSnsFailure{}
	if err := sm.Init(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := sm.ProductsEvents().SendBulk(ctx, reqs); err != nil {
		t.Error(err)
	}

	if err := sm.ProductsEvents().SendBulk(ctx, nil); err == nil {
		t.Error("expected an error")
	}

	if err := serviceImpl.ServiceManager().Close(); err != nil {
		t.Error(err)
	}
}
