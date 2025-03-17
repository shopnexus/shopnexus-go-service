package model

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
