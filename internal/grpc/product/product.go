package product

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service/product"

	common_grpc "shopnexus-go-service/internal/grpc/common"

	"connectrpc.com/connect"
	productv1 "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/product/v1"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/product/v1/productv1connect"
)

var _ productv1connect.ProductServiceHandler = (*ProductServer)(nil)

type ProductServer struct {
	productv1connect.UnimplementedProductServiceHandler
	service *product.ProductService
}

func NewProductServer(service *product.ProductService) *ProductServer {
	return &ProductServer{service: service}
}

func (s *ProductServer) GetProduct(ctx context.Context, req *connect.Request[productv1.GetProductRequest]) (*connect.Response[productv1.GetProductResponse], error) {
	data, err := s.service.GetProduct(ctx, model.ProductIdentifier{
		ID:       req.Msg.Id,
		SerialID: req.Msg.SerialId,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.GetProductResponse{
		Product: modelToProductEntity(data),
	}), nil
}

func (s *ProductServer) ListProducts(ctx context.Context, req *connect.Request[productv1.ListProductsRequest]) (*connect.Response[productv1.ListProductsResponse], error) {
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
		Pagination: common_grpc.ToProtoPaginationResponse(data),
		Products:   products,
	}), nil
}

func (s *ProductServer) CreateProduct(ctx context.Context, req *connect.Request[productv1.CreateProductRequest]) (*connect.Response[productv1.CreateProductResponse], error) {
	data, err := s.service.CreateProduct(ctx, model.Product{
		SerialID:       req.Msg.SerialId,
		ProductModelID: req.Msg.ProductModelId,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.CreateProductResponse{
		Product: modelToProductEntity(data),
	}), nil
}

func (s *ProductServer) UpdateProduct(ctx context.Context, req *connect.Request[productv1.UpdateProductRequest]) (*connect.Response[productv1.UpdateProductResponse], error) {
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

func (s *ProductServer) DeleteProduct(ctx context.Context, req *connect.Request[productv1.DeleteProductRequest]) (*connect.Response[productv1.DeleteProductResponse], error) {
	err := s.service.DeleteProduct(ctx, model.ProductIdentifier{
		ID:       req.Msg.Id,
		SerialID: req.Msg.SerialId,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.DeleteProductResponse{}), nil
}

func modelToProductEntity(data model.Product) *productv1.ProductEntity {
	return &productv1.ProductEntity{
		Id:             data.ID,
		SerialId:       data.SerialID,
		ProductModelId: data.ProductModelID,
		DateCreated:    data.DateCreated,
		DateUpdated:    data.DateUpdated,
	}
}
