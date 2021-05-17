package domains

import (
	"time"

	"github.com/dalmarcogd/mobstore/users/internal/infra/times"
)

type (
	UserCreateRequestV1 struct {
		FirstName string         `json:"first_name" validate:"required"`
		LastName  string         `json:"last_name" validate:"required"`
		BirthDate times.JsonTime `json:"birth_date" validate:"required"`
	}
	UserCreateResponseV1 struct {
		Id        string         `json:"id" validate:"required"`
		FirstName string         `json:"first_name" validate:"required"`
		LastName  string         `json:"last_name" validate:"required"`
		BirthDate times.JsonTime `json:"birth_date" validate:"required"`
	}
	UserUpdateRequestV1 struct {
		FirstName *string         `json:"first_name" validate:"omitempty,min=1"`
		LastName  *string         `json:"last_name" validate:"omitempty,min=1"`
		BirthDate *times.JsonTime `json:"birth_date" validate:"omitempty,min=1"`
	}
	UserUpdateResponseV1 struct {
		Id        string         `json:"id" validate:"required"`
		FirstName string         `json:"first_name" validate:"required"`
		LastName  string         `json:"last_name" validate:"required"`
		BirthDate times.JsonTime `json:"birth_date" validate:"required"`
	}
	UserGetResponseV1 struct {
		Id        string         `json:"id"`
		FirstName string         `json:"first_name" validate:"required"`
		LastName  string         `json:"last_name" validate:"required"`
		BirthDate times.JsonTime `json:"birth_date" validate:"required"`
	}
	UserListResponseV1 struct {
		Users []UserGetResponseV1 `json:"users"`
	}

	User struct {
		Id        *string    `projection:"id"`
		FirstName *string    `projection:"first_name" `
		LastName  *string    `projection:"last_name" `
		BirthDate *time.Time `projection:"birth_date" `
		CreatedAt time.Time  `projection:"created_at"`
		UpdatedAt time.Time  `projection:"updated_at"`
		DeletedAt *time.Time `projection:"deleted_at"`
	}
	UserCreate struct {
		FirstName string    `validate:"required"`
		LastName  string    `validate:"required"`
		BirthDate time.Time `validate:"required"`
	}
	UserUpdate struct {
		Id        string     `validate:"required"`
		FirstName *string    `validate:"omitempty,min=1"`
		LastName  *string    `validate:"omitempty,min=1"`
		BirthDate *time.Time `validate:"omitempty,required"`
	}
	UserSearch struct {
		Filter     UserFilter
		Projection UserProjection
	}

	UserFilter struct {
		Id        *string    `filter:"id"`
		FirstName *string    `filter:"first_name"`
		LastName  *string    `filter:"last_name"`
		BirthDate *time.Time `filter:"birth_Date"`
		DeletedAt *time.Time `filter:"deleted_at"`
	}

	UserProjection struct {
		Id        bool `projection:"id"`
		FirstName bool `projection:"first_name"`
		LastName  bool `projection:"last_name"`
		BirthDate bool `projection:"birth_date"`
		CreatedAt bool `projection:"created_at"`
		UpdatedAt bool `projection:"updated_at"`
		DeletedAt bool `projection:"deleted_at"`
	}
)

func (c *UserCreateRequestV1) UserCreate() UserCreate {
	return UserCreate{
		FirstName: c.FirstName,
		LastName:  c.LastName,
		BirthDate: time.Time(c.BirthDate),
	}
}

func (c *UserUpdateRequestV1) UserUpdate(userId string) UserUpdate {
	var t *time.Time
	jsonTime := c.BirthDate
	if jsonTime != nil {
		j := time.Time(*jsonTime)
		t = &j
	}
	return UserUpdate{
		Id:        userId,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		BirthDate: t,
	}
}
