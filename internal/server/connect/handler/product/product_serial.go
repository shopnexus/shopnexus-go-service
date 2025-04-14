package product

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/repository"
	"shopnexus-go-service/internal/service/product"

	common_grpc "shopnexus-go-service/internal/server/connect/handler/common"
	"shopnexus-go-service/internal/server/connect/interceptor/auth"

	"connectrpc.com/connect"
	productv1 "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/product/v1"
)

func (s *ImplementedProductServiceHandler) GetProductSerial(ctx context.Context, req *connect.Request[productv1.GetProductSerialRequest]) (*connect.Response[productv1.GetProductSerialResponse], error) {
	data, err := s.service.GetProductSerial(ctx, req.Msg.SerialId)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.GetProductSerialResponse{
		Data: modelToProductSerialEntity(data),
	}), nil
}

func (s *ImplementedProductServiceHandler) ListProductSerials(ctx context.Context, req *connect.Request[productv1.ListProductSerialsRequest]) (*connect.Response[productv1.ListProductSerialsResponse], error) {
	data, err := s.service.ListProductSerials(ctx, product.ListProductSerialsParams{
		PaginationParams: model.PaginationParams{
			Page:  req.Msg.GetPagination().GetPage(),
			Limit: req.Msg.GetPagination().GetLimit(),
		},
		SerialID:        req.Msg.SerialId,
		ProductID:       req.Msg.ProductId,
		IsSold:          req.Msg.IsSold,
		IsActive:        req.Msg.IsActive,
		DateCreatedFrom: req.Msg.DateCreatedFrom,
		DateCreatedTo:   req.Msg.DateCreatedTo,
	})
	if err != nil {
		return nil, err
	}

	var serials []*productv1.ProductSerialEntity
	for _, d := range data.Data {
		serials = append(serials, modelToProductSerialEntity(d))
	}

	return connect.NewResponse(&productv1.ListProductSerialsResponse{
		Data:       serials,
		Pagination: common_grpc.ToProtoPaginationResponse(data),
	}), nil
}

func (s *ImplementedProductServiceHandler) CreateProductSerial(ctx context.Context, req *connect.Request[productv1.CreateProductSerialRequest]) (*connect.Response[productv1.CreateProductSerialResponse], error) {
	data, err := s.service.CreateProductSerial(ctx, model.ProductSerial{
		SerialID:  req.Msg.SerialId,
		ProductID: req.Msg.ProductId,
		IsSold:    req.Msg.IsSold,
		IsActive:  req.Msg.IsActive,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.CreateProductSerialResponse{
		Data: modelToProductSerialEntity(data),
	}), nil
}

func (s *ImplementedProductServiceHandler) UpdateProductSerial(ctx context.Context, req *connect.Request[productv1.UpdateProductSerialRequest]) (*connect.Response[productv1.UpdateProductSerialResponse], error) {
	err := s.service.UpdateProductSerial(ctx, repository.UpdateProductSerialParams{
		SerialID: req.Msg.GetSerialId(),
		IsSold:   req.Msg.IsSold,
		IsActive: req.Msg.IsActive,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.UpdateProductSerialResponse{}), nil
}

func (s *ImplementedProductServiceHandler) DeleteProductSerial(ctx context.Context, req *connect.Request[productv1.DeleteProductSerialRequest]) (*connect.Response[productv1.DeleteProductSerialResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	err = s.service.DeleteProductSerial(ctx, product.DeleteProductSerialPParams{
		AccountID: claims.UserID,
		Role:      model.Role(claims.Role),
		SerialID:  req.Msg.SerialId,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.DeleteProductSerialResponse{}), nil
}

func modelToProductSerialEntity(data model.ProductSerial) *productv1.ProductSerialEntity {
	return &productv1.ProductSerialEntity{
		SerialId:    data.SerialID,
		ProductId:   data.ProductID,
		IsSold:      data.IsSold,
		IsActive:    data.IsActive,
		DateCreated: data.DateCreated,
		DateUpdated: data.DateUpdated,
	}
}
