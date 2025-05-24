package product

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service/product"

	common_grpc "shopnexus-go-service/internal/transport/connect/handler/common"
	"shopnexus-go-service/internal/transport/connect/interceptor/auth"

	"connectrpc.com/connect"
	productv1 "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/product/v1"
)

func (s *ImplementedProductServiceHandler) GetSale(ctx context.Context, req *connect.Request[productv1.GetSaleRequest]) (*connect.Response[productv1.GetSaleResponse], error) {
	data, err := s.service.GetSale(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.GetSaleResponse{
		Data: modelToSaleEntity(data),
	}), nil
}

func (s *ImplementedProductServiceHandler) ListSales(ctx context.Context, req *connect.Request[productv1.ListSalesRequest]) (*connect.Response[productv1.ListSalesResponse], error) {
	data, err := s.service.ListSales(ctx, product.ListSalesParams{
		PaginationParams: model.PaginationParams{
			Page:  req.Msg.GetPagination().GetPage(),
			Limit: req.Msg.GetPagination().GetLimit(),
		},
		Tag:             req.Msg.Tag,
		ProductModelID:  req.Msg.ProductModelId,
		BrandID:         req.Msg.BrandId,
		DateStartedFrom: req.Msg.DateStartedFrom,
		DateStartedTo:   req.Msg.DateStartedTo,
		DateEndedFrom:   req.Msg.DateEndedFrom,
		DateEndedTo:     req.Msg.DateEndedTo,
		IsActive:        req.Msg.IsActive,
	})
	if err != nil {
		return nil, err
	}

	var sales []*productv1.SaleEntity
	for _, d := range data.Data {
		sales = append(sales, modelToSaleEntity(d))
	}

	return connect.NewResponse(&productv1.ListSalesResponse{
		Data:       sales,
		Pagination: common_grpc.ToProtoPaginationResponse(data),
	}), nil
}

func (s *ImplementedProductServiceHandler) CreateSale(ctx context.Context, req *connect.Request[productv1.CreateSaleRequest]) (*connect.Response[productv1.CreateSaleResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	data, err := s.service.CreateSale(ctx, product.CreateSaleParams{
		UserID: claims.UserID,
		Sale: model.Sale{
			Tag:              req.Msg.Tag,
			ProductModelID:   req.Msg.ProductModelId,
			BrandID:          req.Msg.BrandId,
			DateStarted:      req.Msg.DateStarted,
			DateEnded:        req.Msg.DateEnded,
			Quantity:         req.Msg.Quantity,
			IsActive:         req.Msg.IsActive,
			DiscountPercent:  req.Msg.DiscountPercent,
			DiscountPrice:    req.Msg.DiscountPrice,
			MaxDiscountPrice: req.Msg.MaxDiscountPrice,
		},
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.CreateSaleResponse{
		Data: modelToSaleEntity(data),
	}), nil
}

func (s *ImplementedProductServiceHandler) UpdateSale(ctx context.Context, req *connect.Request[productv1.UpdateSaleRequest]) (*connect.Response[productv1.UpdateSaleResponse], error) {
	err := s.service.UpdateSale(ctx, product.UpdateSaleParams{
		ID:              req.Msg.Id,
		Tag:             req.Msg.Tag,
		ProductModelID:  req.Msg.ProductModelId,
		BrandID:         req.Msg.BrandId,
		DateStarted:     req.Msg.DateStarted,
		DateEnded:       req.Msg.DateEnded,
		Quantity:        req.Msg.Quantity,
		Used:            req.Msg.Used,
		IsActive:        req.Msg.IsActive,
		DiscountPercent: req.Msg.DiscountPercent,
		DiscountPrice:   req.Msg.DiscountPrice,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.UpdateSaleResponse{}), nil
}

func (s *ImplementedProductServiceHandler) DeleteSale(ctx context.Context, req *connect.Request[productv1.DeleteSaleRequest]) (*connect.Response[productv1.DeleteSaleResponse], error) {
	err := s.service.DeleteSale(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.DeleteSaleResponse{}), nil
}

func modelToSaleEntity(data model.Sale) *productv1.SaleEntity {
	return &productv1.SaleEntity{
		Id:               data.ID,
		Tag:              data.Tag,
		ProductModelId:   data.ProductModelID,
		BrandId:          data.BrandID,
		DateCreated:      data.DateCreated,
		DateStarted:      data.DateStarted,
		DateEnded:        data.DateEnded,
		Quantity:         data.Quantity,
		Used:             data.Used,
		IsActive:         data.IsActive,
		DiscountPercent:  data.DiscountPercent,
		DiscountPrice:    data.DiscountPrice,
		MaxDiscountPrice: data.MaxDiscountPrice,
	}
}
