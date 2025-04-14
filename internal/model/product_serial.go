package model

type ProductSerial struct {
	SerialID    string `json:"serial_id"` /* unique */
	ProductID   int64  `json:"product_id"`
	IsSold      bool   `json:"is_sold"`
	IsActive    bool   `json:"is_active"`
	DateCreated int64  `json:"date_created"`
	DateUpdated int64  `json:"date_updated"`
}
