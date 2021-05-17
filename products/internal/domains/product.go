package domains

import "time"

type (
	ProductCreateRequestV1 struct {
		PriceInCents int64  `json:"price_in_cents" validate:"required"`
		Title        string `json:"title" validate:"required,min=1"`
		Description  string `json:"description" validate:"required,min=1"`
	}
	ProductCreateResponseV1 struct {
		Id           string `json:"id" validate:"required"`
		PriceInCents int64  `json:"price_in_cents" validate:"required"`
		Title        string `json:"title" validate:"required,min=1"`
		Description  string `json:"description" validate:"required,min=1"`
	}
	ProductUpdateRequestV1 struct {
		PriceInCents *int64  `json:"price_in_cents"`
		Title        *string `json:"title" validate:"omitempty,min=1"`
		Description  *string `json:"description" validate:"omitempty,min=1"`
	}
	ProductUpdateResponseV1 struct {
		Id           string `json:"id" validate:"required"`
		PriceInCents int64  `json:"price_in_cents" validate:"required"`
		Title        string `json:"title" validate:"required,min=1"`
		Description  string `json:"description" validate:"omitempty,min=1"`
	}
	ProductGetResponseV1 struct {
		Id           string           `json:"id"`
		PriceInCents int64            `json:"price_in_cents"`
		Title        string           `json:"title"`
		Description  *string          `json:"description"`
		Discount     *ProductDiscount `json:"discount"`
	}
	ProductListResponseV1 struct {
		Products []ProductGetResponseV1 `json:"products"`
	}

	Product struct {
		Id           *string    `projection:"id"`
		PriceInCents *int64     `projection:"price_in_cents"`
		Title        *string    `projection:"title" `
		Description  *string    `projection:"description"`
		CreatedAt    time.Time  `projection:"created_at"`
		UpdatedAt    time.Time  `projection:"updated_at"`
		DeletedAt    *time.Time `projection:"deleted_at"`
		Discount     *ProductDiscount
	}
	ProductCreate struct {
		PriceInCents int64  `validate:"required"`
		Title        string `validate:"required"`
		Description  string `validate:"required"`
	}
	ProductUpdate struct {
		Id           string  `validate:"required"`
		PriceInCents *int64  `validate:"omitempty,required"`
		Title        *string `validate:"omitempty,min=1"`
		Description  *string `validate:"omitempty,min=1"`
	}
	ProductSearch struct {
		Filter     ProductFilter
		Projection ProductProjection
	}
	ProductFilter struct {
		Id           *string    `filter:"id"`
		PriceInCents *int64     `filter:"price_in_cents"`
		Title        *string    `filter:"title"`
		Description  *string    `filter:"description"`
		DeletedAt    *time.Time `filter:"deleted_at"`
	}
	ProductProjection struct {
		Id           bool `projection:"id"`
		PriceInCents bool `projection:"price_in_cents"`
		Title        bool `projection:"title"`
		Description  bool `projection:"description"`
		CreatedAt    bool `projection:"created_at"`
		UpdatedAt    bool `projection:"updated_at"`
		DeletedAt    bool `projection:"deleted_at"`
	}

	ProductDiscount struct {
		Percentage   float64 `json:"percentage"`
		ValueInCents float64   `json:"value_in_cents"`
	}
)

func (c *ProductCreateRequestV1) ProductCreate() ProductCreate {
	return ProductCreate{
		PriceInCents: c.PriceInCents,
		Title:        c.Title,
		Description:  c.Description,
	}
}

func (c *ProductUpdateRequestV1) ProductUpdate(productId string) ProductUpdate {
	return ProductUpdate{
		Id:           productId,
		PriceInCents: c.PriceInCents,
		Title:        c.Title,
		Description:  c.Description,
	}
}
