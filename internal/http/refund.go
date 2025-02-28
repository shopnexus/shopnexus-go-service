package http

import (
	"net/http"
	"shopnexus-go-service/gen/pb"
	"shopnexus-go-service/internal/http/response"

	"github.com/bytedance/sonic"
	"github.com/go-chi/chi/v5"
)

type RefundHandler struct {
	client pb.RefundClient
}

func NewRefundHandler(client pb.RefundClient) http.Handler {
	h := &RefundHandler{client: client}

	r := chi.NewRouter()
	r.Post("/", h.CreateRefund)

	return r
}

func (h *RefundHandler) CreateRefund(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PaymentID int64   `json:"paymentId"`
		Method    string  `json:"method"`
		Reason    string  `json:"reason"`
		Address   *string `json:"address"`
	}

	if err := sonic.ConfigFastest.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	resp, err := h.client.Create(ctx, &pb.CreateRefundRequest{
		PaymentId: req.PaymentID,
		Method:    pb.RefundMethod(pb.RefundMethod_value[req.Method]),
		Reason:    req.Reason,
		Address:   req.Address,
	})
	if err != nil {
		response.FromError(w, err)
		return
	}

	response.FromDTO(w, http.StatusOK, struct {
		Id string `json:"id"`
	}{
		Id: resp.RefundId,
	})
}
