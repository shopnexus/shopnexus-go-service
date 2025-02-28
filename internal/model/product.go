package model

type Brand struct {
	ID          int64  `json:"id"` /* unique */
	Name        string `json:"name"`
	Description string `json:"description"`

	Resources []string `json:"resources"`
}

type ProductModel struct {
	ID               int64  `json:"id"` /* unique */
	BrandID          int64  `json:"brand_id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	ListPrice        int64  `json:"list_price"`
	DateManufactured int64  `json:"date_manufactured"`

	Resources []string `json:"resources"`
	Tags      []string `json:"tags"`
}

// ProductIdentifier is a struct to identify a product, either by ID or SerialID
type ProductIdentifier struct {
	ID       *int64
	SerialID *string
}

type Product struct {
	ID             int64  `json:"id"`        /* unique */
	SerialID       string `json:"serial_id"` /* unique */
	ProductModelID int64  `json:"product_model_id"`
	DateCreated    int64  `json:"date_created"`
	DateUpdated    int64  `json:"date_updated"`
}

type Tag struct {
	Name        string `json:"name"` /* unique */
	Description string `json:"description"`
}

type Sale struct {
	ID              int64   `json:"id"` /* unique */
	Tag             *string `json:"tag"`
	ProductModelID  *int64  `json:"product_model_id"`
	DateStarted     int64   `json:"date_started"`
	DateEnded       *int64  `json:"date_ended"`
	Quantity        int64   `json:"quantity"`
	Used            int64   `json:"used"`
	IsActive        bool    `json:"is_active"`
	DiscountPercent *int64  `json:"discount_percent"`
	DiscountPrice   *int64  `json:"discount_price"`
}
