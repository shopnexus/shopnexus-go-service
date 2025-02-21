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
	ItemQuantityBase[string]
	Price      int64 `json:"price"`
	TotalPrice int64 `json:"total_price"`
}

type Payment struct {
	ID            int64         `json:"id"` /* unique */
	UserID        int64         `json:"user_id"`
	Address       string        `json:"address"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Total         int64         `json:"total"`
	Status        PaymentStatus `json:"status"`
	DateCreated   int64         `json:"date_created"`

	Products []ProductOnPayment `json:"products"`
}
