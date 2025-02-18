package model

type Brand struct {
	ID          []byte `json:"id"` /* unique */
	Name        string `json:"name"`
	Description string `json:"description"`

	Images []string `json:"images"`
}

type ProductModel struct {
	ID          []byte  `json:"id"` /* unique */
	BrandID     []byte  `json:"brand_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ListPrice   float64 `json:"list_price"`

	Images []string `json:"images"`
	Tags   []string `json:"tags"`
}

type Product struct {
	SerialID       []byte `json:"serial_id"` /* unique */
	ProductModelID []byte
	EntryDate      int64
	UpdateDate     int64
}

type Tag struct {
	Name        string `json:"name"` /* unique */
	Description string `json:"description"`
}

type Sale struct {
	ID              []byte   `json:"id"` /* unique */
	Tag             *string  `json:"tag"`
	ProductModelID  *[]byte  `json:"product_model_id"`
	StartDate       int64    `json:"start_date"`
	EndDate         *int64   `json:"end_date"`
	Quantity        int64    `json:"quantity"`
	Used            int64    `json:"used"`
	IsActive        bool     `json:"is_active"`
	DiscountPercent *float64 `json:"discount_percent"`
	DiscountPrice   *float64 `json:"discount_price"`
}
