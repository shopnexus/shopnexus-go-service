package http

import (
	"net/http"
	"shopnexus-go-service/gen/pb"
	"shopnexus-go-service/internal/http/response"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/util"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	client pb.ProductClient
}

func NewProductHandler(client pb.ProductClient) http.Handler {
	h := &ProductHandler{client: client}

	r := chi.NewRouter()

	// Product Model routes
	r.Get("/models/{id}", h.GetProductModel)
	r.Get("/models", h.ListProductModels)
	r.Delete("/models/{id}", h.DeleteProductModel)

	// Product routes
	r.Get("/{id}", h.GetProduct)
	r.Get("/", h.ListProducts)
	r.Post("/", h.CreateProduct)
	r.Put("/{id}", h.UpdateProduct)
	r.Delete("/{id}", h.DeleteProduct)

	return r
}

func (h *ProductHandler) GetProductModel(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.FromMessage(w, http.StatusBadRequest, "invalid id")
		return
	}

	ctx := r.Context()

	resp, err := h.client.GetProductModel(ctx, &pb.GetProductModelRequest{Id: id})
	if err != nil {
		response.FromMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.FromDTO(w, http.StatusOK, struct {
		Id               int64  `json:"id"`
		BrandId          int64  `json:"brandId"`
		Name             string `json:"name"`
		Description      string `json:"description"`
		ListPrice        int64  `json:"listPrice"`
		DateManufactured int64  `json:"dateManufactured"`
	}{
		Id:               resp.Id,
		BrandId:          resp.BrandId,
		Name:             resp.Name,
		Description:      resp.Description,
		ListPrice:        resp.ListPrice,
		DateManufactured: resp.DateManufactured,
	})
}

func (h *ProductHandler) ListProductModels(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Page                 int32   `schema:"page"`
		Limit                int32   `schema:"limit"`
		BrandID              *int64  `schema:"brand_id"`
		Name                 *string `schema:"name"`
		Description          *string `schema:"description"`
		ListPriceFrom        *int64  `schema:"list_price_from"`
		ListPriceTo          *int64  `schema:"list_price_to"`
		DateManufacturedFrom *int64  `schema:"date_manufactured_from"`
		DateManufacturedTo   *int64  `schema:"date_manufactured_to"`
	}
	if err := decode.Decode(&req, r.URL.Query()); err != nil {
		response.FromMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	resp, err := h.client.ListProductModels(ctx, &pb.ListProductModelsRequest{
		Pagination: &pb.PaginationRequest{
			Page:  util.Max(req.Page, 1),
			Limit: util.Max(req.Limit, 1),
		},
		BrandId:              req.BrandID,
		Name:                 req.Name,
		Description:          req.Description,
		ListPriceFrom:        req.ListPriceFrom,
		ListPriceTo:          req.ListPriceTo,
		DateManufacturedFrom: req.DateManufacturedFrom,
		DateManufacturedTo:   req.DateManufacturedTo,
	})
	if err != nil {
		response.FromMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.FromPagination(w, http.StatusOK, model.PaginateResult[*pb.ProductModelEntity]{
		Data:       util.NonEmptySlice(resp.ProductModels),
		Page:       resp.Pagination.Page,
		Limit:      resp.Pagination.Limit,
		Total:      resp.Pagination.Total,
		NextPage:   resp.Pagination.NextPage,
		NextCursor: resp.Pagination.NextCursor,
	})
}

func (h *ProductHandler) DeleteProductModel(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.FromMessage(w, http.StatusBadRequest, "invalid id")
		return
	}

	ctx := r.Context()

	_, err = h.client.DeleteProductModel(ctx, &pb.DeleteProductModelRequest{Id: id})
	if err != nil {
		response.FromMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.FromMessage(w, http.StatusOK, "Product model deleted successfully")
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.FromMessage(w, http.StatusBadRequest, "invalid id")
		return
	}

	ctx := r.Context()

	resp, err := h.client.GetProduct(ctx, &pb.GetProductRequest{
		Id: &id,
	})
	if err != nil {
		response.FromMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.FromDTO(w, http.StatusOK, struct {
		Id             int64  `json:"id"`
		SerialId       string `json:"serial_id"`
		ProductModelId int64  `json:"product_model_id"`
		DateCreated    int64  `json:"date_created"`
		DateUpdated    int64  `json:"date_updated"`
	}{
		Id:             resp.Id,
		SerialId:       resp.SerialId,
		ProductModelId: resp.ProductModelId,
		DateCreated:    resp.DateCreated,
		DateUpdated:    resp.DateUpdated,
	})
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Page            int32  `schema:"page"`
		Limit           int32  `schema:"limit"`
		ProductModelID  *int64 `schema:"product_model_id"`
		DateCreatedFrom *int64 `schema:"date_created_from"`
		DateCreatedTo   *int64 `schema:"date_created_to"`
	}
	if err := decode.Decode(&req, r.URL.Query()); err != nil {
		response.FromMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	resp, err := h.client.ListProducts(ctx, &pb.ListProductsRequest{
		Pagination: &pb.PaginationRequest{
			Page:  util.Max(req.Page, 1),
			Limit: util.Max(req.Limit, 1),
		},
		ProductModelId:  req.ProductModelID,
		DateCreatedFrom: req.DateCreatedFrom,
		DateCreatedTo:   req.DateCreatedTo,
	})
	if err != nil {
		response.FromMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.FromPagination(w, http.StatusOK, model.PaginateResult[*pb.ProductEntity]{
		Data:       util.NonEmptySlice(resp.Products),
		Page:       resp.Pagination.Page,
		Limit:      resp.Pagination.Limit,
		Total:      resp.Pagination.Total,
		NextPage:   resp.Pagination.NextPage,
		NextCursor: resp.Pagination.NextCursor,
	})
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SerialID       string `json:"serial_id"`
		ProductModelID int64  `json:"product_model_id"`
	}
	if err := sonic.ConfigFastest.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	resp, err := h.client.CreateProduct(ctx, &pb.CreateProductRequest{
		SerialId:       req.SerialID,
		ProductModelId: req.ProductModelID,
	})
	if err != nil {
		response.FromMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.FromDTO(w, http.StatusCreated, struct {
		Id             int64  `json:"id"`
		SerialId       string `json:"serial_id"`
		ProductModelId int64  `json:"product_model_id"`
		DateCreated    int64  `json:"date_created"`
		DateUpdated    int64  `json:"date_updated"`
	}{
		Id:             resp.Id,
		SerialId:       resp.SerialId,
		ProductModelId: resp.ProductModelId,
		DateCreated:    resp.DateCreated,
		DateUpdated:    resp.DateUpdated,
	})
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.FromMessage(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		SerialID       *string `json:"serial_id"`
		ProductModelID *int64  `json:"product_model_id"`
	}
	if err := sonic.ConfigFastest.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	_, err = h.client.UpdateProduct(ctx, &pb.UpdateProductRequest{
		Id:             id,
		SerialId:       req.SerialID,
		ProductModelId: req.ProductModelID,
	})
	if err != nil {
		response.FromMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.FromMessage(w, http.StatusOK, "Product updated successfully")
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var (
		id       *int64
		serialId *string = util.ToPtr(chi.URLParam(r, "id"))
	)

	// If id is not a number, use the serialId
	idNum, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err == nil {
		id = &idNum
		serialId = nil
	}

	ctx := r.Context()

	_, err = h.client.DeleteProduct(ctx, &pb.DeleteProductRequest{
		Id:       id,
		SerialId: serialId,
	})
	if err != nil {
		response.FromMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.FromMessage(w, http.StatusOK, "Product deleted successfully")
}
