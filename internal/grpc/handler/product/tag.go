package product

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service/product"

	common_grpc "shopnexus-go-service/internal/grpc/handler/common"

	"connectrpc.com/connect"
	productv1 "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/product/v1"
)

func (s *ImplementedProductServiceHandler) GetTag(ctx context.Context, req *connect.Request[productv1.GetTagRequest]) (*connect.Response[productv1.GetTagResponse], error) {
	data, err := s.service.GetTag(ctx, req.Msg.GetTag())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.GetTagResponse{
		Data: modelToTagEntity(data),
	}), nil
}

func (s *ImplementedProductServiceHandler) ListTags(ctx context.Context, req *connect.Request[productv1.ListTagsRequest]) (*connect.Response[productv1.ListTagsResponse], error) {
	data, err := s.service.ListTags(ctx, product.ListTagsParams{
		PaginationParams: model.PaginationParams{
			Page:  req.Msg.GetPagination().GetPage(),
			Limit: req.Msg.GetPagination().GetLimit(),
		},
		Tag:         req.Msg.Tag,
		Description: req.Msg.Description,
	})
	if err != nil {
		return nil, err
	}

	var tags []*productv1.TagEntity
	for _, d := range data.Data {
		tags = append(tags, modelToTagEntity(d))
	}

	return connect.NewResponse(&productv1.ListTagsResponse{
		Data:       tags,
		Pagination: common_grpc.ToProtoPaginationResponse(data),
	}), nil
}

func (s *ImplementedProductServiceHandler) CreateTag(ctx context.Context, req *connect.Request[productv1.CreateTagRequest]) (*connect.Response[productv1.CreateTagResponse], error) {
	err := s.service.CreateTag(ctx, model.Tag{
		Tag:         req.Msg.Tag,
		Description: req.Msg.Description,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.CreateTagResponse{}), nil
}

func (s *ImplementedProductServiceHandler) UpdateTag(ctx context.Context, req *connect.Request[productv1.UpdateTagRequest]) (*connect.Response[productv1.UpdateTagResponse], error) {
	if req.Msg.NewTag != nil && *req.Msg.NewTag == "" {
		return nil, model.ErrMalformedParams
	}

	err := s.service.UpdateTag(ctx, product.UpdateTagParams{
		Tag:         req.Msg.Tag,
		NewTag:      req.Msg.NewTag,
		Description: req.Msg.Description,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.UpdateTagResponse{}), nil
}

func (s *ImplementedProductServiceHandler) DeleteTag(ctx context.Context, req *connect.Request[productv1.DeleteTagRequest]) (*connect.Response[productv1.DeleteTagResponse], error) {
	err := s.service.DeleteTag(ctx, req.Msg.GetTag())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&productv1.DeleteTagResponse{}), nil
}

type TagResponseParams struct {
	ProductCount int32
}

func modelToTagEntity(data product.TagResponse) *productv1.TagEntity {
	return &productv1.TagEntity{
		Tag:          data.Tag.Tag,
		Description:  data.Tag.Description,
		ProductCount: data.ProductCount,
	}
}
