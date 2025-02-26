package service

import (
	"context"
	"shopnexus-go-service/internal/model"
	repository "shopnexus-go-service/internal/repository"
)

type CartService struct {
	Repo *repository.Repository
}

func NewCartService(repo *repository.Repository) *CartService {
	return &CartService{
		Repo: repo,
	}
}

func (s *CartService) GetCart(ctx context.Context, userID int64) (model.Cart, error) {
	cart, err := s.Repo.GetCart(ctx, userID)
	if err != nil {
		return model.Cart{}, err
	}

	return cart, nil
}

type AddCartItemParams struct {
	UserID         int64
	ProductModelID int64
	Quantity       int64
}

func (s *CartService) AddCartItem(ctx context.Context, params AddCartItemParams) (int64, error) {
	txRepo, err := s.Repo.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer txRepo.Rollback(ctx)

	newQty, err := txRepo.AddCartItem(ctx, repository.AddCartItemParams{
		CartID:         params.UserID,
		ProductModelID: params.ProductModelID,
		Quantity:       params.Quantity,
	})
	if err != nil {
		return 0, err
	}

	if err = txRepo.Commit(ctx); err != nil {
		return 0, err
	}

	return newQty, nil
}

type UpdateCartItemParams struct {
	UserID         int64
	ProductModelID int64
	Quantity       int64
}

func (s *CartService) UpdateCartItem(ctx context.Context, params UpdateCartItemParams) (int64, error) {
	txRepo, err := s.Repo.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer txRepo.Rollback(ctx)

	newQty, err := txRepo.UpdateCartItem(ctx, repository.UpdateCartItemParams{
		CartID:         params.UserID,
		ProductModelID: params.ProductModelID,
		Quantity:       params.Quantity,
	})
	if err != nil {
		return 0, err
	}

	if newQty <= 0 {
		if err = txRepo.RemoveCartItem(ctx, repository.RemoveCartItemParams{
			CartID:         params.UserID,
			ProductModelID: params.ProductModelID,
		}); err != nil {
			return 0, err
		}

		newQty = 0
	}

	if err = txRepo.Commit(ctx); err != nil {
		return 0, err
	}

	return newQty, nil
}

func (s *CartService) ClearCart(ctx context.Context, userID int64) error {
	return s.Repo.ClearCart(ctx, userID)
}
