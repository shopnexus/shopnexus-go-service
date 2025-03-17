package model

type Sale struct {
	ID              int64   `json:"id"` /* unique */
	Tag             *string `json:"tag"`
	ProductModelID  *int64  `json:"product_model_id"`
	BrandID         *int64  `json:"brand_id"`
	DateStarted     int64   `json:"date_started"`
	DateEnded       *int64  `json:"date_ended"`
	Quantity        int64   `json:"quantity"`
	Used            int64   `json:"used"`
	IsActive        bool    `json:"is_active"`
	DiscountPercent *int32  `json:"discount_percent"`
	DiscountPrice   *int64  `json:"discount_price"`
}
