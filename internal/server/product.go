package server

import (
	"context"
	"shopnexus-go-service/gen/pb"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductServer struct {
	pb.UnimplementedProductServer
	service *service.ProductService
}

func NewProductServer(service *service.ProductService) *ProductServer {
	return &ProductServer{service: service}
}

func (s *ProductServer) GetProductModel(ctx context.Context, req *pb.GetProductModelRequest) (*pb.ProductModelEntity, error) {
	data, err := s.service.GetProductModel(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return modelToProductModelEntity(data), nil
}

func (s *ProductServer) ListProductModels(ctx context.Context, req *pb.ListProductModelsRequest) (*pb.ListProductModelsResponse, error) {
	data, err := s.service.ListProductModels(ctx, service.ListProductModelsParams{
		PaginationParams: model.PaginationParams{
			Offset: req.GetPagination().GetOffset(),
			Limit:  req.GetPagination().GetLimit(),
		},
		BrandID:              req.BrandId,
		Name:                 req.Name,
		Description:          req.Description,
		ListPriceFrom:        req.ListPriceFrom,
		ListPriceTo:          req.ListPriceTo,
		DateManufacturedFrom: req.DateManufacturedFrom,
		DateManufacturedTo:   req.DateManufacturedTo,
	})
	if err != nil {
		return nil, err
	}

	var models []*pb.ProductModelEntity
	for _, d := range data.Data {
		models = append(models, modelToProductModelEntity(d))
	}

	return &pb.ListProductModelsResponse{
		Pagination:    modelToPaginationResponse(data),
		ProductModels: models,
	}, nil
}

func (s *ProductServer) DeleteProductModel(ctx context.Context, req *pb.DeleteProductModelRequest) (*emptypb.Empty, error) {
	err := s.service.DeleteProductModel(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *ProductServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductEntity, error) {
	data, err := s.service.GetProduct(ctx, service.ProductIdentifier{
		ID:       req.Id,
		SerialID: req.SerialId,
	})
	if err != nil {
		return nil, err
	}

	return modelToProductEntity(data), nil
}

func (s *ProductServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	data, err := s.service.ListProducts(ctx, service.ListProductsParams{
		PaginationParams: model.PaginationParams{
			Offset: req.GetPagination().GetOffset(),
			Limit:  req.GetPagination().GetLimit(),
		},
		ProductModelID:  req.ProductModelId,
		DateCreatedFrom: req.DateCreatedFrom,
		DateCreatedTo:   req.DateCreatedTo,
	})
	if err != nil {
		return nil, err
	}

	var products []*pb.ProductEntity
	for _, d := range data.Data {
		products = append(products, modelToProductEntity(d))
	}

	return &pb.ListProductsResponse{
		Pagination: modelToPaginationResponse(data),
		Products:   products,
	}, nil
}

func (s *ProductServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductEntity, error) {
	data, err := s.service.CreateProduct(ctx, model.Product{
		SerialID:       req.SerialId,
		ProductModelID: req.ProductModelId,
	})
	if err != nil {
		return nil, err
	}

	return modelToProductEntity(data), nil
}

func (s *ProductServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*emptypb.Empty, error) {
	err := s.service.UpdateProduct(ctx, service.UpdateProductParams{
		ID:             req.GetId(),
		SerialID:       req.SerialId,
		ProductModelID: req.ProductModelId,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *ProductServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*emptypb.Empty, error) {
	err := s.service.DeleteProduct(ctx, service.ProductIdentifier{
		ID:       req.Id,
		SerialID: req.SerialId,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func modelToProductModelEntity(data model.ProductModel) *pb.ProductModelEntity {
	return &pb.ProductModelEntity{
		Id:               data.ID,
		BrandId:          data.BrandID,
		Name:             data.Name,
		Description:      data.Description,
		ListPrice:        data.ListPrice,
		DateManufactured: data.DateManufactured,
	}
}

func modelToProductEntity(data model.Product) *pb.ProductEntity {
	return &pb.ProductEntity{
		Id:             data.ID,
		SerialId:       data.SerialID,
		ProductModelId: data.ProductModelID,
		DateCreated:    data.DateCreated,
		DateUpdated:    data.DateUpdated,
	}
}
