package account

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service/account"

	common_grpc "shopnexus-go-service/internal/transport/connect/handler/common"
	"shopnexus-go-service/internal/transport/connect/interceptor/auth"

	"connectrpc.com/connect"
	accountv1 "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/account/v1"
)

func (s *ImplementedAccountServiceHandler) GetAddress(ctx context.Context, req *connect.Request[accountv1.GetAddressRequest]) (*connect.Response[accountv1.GetAddressResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	data, err := s.service.GetAddress(ctx, account.GetAddressParams{
		ID:     req.Msg.Id,
		UserID: &claims.UserID,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&accountv1.GetAddressResponse{
		Data: modelToAddressEntity(data),
	}), nil
}

func (s *ImplementedAccountServiceHandler) ListAddresses(ctx context.Context, req *connect.Request[accountv1.ListAddressesRequest]) (*connect.Response[accountv1.ListAddressesResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	data, err := s.service.ListAddresses(ctx, account.ListAddressesParams{
		PaginationParams: model.PaginationParams{
			Page:  req.Msg.GetPagination().GetPage(),
			Limit: req.Msg.GetPagination().GetLimit(),
		},
		UserID:   &claims.UserID,
		FullName: req.Msg.FullName,
		Phone:    req.Msg.Phone,
		Address:  req.Msg.Address,
		City:     req.Msg.City,
		Province: req.Msg.Province,
		Country:  req.Msg.Country,
	})
	if err != nil {
		return nil, err
	}

	var addresses []*accountv1.AddressEntity
	for _, d := range data.Data {
		addresses = append(addresses, modelToAddressEntity(d))
	}

	return connect.NewResponse(&accountv1.ListAddressesResponse{
		Data:       addresses,
		Pagination: common_grpc.ToProtoPaginationResponse(data),
	}), nil
}

func (s *ImplementedAccountServiceHandler) CreateAddress(ctx context.Context, req *connect.Request[accountv1.CreateAddressRequest]) (*connect.Response[accountv1.CreateAddressResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	data, err := s.service.CreateAddress(ctx, account.CreateAddressParams{
		UserID:   claims.UserID,
		FullName: req.Msg.FullName,
		Phone:    req.Msg.Phone,
		Address:  req.Msg.Address,
		City:     req.Msg.City,
		Province: req.Msg.Province,
		Country:  req.Msg.Country,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&accountv1.CreateAddressResponse{
		Data: modelToAddressEntity(data),
	}), nil
}

func (s *ImplementedAccountServiceHandler) UpdateAddress(ctx context.Context, req *connect.Request[accountv1.UpdateAddressRequest]) (*connect.Response[accountv1.UpdateAddressResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	_, err = s.service.UpdateAddress(ctx, account.UpdateAddressParams{
		ID:       req.Msg.Id,
		UserID:   &claims.UserID,
		FullName: req.Msg.FullName,
		Phone:    req.Msg.Phone,
		Address:  req.Msg.Address,
		City:     req.Msg.City,
		Province: req.Msg.Province,
		Country:  req.Msg.Country,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&accountv1.UpdateAddressResponse{}), nil
}

func (s *ImplementedAccountServiceHandler) DeleteAddress(ctx context.Context, req *connect.Request[accountv1.DeleteAddressRequest]) (*connect.Response[accountv1.DeleteAddressResponse], error) {
	claims, err := auth.GetClaims(req)
	if err != nil {
		return nil, err
	}

	err = s.service.DeleteAddress(ctx, account.DeleteAddressParams{
		ID:     req.Msg.Id,
		UserID: &claims.UserID,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&accountv1.DeleteAddressResponse{}), nil
}

func modelToAddressEntity(data model.Address) *accountv1.AddressEntity {
	return &accountv1.AddressEntity{
		Id:          data.ID,
		UserId:      data.UserID,
		FullName:    data.FullName,
		Phone:       data.Phone,
		Address:     data.Address,
		City:        data.City,
		Province:    data.Province,
		Country:     data.Country,
		DateCreated: data.DateCreated,
		DateUpdated: data.DateUpdated,
	}
}
