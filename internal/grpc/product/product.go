package product

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service/product"

	"connectrpc.com/connect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/common"
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

func (s *ProductServer) GetProductModel(ctx context.Context, req *connect.Request[productv1.GetProductModelRequest]) (*connect.Response[productv1.GetProductModelResponse], error) {
	data, err := s.service.GetProductModel(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.GetProductModelResponse{
		ProductModel: modelToProductModelEntity(data),
	}), nil
}

func (s *ProductServer) ListProductModels(ctx context.Context, req *connect.Request[productv1.ListProductModelsRequest]) (*connect.Response[productv1.ListProductModelsResponse], error) {
	data, err := s.service.ListProductModels(ctx, product.ListProductModelsParams{
		PaginationParams: model.PaginationParams{
			Page:  req.Msg.GetPagination().GetPage(),
			Limit: req.Msg.GetPagination().GetLimit(),
		},
		BrandID:              req.Msg.BrandId,
		Name:                 req.Msg.Name,
		Description:          req.Msg.Description,
		ListPriceFrom:        req.Msg.ListPriceFrom,
		ListPriceTo:          req.Msg.ListPriceTo,
		DateManufacturedFrom: req.Msg.DateManufacturedFrom,
		DateManufacturedTo:   req.Msg.DateManufacturedTo,
	})
	if err != nil {
		return nil, err
	}

	var models []*productv1.ProductModelEntity
	for _, d := range data.Data {
		models = append(models, modelToProductModelEntity(d))
	}

	return connect.NewResponse(&productv1.ListProductModelsResponse{
		Pagination:    modelToPaginationResponse(data),
		ProductModels: models,
	}), nil
}

func (s *ProductServer) CreateProductModel(ctx context.Context, req *connect.Request[productv1.CreateProductModelRequest]) (*connect.Response[productv1.CreateProductModelResponse], error) {
	data, err := s.service.CreateProductModel(ctx, product.CreateProductModelParams{
		ProductModel: model.ProductModel{
			BrandID:          req.Msg.BrandId,
			Name:             req.Msg.Name,
			Description:      req.Msg.Description,
			ListPrice:        req.Msg.ListPrice,
			DateManufactured: req.Msg.DateManufactured,
			Resources:        req.Msg.Resources,
			Tags:             req.Msg.Tags,
		},
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.CreateProductModelResponse{
		ProductModel: modelToProductModelEntity(data),
	}), nil
}

func (s *ProductServer) UpdateProductModel(ctx context.Context, req *connect.Request[productv1.UpdateProductModelRequest]) (*connect.Response[productv1.UpdateProductModelResponse], error) {
	err := s.service.UpdateProductModel(ctx, product.UpdateProductModelParams{
		ID:               req.Msg.Id,
		BrandID:          req.Msg.BrandId,
		Name:             req.Msg.Name,
		Description:      req.Msg.Description,
		ListPrice:        req.Msg.ListPrice,
		DateManufactured: req.Msg.DateManufactured,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.UpdateProductModelResponse{}), nil
}

func (s *ProductServer) DeleteProductModel(ctx context.Context, req *connect.Request[productv1.DeleteProductModelRequest]) (*connect.Response[productv1.DeleteProductModelResponse], error) {
	err := s.service.DeleteProductModel(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.DeleteProductModelResponse{}), nil
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
		Pagination: modelToPaginationResponse(data),
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

func modelToProductModelEntity(data model.ProductModel) *productv1.ProductModelEntity {
	return &productv1.ProductModelEntity{
		Id:               data.ID,
		BrandId:          data.BrandID,
		Name:             data.Name,
		Description:      data.Description,
		ListPrice:        data.ListPrice,
		DateManufactured: data.DateManufactured,
		Resources:        data.Resources,
		Tags:             data.Tags,
	}
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

func modelToPaginationResponse[T any](data model.PaginateResult[T]) *common.PaginationResponse {
	return &common.PaginationResponse{
		Page:     data.Page,
		Limit:    data.Limit,
		Total:    data.Total,
		NextPage: data.NextPage,
	}
}
