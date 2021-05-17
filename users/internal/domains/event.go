package domains

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/dalmarcogd/mobstore/users/internal/infra/ctxs"
	"github.com/dalmarcogd/mobstore/users/internal/infra/ptrs"
	"github.com/dalmarcogd/mobstore/users/internal/infra/times"
)

var (
	UsersTypeEvent TypeEvent = "users"
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
	UserEventRequest struct {
		event
		UserId    string         `json:"user_id" validate:"required"`
		FirstName string         `json:"first_name" validate:"required"`
		LastName  string         `json:"last_name" validate:"required"`
		BirthDate times.JsonTime `json:"birth_date" validate:"required"`
	}
)

func (p UserEventRequest) EventId() string {
	return p.event.EventId
}

func (p UserEventRequest) EventType() string {
	return string(p.event.EventType)
}

func (p UserEventRequest) Operation() string {
	return p.event.Operation
}

func NewUserCreatedEventRequest(ctx context.Context, userId, firstName, lastName string, birthDate time.Time) UserEventRequest {
	return UserEventRequest{
		event: event{
			EventId:   uuid.NewString(),
			EventType: UsersTypeEvent,
			Operation: "create",
			Cid:       ptrs.StringValue(ctxs.GetCidFromContext(ctx)),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		UserId:    userId,
		FirstName: firstName,
		LastName:  lastName,
		BirthDate: times.JsonTime(birthDate),
	}
}

func NewUserUpdatedEventRequest(ctx context.Context, userId, firstName, lastName string, birthDate time.Time) UserEventRequest {
	return UserEventRequest{
		event: event{
			EventId:   uuid.NewString(),
			EventType: UsersTypeEvent,
			Operation: "update",
			Cid:       ptrs.StringValue(ctxs.GetCidFromContext(ctx)),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		UserId:    userId,
		FirstName: firstName,
		LastName:  lastName,
		BirthDate: times.JsonTime(birthDate),
	}
}

func NewUserDeletedEventRequest(ctx context.Context, userId, firstName, lastName string, birthDate time.Time) UserEventRequest {
	return UserEventRequest{
		event: event{
			EventId:   uuid.NewString(),
			EventType: UsersTypeEvent,
			Operation: "delete",
			Cid:       ptrs.StringValue(ctxs.GetCidFromContext(ctx)),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		UserId:    userId,
		FirstName: firstName,
		LastName:  lastName,
		BirthDate: times.JsonTime(birthDate),
	}
}
