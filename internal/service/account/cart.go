package account

import (
	"context"
	"shopnexus-go-service/internal/model"
	repository "shopnexus-go-service/internal/repository"
)

func (s *AccountService) GetCart(ctx context.Context, userID int64) (model.Cart, error) {
	cart, err := s.repo.GetCart(ctx, userID)
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

func (s *AccountService) AddCartItem(ctx context.Context, params AddCartItemParams) (int64, error) {
	txRepo, err := s.repo.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer txRepo.Rollback(ctx)

	// Check if cart exists for the user, create one if it doesn't
	exists, err := txRepo.ExistsCart(ctx, params.UserID)
	if err != nil {
		return 0, err
	}
	if !exists {
		if err = txRepo.CreateCart(ctx, params.UserID); err != nil {
			return 0, err
		}
	}

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

func (s *AccountService) UpdateCartItem(ctx context.Context, params UpdateCartItemParams) (int64, error) {
	txRepo, err := s.repo.Begin(ctx)
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
		if err = txRepo.RemoveCartItem(ctx, params.UserID, params.ProductModelID); err != nil {
			return 0, err
		}

		newQty = 0
	}

	if err = txRepo.Commit(ctx); err != nil {
		return 0, err
	}

	return newQty, nil
}

func (s *AccountService) ClearCart(ctx context.Context, userID int64) error {
	return s.repo.ClearCart(ctx, userID)
}
