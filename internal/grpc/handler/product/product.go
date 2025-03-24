package product

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service/product"

	common_grpc "shopnexus-go-service/internal/grpc/handler/common"

	"connectrpc.com/connect"
	productv1 "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/product/v1"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/product/v1/productv1connect"
)

var _ productv1connect.ProductServiceHandler = (*ImplementedProductServiceHandler)(nil)

type ImplementedProductServiceHandler struct {
	productv1connect.UnimplementedProductServiceHandler
	service *product.ProductService
}

func NewProductServiceHandler(service *product.ProductService) productv1connect.ProductServiceHandler {
	return &ImplementedProductServiceHandler{service: service}
}

func (s *ImplementedProductServiceHandler) GetProduct(ctx context.Context, req *connect.Request[productv1.GetProductRequest]) (*connect.Response[productv1.GetProductResponse], error) {
	data, err := s.service.GetProduct(ctx, model.ProductIdentifier{
		ID:       req.Msg.Id,
		SerialID: req.Msg.SerialId,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.GetProductResponse{
		Data: modelToProductEntity(data),
	}), nil
}

func (s *ImplementedProductServiceHandler) ListProducts(ctx context.Context, req *connect.Request[productv1.ListProductsRequest]) (*connect.Response[productv1.ListProductsResponse], error) {
	data, err := s.service.ListProducts(ctx, product.ListProductsParams{
		PaginationParams: model.PaginationParams{
			Page:  req.Msg.GetPagination().GetPage(),
			Limit: req.Msg.GetPagination().GetLimit(),
		},
		ProductModelID:  req.Msg.ProductModelId,
		DateCreatedFrom: req.Msg.DateCreatedFrom,
		DateCreatedTo:   req.Msg.DateCreatedTo,
	})
	if err != nil {
		return nil, err
	}

	var products []*productv1.ProductEntity
	for _, d := range data.Data {
		products = append(products, modelToProductEntity(d))
	}

	return connect.NewResponse(&productv1.ListProductsResponse{
		Data:       products,
		Pagination: common_grpc.ToProtoPaginationResponse(data),
	}), nil
}

func (s *ImplementedProductServiceHandler) CreateProduct(ctx context.Context, req *connect.Request[productv1.CreateProductRequest]) (*connect.Response[productv1.CreateProductResponse], error) {
	data, err := s.service.CreateProduct(ctx, model.Product[any]{
		SerialID:       req.Msg.SerialId,
		ProductModelID: req.Msg.ProductModelId,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.CreateProductResponse{
		Data: modelToProductEntity(data),
	}), nil
}

func (s *ImplementedProductServiceHandler) UpdateProduct(ctx context.Context, req *connect.Request[productv1.UpdateProductRequest]) (*connect.Response[productv1.UpdateProductResponse], error) {
	err := s.service.UpdateProduct(ctx, product.UpdateProductParams{
		ID:             req.Msg.GetId(),
		SerialID:       req.Msg.SerialId,
		ProductModelID: req.Msg.ProductModelId,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.UpdateProductResponse{}), nil
}

func (s *ImplementedProductServiceHandler) DeleteProduct(ctx context.Context, req *connect.Request[productv1.DeleteProductRequest]) (*connect.Response[productv1.DeleteProductResponse], error) {
	err := s.service.DeleteProduct(ctx, model.ProductIdentifier{
		ID:       req.Msg.Id,
		SerialID: req.Msg.SerialId,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.DeleteProductResponse{}), nil
}

func modelToProductEntity(data model.Product[any]) *productv1.ProductEntity {
	return &productv1.ProductEntity{
		Id:             data.ID,
		SerialId:       data.SerialID,
		ProductModelId: data.ProductModelID,
		DateCreated:    data.DateCreated,
		DateUpdated:    data.DateUpdated,
	}
}
