package product

import (
	"context"
	"shopnexus-go-service/internal/model"
	common_grpc "shopnexus-go-service/internal/server/connect/handler/common"
	"shopnexus-go-service/internal/server/connect/interceptor/auth"
	"shopnexus-go-service/internal/service/product"

	"connectrpc.com/connect"
	commentv1 "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/product/v1"
)

func (s *ImplementedProductServiceHandler) GetComment(ctx context.Context, req *connect.Request[commentv1.GetCommentRequest]) (*connect.Response[commentv1.GetCommentResponse], error) {
	data, err := s.service.GetComment(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&commentv1.GetCommentResponse{
		Data: modelToCommentEntity(data),
	}), nil
}

func (s *ImplementedProductServiceHandler) ListComments(ctx context.Context, req *connect.Request[commentv1.ListCommentsRequest]) (*connect.Response[commentv1.ListCommentsResponse], error) {
	data, err := s.service.ListComments(ctx, product.ListCommentsParams{
		PaginationParams: model.PaginationParams{
			Page:  req.Msg.GetPagination().GetPage(),
			Limit: req.Msg.GetPagination().GetLimit(),
		},
		AccountID:       req.Msg.UserId,
		DestID:          req.Msg.DestId,
		Body:            req.Msg.Body,
		UpvoteFrom:      req.Msg.UpvoteFrom,
		UpvoteTo:        req.Msg.UpvoteTo,
		DownvoteFrom:    req.Msg.DownvoteFrom,
		DownvoteTo:      req.Msg.DownvoteTo,
		ScoreFrom:       req.Msg.ScoreFrom,
		ScoreTo:         req.Msg.ScoreTo,
		DateCreatedFrom: req.Msg.DateCreatedFrom,
		DateCreatedTo:   req.Msg.DateCreatedTo,
	})
	if err != nil {
		return nil, err
	}

	var comments []*commentv1.CommentEntity
	for _, d := range data.Data {
		comments = append(comments, modelToCommentEntity(d))
	}

	return connect.NewResponse(&commentv1.ListCommentsResponse{
		Data:       comments,
		Pagination: common_grpc.ToProtoPaginationResponse(data),
	}), nil
}

func (s *ImplementedProductServiceHandler) CreateComment(ctx context.Context, req *connect.Request[commentv1.CreateCommentRequest]) (*connect.Response[commentv1.CreateCommentResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	err = s.service.CreateComment(ctx, product.CreateCommentParams{
		AccountID: claims.UserID,
		DestID:    req.Msg.DestId,
		Body:      req.Msg.Body,
		Resources: req.Msg.Resources,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&commentv1.CreateCommentResponse{}), nil
}

func (s *ImplementedProductServiceHandler) UpdateComment(ctx context.Context, req *connect.Request[commentv1.UpdateCommentRequest]) (*connect.Response[commentv1.UpdateCommentResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	err = s.service.UpdateComment(ctx, product.UpdateCommentParams{
		Role:      claims.Role,
		AccountID: claims.UserID,
		ID:        req.Msg.Id,
		Body:      req.Msg.Body,
		Resources: req.Msg.Resources,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&commentv1.UpdateCommentResponse{}), nil
}

func (s *ImplementedProductServiceHandler) DeleteComment(ctx context.Context, req *connect.Request[commentv1.DeleteCommentRequest]) (*connect.Response[commentv1.DeleteCommentResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	err = s.service.DeleteComment(ctx, product.DeleteCommentParams{
		Role:      claims.Role,
		ID:        req.Msg.Id,
		AccountID: claims.UserID,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&commentv1.DeleteCommentResponse{}), nil
}

func modelToCommentEntity(data model.Comment) *commentv1.CommentEntity {
	return &commentv1.CommentEntity{
		Id:          data.ID,
		UserId:      data.AccountID,
		DestId:      data.DestID,
		Body:        data.Body,
		Upvote:      data.Upvote,
		Downvote:    data.Downvote,
		Score:       data.Score,
		DateCreated: data.DateCreated,
		DateUpdated: data.DateUpdated,
		Resources:   data.Resources,
	}
}
