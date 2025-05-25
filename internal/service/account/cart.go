package account

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service/storage"
)

func (s *ServiceImpl) GetCart(ctx context.Context, userID int64) (model.Cart, error) {
	cart, err := s.storage.GetCart(ctx, storage.GetCartParams{
		CartID: userID,
	})
	if err != nil {
		return model.Cart{}, err
	}

	return cart, nil
}

type AddCartItemParams struct {
	UserID    int64
	ProductID int64
	Quantity  int64
}

func (s *ServiceImpl) AddCartItem(ctx context.Context, params AddCartItemParams) (int64, error) {
	txStorage, err := s.storage.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer txStorage.Rollback(ctx)

	// Check if cart exists for the user, create one if it doesn't
	exists, err := txStorage.ExistsCart(ctx, params.UserID)
	if err != nil {
		return 0, err
	}
	if !exists {
		if err = txStorage.CreateCart(ctx, params.UserID); err != nil {
			return 0, err
		}
	}

	newQty, err := txStorage.AddCartItem(ctx, storage.AddCartItemParams{
		CartID:    params.UserID,
		ProductID: params.ProductID,
		Quantity:  params.Quantity,
	})
	if err != nil {
		return 0, err
	}

	if err = txStorage.Commit(ctx); err != nil {
		return 0, err
	}

	return newQty, nil
}

type UpdateCartItemParams struct {
	UserID    int64
	ProductID int64
	Quantity  int64
}

func (s *ServiceImpl) UpdateCartItem(ctx context.Context, params UpdateCartItemParams) (int64, error) {
	txStorage, err := s.storage.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer txStorage.Rollback(ctx)

	newQty, err := txStorage.UpdateCartItem(ctx, storage.UpdateCartItemParams{
		CartID:    params.UserID,
		ProductID: params.ProductID,
		Quantity:  params.Quantity,
	})
	if err != nil {
		return 0, err
	}

	if newQty <= 0 {
		if err = txStorage.RemoveCartItem(ctx, params.UserID, []int64{params.ProductID}); err != nil {
			return 0, err
		}

		newQty = 0
	}

	if err = txStorage.Commit(ctx); err != nil {
		return 0, err
	}

	return newQty, nil
}

func (s *ServiceImpl) ClearCart(ctx context.Context, userID int64) error {
	return s.storage.ClearCart(ctx, userID)
}
