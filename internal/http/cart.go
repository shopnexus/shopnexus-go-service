package http

import (
	"net/http"
	"shopnexus-go-service/gen/pb"
	"shopnexus-go-service/internal/http/response"
	"shopnexus-go-service/internal/util"

	"github.com/bytedance/sonic"
	"github.com/go-chi/chi/v5"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CartHandler struct {
	client pb.CartClient
}

func NewCartHandler(client pb.CartClient) http.Handler {
	h := &CartHandler{client: client}

	r := chi.NewRouter()
	r.Post("/items", h.AddCartItem)
	r.Get("/", h.GetCart)
	r.Put("/items", h.UpdateCartItem)
	r.Delete("/", h.ClearCart)

	return r
}

func (h *CartHandler) AddCartItem(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProductModelID int64 `json:"product_model_id"`
		Quantity       int64 `json:"quantity"`
	}

	if err := sonic.ConfigFastest.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	// Forward the auth token to gRPC metadata
	ctx := r.Context()

	resp, err := h.client.AddCartItem(ctx, &pb.AddCartItemRequest{
		ProductModelId: req.ProductModelID,
		Quantity:       req.Quantity,
	})
	if err != nil {
		response.FromError(w, err)
		return
	}

	response.FromDTO(w, http.StatusOK, struct {
		Quantity int64 `json:"quantity"`
	}{
		Quantity: resp.Quantity,
	})
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := h.client.GetCart(ctx, &emptypb.Empty{})
	if err != nil {
		response.FromError(w, err)
		return
	}

	response.FromDTO(w, http.StatusOK, struct {
		Items []*pb.ItemQuantity `json:"items"`
	}{
		Items: util.NonEmptySlice(resp.Items),
	})
}

func (h *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProductModelID int64 `json:"product_model_id"`
		Quantity       int64 `json:"quantity"`
	}

	if err := sonic.ConfigFastest.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	_, err := h.client.UpdateCartItem(ctx, &pb.UpdateCartItemRequest{
		ProductModelId: req.ProductModelID,
		Quantity:       req.Quantity,
	})
	if err != nil {
		response.FromError(w, err)
		return
	}

	response.FromMessage(w, http.StatusOK, "Cart item updated successfully")
}

func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, err := h.client.ClearCart(ctx, &emptypb.Empty{})
	if err != nil {
		response.FromError(w, err)
		return
	}

	response.FromMessage(w, http.StatusOK, "Cart cleared successfully")
}
