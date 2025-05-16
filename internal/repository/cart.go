package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	"shopnexus-go-service/internal/model"
)

// ExistsCart checks if a cart exists for the given user ID
func (r *RepositoryImpl) ExistsCart(ctx context.Context, userID int64) (bool, error) {
	return r.sqlc.ExistsCart(ctx, userID)
}

// CreateCart creates a new cart for the user
func (r *RepositoryImpl) CreateCart(ctx context.Context, userID int64) error {
	return r.sqlc.CreateCart(ctx, userID)
}

type GetCartParams struct {
	CartID     int64
	ProductIDs []int64 // List of product IDs to retrieve, if empty, retrieves all items in the cart
}

// GetCart retrieves the cart with the given ID
func (r *RepositoryImpl) GetCart(ctx context.Context, params GetCartParams) (model.Cart, error) {
	itemRows, err := r.sqlc.GetCartItems(ctx, sqlc.GetCartItemsParams{
		CartID:     params.CartID,
		ProductIds: params.ProductIDs,
	})
	if err != nil {
		return model.Cart{}, err
	}

	items := make([]model.ItemQuantity[int64], len(itemRows))
	for i, row := range itemRows {
		items[i] = model.ItemOnCart{
			ItemQuantityBase: model.ItemQuantityBase[int64]{
				ItemID:   row.ProductID,
				Quantity: row.Quantity,
			},
			CartID: row.CartID,
		}
	}

	return model.Cart{
		ID:       params.CartID,
		Products: items,
	}, nil
}

type AddCartItemParams struct {
	CartID    int64
	ProductID int64
	Quantity  int64
}

// AddCartItem adds a new item to the cart
func (r *RepositoryImpl) AddCartItem(ctx context.Context, params AddCartItemParams) (int64, error) {
	return r.sqlc.AddCartItem(ctx, sqlc.AddCartItemParams{
		CartID:    params.CartID,
		ProductID: params.ProductID,
		Quantity:  params.Quantity,
	})
}

type UpdateCartItemParams struct {
	CartID    int64
	ProductID int64
	Quantity  int64
}

func (r *RepositoryImpl) UpdateCartItem(ctx context.Context, params UpdateCartItemParams) (int64, error) {
	return r.sqlc.UpdateCartItem(ctx, sqlc.UpdateCartItemParams{
		CartID:    params.CartID,
		ProductID: params.ProductID,
		Quantity:  params.Quantity,
	})
}

func (r *RepositoryImpl) RemoveCartItem(ctx context.Context, cartID int64, productIDs []int64) error {
	return r.sqlc.RemoveCartItem(ctx, sqlc.RemoveCartItemParams{
		CartID:     cartID,
		ProductIds: productIDs,
	})
}

func (r *RepositoryImpl) ClearCart(ctx context.Context, cartID int64) error {
	return r.sqlc.ClearCart(ctx, cartID)
}
