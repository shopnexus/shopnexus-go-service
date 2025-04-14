package model

type Sale struct {
	ID               int64   `json:"id"` /* unique */
	Tag              *string `json:"tag"`
	ProductModelID   *int64  `json:"product_model_id"`
	BrandID          *int64  `json:"brand_id"`
	DateCreated      int64   `json:"date_created"`
	DateStarted      int64   `json:"date_started"`
	DateEnded        *int64  `json:"date_ended"`
	Quantity         int64   `json:"quantity"`
	Used             int64   `json:"used"`
	IsActive         bool    `json:"is_active"`
	DiscountPercent  *int32  `json:"discount_percent"`
	DiscountPrice    *int64  `json:"discount_price"`
	MaxDiscountPrice int64   `json:"max_discount_price"`
}

// Apply calculates the final sale discount (not the final price)
func (s Sale) Apply(price int64) int64 {
	var discount int64

	if s.DiscountPercent != nil {
		discount = price - price*int64(*s.DiscountPercent)/100
	}

	if s.DiscountPrice != nil {
		discount = *s.DiscountPrice
	}

	if s.MaxDiscountPrice > 0 && discount > s.MaxDiscountPrice {
		discount = s.MaxDiscountPrice
	}

	return discount
}
