package model

type Refund struct {
	ID          int64        `json:"id"` /* unique */
	PaymentID   int64        `json:"payment_id"`
	Method      RefundMethod `json:"method"`
	Status      Status       `json:"status"`
	Reason      string       `json:"reason"`
	Address     string       `json:"address"`
	DateCreated int64        `json:"date_created"`
	DateUpdated int64        `json:"date_updated"`
	Resources   []string     `json:"resources"`
}
