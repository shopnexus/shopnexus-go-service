package account

import (
	"context"
	"shopnexus-go-service/internal/service/account"
	"shopnexus-go-service/internal/transport/connect/interceptor/auth"

	"connectrpc.com/connect"
	accountv1 "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/account/v1"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/common"
)

func (s *ImplementedAccountServiceHandler) GetCart(ctx context.Context, req *connect.Request[accountv1.GetCartRequest]) (*connect.Response[accountv1.GetCartResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	cart, err := s.service.GetCart(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	items := make([]*common.ItemQuantityInt64, 0, len(cart.Products))
	for _, item := range cart.Products {
		items = append(items, &common.ItemQuantityInt64{
			ItemId:   item.GetID(),
			Quantity: item.GetQuantity(),
		})
	}

	return connect.NewResponse(&accountv1.GetCartResponse{
		Items: items,
	}), nil
}

func (s *ImplementedAccountServiceHandler) AddCartItem(ctx context.Context, req *connect.Request[accountv1.AddCartItemRequest]) (*connect.Response[accountv1.AddCartItemResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	for _, item := range req.Msg.GetItems() {
		_, err := s.service.AddCartItem(ctx, account.AddCartItemParams{
			UserID:    claims.UserID,
			ProductID: item.GetItemId(),
			Quantity:  item.GetQuantity(),
		})
		if err != nil {
			return nil, err
		}
	}

	return connect.NewResponse(&accountv1.AddCartItemResponse{}), nil
}

func (s *ImplementedAccountServiceHandler) UpdateCartItem(ctx context.Context, req *connect.Request[accountv1.UpdateCartItemRequest]) (*connect.Response[accountv1.UpdateCartItemResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	for _, item := range req.Msg.GetItems() {
		_, err := s.service.UpdateCartItem(ctx, account.UpdateCartItemParams{
			UserID:    claims.UserID,
			ProductID: item.GetItemId(),
			Quantity:  item.GetQuantity(),
		})
		if err != nil {
			return nil, err
		}
	}

	return connect.NewResponse(&accountv1.UpdateCartItemResponse{}), nil
}

func (s *ImplementedAccountServiceHandler) ClearCart(ctx context.Context, req *connect.Request[accountv1.ClearCartRequest]) (*connect.Response[accountv1.ClearCartResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	err = s.service.ClearCart(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&accountv1.ClearCartResponse{}), nil
}
