package payment

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"
	"shopnexus-go-service/config"
	"time"
)

type VnpayPlatform struct {
}

func (s *VnpayPlatform) CreateOrder(ctx context.Context, params CreateOrderParams) (url string, err error) {
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
	q.Add("vnp_ReturnUrl", "http://localhost:8080/payment/vnpay/callback")
	q.Add("vnp_ExpireDate", formatTime(time.Now().Add(30*time.Minute)))
	q.Add("vnp_TxnRef", fmt.Sprintf("%d", params.PaymentID))
	// q.Add("vnp_SecureHashType", "HMACSHA512")

	encodedQuery := q.Encode()
	secureHash := sign(encodedQuery, []byte(config.GetConfig().Vnpay.HashSecret))
	q.Add("vnp_SecureHash", secureHash)

	return req.URL.String() + "?" + encodedQuery + "&vnp_SecureHash=" + secureHash, nil

	// req.URL.RawQuery = q.Encode()

	// resp, err := httpClient.Do(req)
	// if err != nil {
	// 	return "", err
	// }

	// var data any
	// sonic.ConfigFastest.NewDecoder(resp.Body).Decode(&data)

	// defer resp.Body.Close()
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return "", err
	// }

	// return "", nil

}

// formatTime formats time to string in format yyyyMMddHHmmss
func formatTime(t time.Time) string {
	return t.Format("20060102150405")
}

func sign(message string, key []byte) string {
	sig := hmac.New(sha512.New, key)
	sig.Write([]byte(message))
	return hex.EncodeToString(sig.Sum(nil))
}
