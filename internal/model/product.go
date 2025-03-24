package model

// ProductIdentifier is a struct to identify a product, either by ID or SerialID
type ProductIdentifier struct {
	ID       *int64
	SerialID *string
}

type Product[T any] struct {
	ID             int64  `json:"id"`        /* unique */
	SerialID       string `json:"serial_id"` /* unique */
	Quantity       int64  `json:"quantity,omitempty"`
	Sold           int64  `json:"sold,omitempty"`
	AddPrice       int64  `json:"add_price,omitempty"`
	IsActive       bool   `json:"is_active"`
	ProductModelID int64  `json:"product_model_id"`
	Metadata       T      `json:"metadata"`
	DateCreated    int64  `json:"date_created"`
	DateUpdated    int64  `json:"date_updated"`
}

type ProductMetadataShoe struct {
	Size  int    `json:"size"`
	Color string `json:"color"` // hex color code
}
