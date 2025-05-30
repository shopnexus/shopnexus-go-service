package model

type (
	Status        string
	PaymentMethod string
	RefundMethod  string
)

const (
	StatusPending  Status = "PENDING"
	StatusSuccess  Status = "SUCCESS"
	StatusCanceled Status = "CANCELED"
	StatusFailed   Status = "FAILED"

	PaymentMethodCash  PaymentMethod = "CASH"
	PaymentMethodMomo  PaymentMethod = "MOMO"
	PaymentMethodVnpay PaymentMethod = "VNPAY"

	RefundMethodDropOff RefundMethod = "DROP_OFF"
	RefundMethodPickUp  RefundMethod = "PICK_UP"
)

type ProductOnPayment struct {
	ItemQuantityBase[int64]          // item_id is product_id
	ID                      int64    `json:"id"`          // unique
	SerialIDs               []string `json:"serial_ids"`  // List of serial IDs
	Price                   int64    `json:"price"`       // Single price, maybe have discount
	TotalPrice              int64    `json:"total_price"` // Total price, maybe have discount if reach certain quantity
}

type Payment struct {
	ID          int64         `json:"id"` /* unique */
	UserID      int64         `json:"user_id"`
	Method      PaymentMethod `json:"payment_method"`
	Address     string        `json:"address"`
	Total       int64         `json:"total"`
	Status      Status        `json:"status"`
	DateCreated int64         `json:"date_created"`

	Products []ProductOnPayment `json:"products"`
}
