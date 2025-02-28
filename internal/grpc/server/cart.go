package server

import (
	"context"
	"shopnexus-go-service/gen/pb"
	"shopnexus-go-service/internal/interceptor"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service"

	"google.golang.org/protobuf/types/known/emptypb"
)

type CartServer struct {
	pb.UnimplementedCartServer
	service *service.CartService
}

func NewCartServer(service *service.CartService) *CartServer {
	return &CartServer{service: service}
}

func (s *CartServer) AddCartItem(ctx context.Context, req *pb.AddCartItemRequest) (*pb.AddCartItemResponse, error) {
	claims, ok := ctx.Value(interceptor.CtxKeyUser).(model.Claims)
	if !ok {
		return nil, model.ErrTokenInvalid
	}

	quantity, err := s.service.AddCartItem(ctx, service.AddCartItemParams{
		UserID:         claims.UserID,
		ProductModelID: req.ProductModelId,
		Quantity:       req.Quantity,
	})
	if err != nil {
		return nil, err
	}

	return &pb.AddCartItemResponse{Quantity: quantity}, nil
}

func (s *CartServer) GetCart(ctx context.Context, req *emptypb.Empty) (*pb.GetCartResponse, error) {
	claims, ok := ctx.Value(interceptor.CtxKeyUser).(model.Claims)
	if !ok {
		return nil, model.ErrTokenInvalid
	}

	cart, err := s.service.GetCart(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	items := make([]*pb.ItemQuantity, len(cart.ProductModels))
	for i, item := range cart.ProductModels {
		items[i] = &pb.ItemQuantity{
			ItemId:   item.GetID(),
			Quantity: item.GetQuantity(),
		}
	}

	return &pb.GetCartResponse{
		Items: items,
	}, nil
}

func (s *CartServer) UpdateCartItem(ctx context.Context, req *pb.UpdateCartItemRequest) (*emptypb.Empty, error) {
	claims, ok := ctx.Value(interceptor.CtxKeyUser).(model.Claims)
	if !ok {
		return nil, model.ErrTokenInvalid
	}

	_, err := s.service.UpdateCartItem(ctx, service.UpdateCartItemParams{
		UserID:         claims.UserID,
		ProductModelID: req.ProductModelId,
		Quantity:       req.Quantity,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *CartServer) ClearCart(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	claims, ok := ctx.Value(interceptor.CtxKeyUser).(model.Claims)
	if !ok {
		return nil, model.ErrTokenInvalid
	}

	if err := s.service.ClearCart(ctx, claims.UserID); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
