package model

type ProductModel struct {
	ID               int64  `json:"id"` /* unique */
	Type             int64  `json:"type"`
	BrandID          int64  `json:"brand_id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	ListPrice        int64  `json:"list_price"`
	DateManufactured int64  `json:"date_manufactured"`

	Resources []string `json:"resources"`
	Tags      []string `json:"tags"`
}
