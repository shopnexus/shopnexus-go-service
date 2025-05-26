package http

import (
	"fmt"
	"net/http"
	"shopnexus-go-service/config"
	"shopnexus-go-service/internal/logger"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service/payment"
	"shopnexus-go-service/internal/utils/ptr"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) IPNVNPAY(c echo.Context) error {
	var query map[string]any

	if err := c.Bind(&query); err != nil {
		logger.Log.Error("Failed to bind query parameters" + "error" + err.Error())
		c.NoContent(http.StatusBadRequest)
		return nil
	}

	fmt.Println("Received VNPAY IPN:", query)

	// Verify the checksum hash
	if err := h.services.Vnpay.VerifyPayment(c.Request().Context(), query); err != nil {
		c.NoContent(http.StatusBadRequest)
		return nil
	}

	// After verifying the checksum, always return 200 OK
	defer c.NoContent(http.StatusOK)

	tmnCode, ok := query["vnp_TmnCode"].(string)
	if !ok {
		logger.Log.Error("TmnCode not found in query parameters")
		c.NoContent(http.StatusBadRequest)
		return nil
	}

	txnStatus, ok := query["vnp_TransactionStatus"].(string)
	if !ok {
		logger.Log.Error("Transaction status not found in query parameters")
		c.NoContent(http.StatusBadRequest)
		return nil
	}

	txnRef, ok := query["vnp_TxnRef"].(string)
	if !ok {
		logger.Log.Error("Transaction reference not found in query parameters")
		c.NoContent(http.StatusBadRequest)
		return nil
	}

	// Verify the TmnCode
	if tmnCode != config.GetConfig().Vnpay.TmnCode {
		logger.Log.Error("Invalid TmnCode" + "tmn_code" + tmnCode)
		c.NoContent(http.StatusBadRequest)
		return nil
	}

	// Verify the transaction status
	if txnStatus != "00" {
		logger.Log.Error("Transaction failed" + "transaction_status" + txnStatus)
		return nil
	}

	paymentID, err := strconv.ParseInt(txnRef, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payment ID")
	}

	if err = h.services.Payment.UpdatePayment(c.Request().Context(), payment.UpdatePaymentParams{
		ID:     paymentID,
		Role:   model.RoleAdmin,
		Status: ptr.ToPtr(model.StatusSuccess),
	}); err != nil {
		logger.Log.Error("Failed to update payment status" + "error" + err.Error())
	}

	logger.Log.Info("Payment successfully updated" + "payment_id" + strconv.FormatInt(paymentID, 10))

	return nil
}
