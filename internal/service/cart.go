package service

import (
	"context"
	"shopnexus-go-service/gen/pb"
	"shopnexus-go-service/internal/model"
	repository "shopnexus-go-service/internal/reposistory"
)

type CartService struct {
	Repo *repository.Repository
	pb.UnimplementedCartServer
}

func NewCartService() *CartService {
	return &CartService{}
}

func (s *CartService) AddItem(ctx context.Context, params *pb.AddItemRequest) (*pb.AddItemResponse, error) {
	txRepo, err := s.Repo.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer txRepo.Rollback(ctx)

	isOwn, err := s.Repo.OwnCart(ctx, params.UserId, params.CartId)
	if err != nil {
		return nil, err
	}
	if !isOwn {
		return nil, model.ErrForbidden
	}

	newQty, err := txRepo.AddCartItem(ctx, repository.AddCartItemParams{
		CartID:         params.UserId,
		ProductModelID: params.ProductModelId,
		Quantity:       params.Quantity,
	})

	if err = txRepo.Commit(ctx); err != nil {
		return nil, err
	}

	return &pb.AddItemResponse{
		Quantity: newQty,
	}, nil
}

func (s *CartService) GetCart(ctx context.Context, params *pb.GetCartRequest) (*pb.GetCartResponse, error) {
	isOwn, err := s.Repo.OwnCart(ctx, params.UserId, params.CartId)
	if err != nil {
		return nil, err
	}
	if !isOwn {
		return nil, model.ErrForbidden
	}

	cart, err := s.Repo.GetCart(ctx, params.CartId)
	if err != nil {
		return nil, err
	}

	return &pb.GetCartResponse{
		CartId: cart.ID,
		Products: func() []*pb.ItemQuantity {
			products := make([]*pb.ItemQuantity, len(cart.Products))
			for i, p := range cart.Products {
				products[i] = &pb.ItemQuantity{
					ItemId:   p.GetID(),
					Quantity: p.GetQuantity(),
				}
			}
			return products
		}(),
	}, nil
}

func (s *CartService) RemoveItem(ctx context.Context, params *pb.RemoveItemRequest) (*pb.RemoveItemResponse, error) {
	// isOwn, err := s.Repo.OwnCart(ctx, params, params.CartId)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

func (s *CartService) Clear(ctx context.Context, params *pb.GetCartRequest) (*pb.GetCartResponse, error) {
	return nil, nil
}

func (s *CartService) Checkout(ctx context.Context, params *pb.GetCartRequest) (*pb.GetCartResponse, error) {
	return nil, nil
}
