package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	"shopnexus-go-service/internal/model"
)

// ExistsCart checks if a cart exists for the given user ID
func (r *Repository) ExistsCart(ctx context.Context, userID int64) (bool, error) {
	return r.sqlc.ExistsCart(ctx, userID)
}

// CreateCart creates a new cart for the user
func (r *Repository) CreateCart(ctx context.Context, userID int64) error {
	return r.sqlc.CreateCart(ctx, userID)
}

// GetCart retrieves the cart with the given ID
func (r *Repository) GetCart(ctx context.Context, cartID int64) (model.Cart, error) {
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
		ID:            cartID,
		ProductModels: items,
	}, nil
}

type AddCartItemParams struct {
	CartID         int64
	ProductModelID int64
	Quantity       int64
}

// AddCartItem adds a new item to the cart
func (r *Repository) AddCartItem(ctx context.Context, params AddCartItemParams) (int64, error) {
	return r.sqlc.AddCartItem(ctx, sqlc.AddCartItemParams{
		CartID:         params.CartID,
		ProductModelID: params.ProductModelID,
		Quantity:       params.Quantity,
	})
}

type UpdateCartItemParams struct {
	CartID         int64
	ProductModelID int64
	Quantity       int64
}

func (r *Repository) UpdateCartItem(ctx context.Context, params UpdateCartItemParams) (int64, error) {
	return r.sqlc.UpdateCartItem(ctx, sqlc.UpdateCartItemParams{
		CartID:         params.CartID,
		ProductModelID: params.ProductModelID,
		Quantity:       params.Quantity,
	})
}

func (r *Repository) RemoveCartItem(ctx context.Context, cartID, productModelID int64) error {
	return r.sqlc.RemoveCartItem(ctx, sqlc.RemoveCartItemParams{
		CartID:         cartID,
		ProductModelID: productModelID,
	})
}

func (r *Repository) ClearCart(ctx context.Context, cartID int64) error {
	return r.sqlc.ClearCart(ctx, cartID)
}
