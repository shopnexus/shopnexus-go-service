package account

// func (s *AccountServer) GetCart(ctx context.Context, req *connect.Request[accountv1.GetCartRequest]) (*connect.Response[accountv1.GetCartResponse], error) {
// 	cart, err := s.service.GetCart(ctx, req.Msg.GetUserId())
// 	if err != nil {
// 		return nil, err
// 	}

// 	items := make([]*common.ItemQuantity, 0, len(cart.ProductModels))
// 	for _, item := range cart.ProductModels {
// 		items = append(items, &common.ItemQuantity{
// 			ItemId:   item.GetID(),
// 			Quantity: item.GetQuantity(),
// 		})
// 	}

// 	return connect.NewResponse(&accountv1.GetCartResponse{
// 		Items: items,
// 	}), nil
// }

// func (s *AccountServer) AddCartItem(ctx context.Context, req *connect.Request[accountv1.AddCartItemRequest]) (*connect.Response[accountv1.AddCartItemResponse], error) {
// 	quantity, err := s.service.AddCartItem(ctx, req.Msg.GetUserId(), model.ItemOnCart{
// 		ItemQuantityBase: model.ItemQuantityBase[int64]{
// 			ItemID:   req.Msg.GetProductModelId(),
// 			Quantity: req.Msg.GetQuantity(),
// 		},
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return connect.NewResponse(&accountv1.AddCartItemResponse{}), nil
// }

// func (s *AccountServer) UpdateCartItem(ctx context.Context, req *connect.Request[accountv1.UpdateCartItemRequest]) (*connect.Response[accountv1.UpdateCartItemResponse], error) {
// 	quantity, err := s.service.UpdateCartItem(ctx, req.Msg.GetUserId(), model.ItemOnCart{
// 		ItemQuantityBase: model.ItemQuantityBase[int64]{
// 			ItemID:   req.Msg.GetProductModelId(),
// 			Quantity: req.Msg.GetQuantity(),
// 		},
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return connect.NewResponse(&accountv1.UpdateCartItemResponse{}), nil
// }
