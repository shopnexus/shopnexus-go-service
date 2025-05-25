package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) IPNVNPAY(c echo.Context) error {
	var query struct {
		VnpAmount            string `query:"vnp_Amount"`
		VnpBankCode          string `query:"vnp_BankCode"`
		VnpCardType          string `query:"vnp_CardType"`
		VnpOrderInfo         string `query:"vnp_OrderInfo"`
		VnpPayDate           string `query:"vnp_PayDate"`
		VnpResponseCode      string `query:"vnp_ResponseCode"`
		VnpSecureHash        string `query:"vnp_SecureHash"`
		VnpTmnCode           string `query:"vnp_TmnCode"`
		VnpTransactionNo     string `query:"vnp_TransactionNo"`
		VnpTransactionStatus string `query:"vnp_TransactionStatus"`
		VnpTxnRef            string `query:"vnp_TxnRef"`
	}

	if err := c.Bind(&query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request parameters")
	}

	c.JSON(http.StatusOK, query)
	return nil
}
