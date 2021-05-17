package domains

type (
	DiscountRequest struct {
		UserId    string `json:"user_id" validate:"required"`
		ProductId string `json:"product_id" validate:"required"`
	}
	DiscountResponse struct {
		UserId       string  `json:"user_id" validate:"required"`
		ProductId    string  `json:"product_id" validate:"required"`
		Percentage   float32 `json:"percentage" `
		ValueInCents float32 `json:"value_in_cents" `
	}
)
