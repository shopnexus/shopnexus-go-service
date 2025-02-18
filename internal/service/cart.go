package service

import (
	"context"
	"shopnexus-go-service/gen/pb"
)

type CartService struct {
	pb.UnimplementedCartServer
}

func NewCartService() *CartService {
	return &CartService{}
}

func (s *CartService) AddItem(ctx context.Context, params *pb.AddItemRequest) (*pb.AddItemResponse, error) {
	return nil, nil
}

func (s *CartService) GetCart(ctx context.Context, params *pb.GetCartRequest) (*pb.GetCartResponse, error) {
	return nil, nil
}

func (s *CartService) RemoveItem(ctx context.Context, params *pb.RemoveItemRequest) (*pb.RemoveItemResponse, error) {
	return nil, nil
}

func (s *CartService) Clear(ctx context.Context, params *pb.GetCartRequest) (*pb.GetCartResponse, error) {
	return nil, nil
}

func (s *CartService) Checkout(ctx context.Context, params *pb.GetCartRequest) (*pb.GetCartResponse, error) {
	return nil, nil
}
