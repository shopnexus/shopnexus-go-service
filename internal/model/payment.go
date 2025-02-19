package model

type (
	PaymentStatus string
	PaymentMethod string
)

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusSuccess   PaymentStatus = "SUCCESS"
	PaymentStatusCancelled PaymentStatus = "CANCELLED"
	PaymentStatusFailed    PaymentStatus = "FAILED"

	PaymentMethodCash     PaymentMethod = "CASH"
	PaymentMethodBankMomo PaymentMethod = "MOMO"
	PaymentMethodVnpay    PaymentMethod = "VNPAY"
)

type ProductOnPayment struct {
	ItemQuantityBase
	InvoiceID  []byte `json:"invoice_id"`
	Price      float64
	TotalPrice float64
}

type Payment struct {
	ID            []byte        `json:"id"` /* unique */
	UserID        []byte        `json:"user_id"`
	Address       string        `json:"address"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Total         float64       `json:"total"`
	Status        PaymentStatus `json:"status"`
	DateCreated   int64         `json:"date_created"`
}
