package vnpay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ClientImpl struct {
	tmnCode    string
	hashSecret string
}

type Client interface {
	CreateOrder(ctx context.Context, params CreateOrderParams) (url string, err error)
	VerifyPayment(ctx context.Context, rawIpn json.RawMessage) error
}

type ClientOptions struct {
	TmnCode    string
	HashSecret string
}

func NewClient(cfg ClientOptions) Client {
	return &ClientImpl{
		tmnCode:    cfg.TmnCode,
		hashSecret: cfg.HashSecret,
	}
}

type CreateOrderParams struct {
	PaymentID int64
	Amount    int64
	Info      string
}

func (c *ClientImpl) CreateOrder(ctx context.Context, params CreateOrderParams) (url string, err error) {
	// httpClient := &http.Client{}
	req, err := http.NewRequest("GET", "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html", nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("vnp_Version", "2.1.0")
	q.Add("vnp_Command", "pay")
	q.Add("vnp_TmnCode", c.tmnCode)
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
	secureHash := sign(encodedQuery, []byte(c.hashSecret))
	q.Add("vnp_SecureHash", secureHash)

	return req.URL.String() + "?" + encodedQuery + "&vnp_SecureHash=" + secureHash, nil
}

type IPNReturn struct {
	VnpAmount            string `json:"vnp_Amount"`
	VnpBankCode          string `json:"vnp_BankCode"`
	VnpCardType          string `json:"vnp_CardType"`
	VnpOrderInfo         string `json:"vnp_OrderInfo"`
	VnpPayDate           string `json:"vnp_PayDate"`
	VnpResponseCode      string `json:"vnp_ResponseCode"`
	VnpSecureHash        string `json:"vnp_SecureHash"`
	VnpTmnCode           string `json:"vnp_TmnCode"`
	VnpTransactionNo     string `json:"vnp_TransactionNo"`
	VnpTransactionStatus string `json:"vnp_TransactionStatus"`
	VnpTxnRef            string `json:"vnp_TxnRef"`
}

func (c *ClientImpl) VerifyPayment(ctx context.Context, rawIpn json.RawMessage) error {
	var ipn map[string]any
	if err := json.Unmarshal(rawIpn, &ipn); err != nil {
		return fmt.Errorf("failed to unmarshal IPN: %w", err)
	}

	expectedHash, ok := ipn["vnp_SecureHash"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid vnp_SecureHash")
	}

	// Remove the secure hash from the IPN data
	delete(ipn, "vnp_SecureHash")

	hashData := buildSortedQuery(ipn)
	hash := sign(c.hashSecret, []byte(hashData))

	if hash != expectedHash {
		return fmt.Errorf("hash mismatch: expected %s, got %s", expectedHash, hash)
	}

	return nil
}
