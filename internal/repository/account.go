package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	"shopnexus-go-service/internal/model"
)

type AddCartItemParams struct {
	CartID         int64
	ProductModelID int64
	Quantity       int64
}

func (r *Repository) OwnCart(ctx context.Context, userID int64, cartID int64) (bool, error) {
	return userID == cartID, nil
}

func (r *Repository) AddCartItem(ctx context.Context, params AddCartItemParams) (int64, error) {
	return r.sqlc.AddCartItem(ctx, sqlc.AddCartItemParams{
		CartID:         params.CartID,
		ProductModelID: params.ProductModelID,
		Quantity:       params.Quantity,
	})
}

type DeductCartItemParams struct {
	CartID         int64
	ProductModelID int64
	Quantity       int64
}

func (r *Repository) DeductCartItem(ctx context.Context, params DeductCartItemParams) (int64, error) {
	return r.sqlc.DeductCartItem(ctx, sqlc.DeductCartItemParams{
		CartID:         params.CartID,
		ProductModelID: params.ProductModelID,
		Quantity:       params.Quantity,
	})
}

type RemoveCartItemParams struct {
	CartID         int64
	ProductModelID int64
}

func (r *Repository) RemoveCartItem(ctx context.Context, params RemoveCartItemParams) error {
	return r.sqlc.RemoveCartItem(ctx, sqlc.RemoveCartItemParams{
		CartID:         params.CartID,
		ProductModelID: params.ProductModelID,
	})
}

func (r *Repository) GetCart(ctx context.Context, cartID int64) (model.Cart, error) {
	cartRow, err := r.sqlc.GetCart(ctx, cartID)
	if err != nil {
		return model.Cart{}, err
	}

	itemRows, err := r.sqlc.GetCartItems(ctx, cartID)
	if err != nil {
		return model.Cart{}, err
	}

	items := make([]model.ItemQuantity[int64], len(itemRows))
	for i, row := range itemRows {
		items[i] = model.ItemOnCart{
			ItemQuantityBase: model.ItemQuantityBase[int64]{
				ItemID:   row.ProductModelID,
				Quantity: row.Quantity,
			},
			CartID: row.CartID,
		}
	}

	return model.Cart{
		ID:       cartRow,
		Products: items,
	}, nil
}

func (r *Repository) CreateCart(ctx context.Context, userID int64) error {
	return r.sqlc.CreateCart(ctx, userID)
}
