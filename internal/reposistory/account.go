package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	"shopnexus-go-service/internal/model"
)

type AddCartItemParams struct {
	CartID         []byte
	ProductModelID []byte
	Quantity       int64
}

func (r *Repository) AddCartItem(ctx context.Context, params AddCartItemParams) (int64, error) {
	return r.sqlc.AddCartItem(ctx, sqlc.AddCartItemParams{
		CartID:         params.CartID,
		ProductModelID: params.ProductModelID,
		Quantity:       params.Quantity,
	})
}

type DeductCartItemParams struct {
	CartID         []byte
	ProductModelID []byte
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
	CartID         []byte
	ProductModelID []byte
}

func (r *Repository) RemoveCartItem(ctx context.Context, params RemoveCartItemParams) error {
	return r.sqlc.RemoveCartItem(ctx, sqlc.RemoveCartItemParams{
		CartID:         params.CartID,
		ProductModelID: params.ProductModelID,
	})
}

func (r *Repository) GetCartItems(ctx context.Context, cartID []byte) ([]model.ItemOnCart, error) {
	rows, err := r.sqlc.GetCartItems(ctx, cartID)
	if err != nil {
		return nil, err
	}

	items := make([]model.ItemOnCart, len(rows))
	for i, row := range rows {
		items[i] = model.ItemOnCart{
			ItemQuantityBase: model.ItemQuantityBase{
				ItemID:   row.ProductModelID,
				Quantity: row.Quantity,
			},
			CartID: row.CartID,
		}
	}

	return items, nil
}
