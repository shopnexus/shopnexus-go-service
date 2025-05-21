package vnpay

import (
	"context"
	"fmt"
	"net/http"
	"shopnexus-go-service/config"
	"time"
)

type Service struct {
}

type ServiceInterface interface {
	CreateOrder(ctx context.Context, params CreateOrderParams) (url string, err error)
	VerifyPayment(ctx context.Context, query IPNObject) error
}

type CreateOrderParams struct {
	PaymentID int64
	Amount    int64
	Info      string
}

func (s *Service) CreateOrder(ctx context.Context, params CreateOrderParams) (url string, err error) {
	// httpClient := &http.Client{}
	req, err := http.NewRequest("GET", "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html", nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("vnp_Version", "2.1.0")
	q.Add("vnp_Command", "pay")
	q.Add("vnp_TmnCode", config.GetConfig().Vnpay.TmnCode)
	q.Add("vnp_Amount", fmt.Sprintf("%d", params.Amount*100))
	// q.Add("vnp_BankCode", string(BankCodeVNPAYQR))
	q.Add("vnp_CreateDate", formatTime(time.Now()))
	q.Add("vnp_CurrCode", "VND")
	q.Add("vnp_IpAddr", "192.168.1.1")
	q.Add("vnp_Locale", "vn")
	q.Add("vnp_OrderInfo", params.Info)
	q.Add("vnp_OrderType", "billpayment")
	q.Add("vnp_ReturnUrl", "/payment-done")
	q.Add("vnp_ExpireDate", formatTime(time.Now().Add(30*time.Minute)))
	q.Add("vnp_TxnRef", fmt.Sprintf("%d", params.PaymentID))
	// q.Add("vnp_SecureHashType", "HMACSHA512")

	encodedQuery := q.Encode()
	secureHash := sign(encodedQuery, []byte(config.GetConfig().Vnpay.HashSecret))
	q.Add("vnp_SecureHash", secureHash)

	return req.URL.String() + "?" + encodedQuery + "&vnp_SecureHash=" + secureHash, nil
}

type IPNObject struct {
	MerchantID              string `json:"merchant_id"`
	MerchantName            string `json:"merchant_name"`
	MerchantTransactionRef  string `json:"merchant_transaction_ref"`
	TransactionInfo         string `json:"transaction_info"`
	Amount                  string `json:"amount"`
	CurrentCode             string `json:"current_code"`
	TransactionResponseCode string `json:"transaction_response_code"`
	Message                 string `json:"message"`
	TransactionNumber       string `json:"transaction_number"`
	Bank                    string `json:"bank"`
}

func (s *Service) VerifyPayment(ctx context.Context, query IPNObject) error {
	return nil
}
