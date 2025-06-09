package model

import "encoding/json"

type Product struct {
	ID              int64           `json:"id"` /* unique */
	ProductModelID  int64           `json:"product_model_id"`
	CurrentStock    int64           `json:"current_stock,omitempty"`
	Sold            int64           `json:"sold,omitempty"`
	AdditionalPrice int64           `json:"additional_price,omitempty"`
	IsActive        bool            `json:"is_active"`
	CanCombine      bool            `json:"can_combine"`
	Metadata        json.RawMessage `json:"metadata"`
	DateCreated     int64           `json:"date_created"`

	Resources []string `json:"resources"`
}
