package model

type (
	Status        string
	PaymentMethod string
	RefundMethod  string
)

const (
	StatusPending   Status = "PENDING"
	StatusSuccess   Status = "SUCCESS"
	StatusCancelled Status = "CANCELLED"
	StatusFailed    Status = "FAILED"

	PaymentMethodCash  PaymentMethod = "CASH"
	PaymentMethodMomo  PaymentMethod = "MOMO"
	PaymentMethodVnpay PaymentMethod = "VNPAY"

	RefundMethodDropOff RefundMethod = "DROP_OFF"
	RefundMethodPickUp  RefundMethod = "PICK_UP"
)

type ProductOnPayment struct {
	ItemQuantityBase[string]
	Price      int64 `json:"price"`       // Single price, maybe have discount
	TotalPrice int64 `json:"total_price"` // Total price, maybe have discount if reach certain quantity
}

type Payment struct {
	ID            int64         `json:"id"` /* unique */
	UserID        int64         `json:"user_id"`
	Address       string        `json:"address"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Total         int64         `json:"total"`
	Status        Status        `json:"status"`
	DateCreated   int64         `json:"date_created"`

	Products []ProductOnPayment `json:"products"`
}

type Refund struct {
	ID          int64        `json:"id"` /* unique */
	PaymentID   int64        `json:"payment_id"`
	Method      RefundMethod `json:"method"`
	Status      Status       `json:"status"`
	Reason      string       `json:"reason"`
	Address     *string      `json:"address"`
	DateCreated int64        `json:"date_created"`
	DateUpdated int64        `json:"date_updated"`

	Resources []string `json:"resources"`
}
