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

type Invoice struct {
	ID            []byte        `json:"id"` /* unique */
	UserID        []byte        `json:"user_id"`
	Address       string        `json:"address"`
	Total         float64       `json:"total"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	DateCreated   int64         `json:"date_created"`

	Items []ItemQuantity `json:"items"`
}

type ProductOnInvoice struct {
	ItemQuantityBase
	InvoiceID  []byte `json:"invoice_id"`
	Price      float64
	TotalPrice float64
}

type Payment struct {
	ID            []byte        `json:"id"` /* unique */
	Status        PaymentStatus `json:"status"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	InvoiceID     []byte        `json:"invoice_id"`
	DateCreated   int64         `json:"date_created"`
	DateExpired   int64         `json:"date_expired"`

	Invoice Invoice `json:"invoice"`
}
