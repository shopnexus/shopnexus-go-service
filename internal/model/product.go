package model

import "encoding/json"

type Product struct {
	ID             int64           `json:"id"` /* unique */
	ProductModelID int64           `json:"product_model_id"`
	Quantity       int64           `json:"quantity,omitempty"`
	Sold           int64           `json:"sold,omitempty"`
	AddPrice       int64           `json:"add_price,omitempty"`
	IsActive       bool            `json:"is_active"`
	CanCombine     bool            `json:"can_combine"`
	Metadata       json.RawMessage `json:"metadata"`
	DateCreated    int64           `json:"date_created"`
	DateUpdated    int64           `json:"date_updated"`

	Resources []string `json:"resources"`
}
