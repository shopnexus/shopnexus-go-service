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

func (s *ImplementedProductServiceHandler) GetBrand(ctx context.Context, req *connect.Request[productv1.GetBrandRequest]) (*connect.Response[productv1.GetBrandResponse], error) {
	data, err := s.service.GetBrand(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.GetBrandResponse{
		Data: modelToBrandEntity(data),
	}), nil
}

func (s *ImplementedProductServiceHandler) ListBrands(ctx context.Context, req *connect.Request[productv1.ListBrandsRequest]) (*connect.Response[productv1.ListBrandsResponse], error) {
	data, err := s.service.ListBrands(ctx, product.ListBrandsParams{
		PaginationParams: model.PaginationParams{
			Page:  req.Msg.GetPagination().GetPage(),
			Limit: req.Msg.GetPagination().GetLimit(),
		},
		Name:        req.Msg.Name,
		Description: req.Msg.Description,
	})
	if err != nil {
		return nil, err
	}

	var brands []*productv1.BrandEntity
	for _, d := range data.Data {
		brands = append(brands, modelToBrandEntity(d))
	}

	return connect.NewResponse(&productv1.ListBrandsResponse{
		Data:       brands,
		Pagination: common_grpc.ToProtoPaginationResponse(data),
	}), nil
}

func (s *ImplementedProductServiceHandler) CreateBrand(ctx context.Context, req *connect.Request[productv1.CreateBrandRequest]) (*connect.Response[productv1.CreateBrandResponse], error) {
	claims, err := auth.GetAccount(req)
	if err != nil {
		return nil, err
	}

	data, err := s.service.CreateBrand(ctx, product.CreateBrandParams{
		UserID: claims.UserID,
		Brand: model.Brand{
			Name:        req.Msg.Name,
			Description: req.Msg.Description,
			Resources:   req.Msg.Resources,
		},
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.CreateBrandResponse{
		Data: modelToBrandEntity(data),
	}), nil
}

func (s *ImplementedProductServiceHandler) UpdateBrand(ctx context.Context, req *connect.Request[productv1.UpdateBrandRequest]) (*connect.Response[productv1.UpdateBrandResponse], error) {
	err := s.service.UpdateBrand(ctx, product.UpdateBrandParams{
		RepoParams: repository.UpdateBrandParams{
			ID:          req.Msg.Id,
			Name:        req.Msg.Name,
			Description: req.Msg.Description,
		},
		Resources: req.Msg.Resources,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.UpdateBrandResponse{}), nil
}

func (s *ImplementedProductServiceHandler) DeleteBrand(ctx context.Context, req *connect.Request[productv1.DeleteBrandRequest]) (*connect.Response[productv1.DeleteBrandResponse], error) {
	err := s.service.DeleteBrand(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.DeleteBrandResponse{}), nil
}

func modelToBrandEntity(data model.Brand) *productv1.BrandEntity {
	return &productv1.BrandEntity{
		Id:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		Resources:   data.Resources,
	}
}
