package domains

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/dalmarcogd/mobstore/products/internal/infra/ctxs"
	"github.com/dalmarcogd/mobstore/products/internal/infra/ptrs"
)

var (
	ProductsTypeEvent TypeEvent = "products"
)

type (
	TypeEvent string

	event struct {
		EventId   string    `json:"event_id"`
		EventType TypeEvent `json:"event_type" validate:"required"`
		Operation string    `json:"operation" validate:"required"`
		Cid       string    `json:"cid" validate:"required"`
		Timestamp string    `json:"timestamp" validate:"required"`
	}

	EventRequest interface {
		EventId() string
		EventType() string
		Operation() string
	}
	ProductCreatedUpdatedEventRequest struct {
		event
		ProductId    string `json:"product_id" validate:"required"`
		Title        string `json:"title" validate:"required"`
		Description  string `json:"description" validate:"required"`
		PriceInCents int64  `json:"value_in_cents" validate:"required"`
	}
)

func (p ProductCreatedUpdatedEventRequest) EventId() string {
	return p.event.EventId
}

func (p ProductCreatedUpdatedEventRequest) EventType() string {
	return string(p.event.EventType)
}

func (p ProductCreatedUpdatedEventRequest) Operation() string {
	return p.event.Operation
}

func NewProductCreatedEventRequest(ctx context.Context, productId, title, description string, priceInCents int64) ProductCreatedUpdatedEventRequest {
	return ProductCreatedUpdatedEventRequest{
		event: event{
			EventId:   uuid.NewString(),
			EventType: ProductsTypeEvent,
			Operation: "create",
			Cid:       ptrs.StringValue(ctxs.GetCidFromContext(ctx)),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		ProductId:    productId,
		Title:        title,
		Description:  description,
		PriceInCents: priceInCents,
	}
}

func NewProductUpdatedEventRequest(ctx context.Context, productId, title, description string, priceInCents int64) ProductCreatedUpdatedEventRequest {
	return ProductCreatedUpdatedEventRequest{
		event: event{
			EventId:   uuid.NewString(),
			EventType: ProductsTypeEvent,
			Operation: "update",
			Cid:       ptrs.StringValue(ctxs.GetCidFromContext(ctx)),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		ProductId:    productId,
		Title:        title,
		Description:  description,
		PriceInCents: priceInCents,
	}
}

func NewProductDeletedEventRequest(ctx context.Context, productId, title, description string, priceInCents int64) ProductCreatedUpdatedEventRequest {
	return ProductCreatedUpdatedEventRequest{
		event: event{
			EventId:   uuid.NewString(),
			EventType: ProductsTypeEvent,
			Operation: "delete",
			Cid:       ptrs.StringValue(ctxs.GetCidFromContext(ctx)),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		ProductId:    productId,
		Title:        title,
		Description:  description,
		PriceInCents: priceInCents,
	}
}
