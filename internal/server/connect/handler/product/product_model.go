package product

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service/product"

	common_grpc "shopnexus-go-service/internal/server/connect/handler/common"

	"connectrpc.com/connect"
	productv1 "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/product/v1"
)

func (s *ImplementedProductServiceHandler) GetProductModel(ctx context.Context, req *connect.Request[productv1.GetProductModelRequest]) (*connect.Response[productv1.GetProductModelResponse], error) {
	data, err := s.service.GetProductModel(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	serialIds, err := s.service.GetProductSerialIDs(ctx, data.ID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.GetProductModelResponse{
		Data: modelToProductModelEntity(data, serialIds),
	}), nil
}

func (s *ImplementedProductServiceHandler) ListProductModels(ctx context.Context, req *connect.Request[productv1.ListProductModelsRequest]) (*connect.Response[productv1.ListProductModelsResponse], error) {
	data, err := s.service.ListProductModels(ctx, product.ListProductModelsParams{
		PaginationParams: model.PaginationParams{
			Page:  req.Msg.GetPagination().GetPage(),
			Limit: req.Msg.GetPagination().GetLimit(),
		},
		Type:                 req.Msg.Type,
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
		// serialIds, err := s.service.GetProductSerialIDs(ctx, d.ID)
		// if err != nil {
		// 	return nil, err
		// }
		models = append(models, modelToProductModelEntity(d, []string{}))
	}

	return connect.NewResponse(&productv1.ListProductModelsResponse{
		Data:       models,
		Pagination: common_grpc.ToProtoPaginationResponse(data),
	}), nil
}

func (s *ImplementedProductServiceHandler) CreateProductModel(ctx context.Context, req *connect.Request[productv1.CreateProductModelRequest]) (*connect.Response[productv1.CreateProductModelResponse], error) {
	data, err := s.service.CreateProductModel(ctx, product.CreateProductModelParams{
		ProductModel: model.ProductModel{
			Type:             req.Msg.Type,
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

	serialIds, err := s.service.GetProductSerialIDs(ctx, data.ID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.CreateProductModelResponse{
		Data: modelToProductModelEntity(data, serialIds),
	}), nil
}

func (s *ImplementedProductServiceHandler) UpdateProductModel(ctx context.Context, req *connect.Request[productv1.UpdateProductModelRequest]) (*connect.Response[productv1.UpdateProductModelResponse], error) {
	var (
		resources *[]string
		tags      *[]string
	)
	if req.Msg.Resources != nil {
		resources = &req.Msg.Resources
	}
	if req.Msg.Tags != nil {
		tags = &req.Msg.Tags
	}

	err := s.service.UpdateProductModel(ctx, product.UpdateProductModelParams{
		ID:               req.Msg.Id,
		Type:             req.Msg.Type,
		BrandID:          req.Msg.BrandId,
		Name:             req.Msg.Name,
		Description:      req.Msg.Description,
		ListPrice:        req.Msg.ListPrice,
		DateManufactured: req.Msg.DateManufactured,
		Resources:        resources,
		Tags:             tags,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.UpdateProductModelResponse{}), nil
}

func (s *ImplementedProductServiceHandler) DeleteProductModel(ctx context.Context, req *connect.Request[productv1.DeleteProductModelRequest]) (*connect.Response[productv1.DeleteProductModelResponse], error) {
	err := s.service.DeleteProductModel(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.DeleteProductModelResponse{}), nil
}

func (s *ImplementedProductServiceHandler) ListProductTypes(ctx context.Context, req *connect.Request[productv1.ListProductTypesRequest]) (*connect.Response[productv1.ListProductTypesResponse], error) {
	data, err := s.service.ListProductTypes(ctx, product.ListProductTypesParams{
		PaginationParams: model.PaginationParams{
			Page:  req.Msg.GetPagination().GetPage(),
			Limit: req.Msg.GetPagination().GetLimit(),
		},
		Name: req.Msg.Name,
	})
	if err != nil {
		return nil, err
	}

	var types []*productv1.ProductTypeEntity
	for _, d := range data {
		types = append(types, modelToProductTypeEntity(d))
	}

	return connect.NewResponse(&productv1.ListProductTypesResponse{
		Data: types,
	}), nil
}

func modelToProductModelEntity(data model.ProductModel, serialIds []string) *productv1.ProductModelEntity {
	return &productv1.ProductModelEntity{
		Id:               data.ID,
		Type:             data.Type,
		BrandId:          data.BrandID,
		Name:             data.Name,
		Description:      data.Description,
		ListPrice:        data.ListPrice,
		DateManufactured: data.DateManufactured,
		Resources:        data.Resources,
		Tags:             data.Tags,
		SerialIds:        serialIds,
	}
}

func modelToProductTypeEntity(data model.ProductType) *productv1.ProductTypeEntity {
	return &productv1.ProductTypeEntity{
		Id:   data.ID,
		Name: data.Name,
	}
}
