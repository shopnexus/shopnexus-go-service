package http

import (
	"net/http"
	"shopnexus-go-service/gen/pb"
	"shopnexus-go-service/internal/http/response"

	"github.com/bytedance/sonic"
	"github.com/go-chi/chi/v5"
)

type PaymentHandler struct {
	client pb.PaymentClient
}

func NewPaymentHandler(client pb.PaymentClient) http.Handler {
	h := &PaymentHandler{client: client}

	r := chi.NewRouter()
	r.Post("/", h.CreatePayment)

	return r
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Address       string `json:"address"`
		PaymentMethod string `json:"payment_method"`
	}

	if err := sonic.ConfigFastest.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	resp, err := h.client.Create(ctx, &pb.CreatePaymentRequest{
		Address:       req.Address,
		PaymentMethod: pb.PaymentMethod(pb.PaymentMethod_value[req.PaymentMethod]),
	})
	if err != nil {
		response.FromError(w, err)
		return
	}

	response.FromDTO(w, http.StatusOK, struct {
		Id  string `json:"id"`
		Url string `json:"url"`
	}{
		Id:  resp.PaymentId,
		Url: resp.PaymentUrl,
	})
}
