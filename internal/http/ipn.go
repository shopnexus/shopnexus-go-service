package http

import (
	"fmt"
	"net/http"
	"shopnexus-go-service/gen/pb"
	"shopnexus-go-service/internal/http/response"

	"github.com/go-chi/chi/v5"
)

type IPNHandler struct {
	client pb.PaymentClient
}

func NewIPNHandler(client pb.PaymentClient) http.Handler {
	h := &IPNHandler{client: client}

	r := chi.NewRouter()
	r.Get("/", h.IPN)

	return r
}

func (h *IPNHandler) IPN(w http.ResponseWriter, r *http.Request) {
	var query struct {
		VnpAmount            string `schema:"vnp_Amount"`
		VnpBankCode          string `schema:"vnp_BankCode"`
		VnpBankTranNo        string `schema:"vnp_BankTranNo"`
		VnpCardType          string `schema:"vnp_CardType"`
		VnpOrderInfo         string `schema:"vnp_OrderInfo"`
		VnpPayDate           string `schema:"vnp_PayDate"`
		VnpResponseCode      string `schema:"vnp_ResponseCode"`
		VnpSecureHash        string `schema:"vnp_SecureHash"`
		VnpTmnCode           string `schema:"vnp_TmnCode"`
		VnpTransactionNo     string `schema:"vnp_TransactionNo"`
		VnpTransactionStatus string `schema:"vnp_TransactionStatus"`
		VnpTxnRef            string `schema:"vnp_TxnRef"`
	}

	if err := decode.Decode(&query, r.URL.Query()); err != nil {
		response.FromMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(query)

	response.FromMessage(w, http.StatusOK, "ok")
}
