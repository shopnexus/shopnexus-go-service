package account

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service/storage"
)

type GetAddressParams struct {
	ID     int64
	UserID *int64
}

func (s *ServiceImpl) GetAddress(ctx context.Context, params GetAddressParams) (model.Address, error) {
	return s.storage.GetAddress(ctx, storage.GetAddressParams{
		ID:     params.ID,
		UserID: params.UserID,
	})
}

type ListAddressesParams = storage.ListAddressesParams

func (s *ServiceImpl) ListAddresses(ctx context.Context, params ListAddressesParams) (result model.PaginateResult[model.Address], err error) {
	total, err := s.storage.CountAddresses(ctx, params)
	if err != nil {
		return result, err
	}

	addresses, err := s.storage.ListAddresses(ctx, params)
	if err != nil {
		return result, err
	}

	return model.PaginateResult[model.Address]{
		Data:       addresses,
		Limit:      params.Limit,
		Page:       params.Page,
		Total:      total,
		NextPage:   params.NextPage(total),
		NextCursor: nil,
	}, nil
}

type CreateAddressParams struct {
	UserID   int64
	FullName string
	Phone    string
	Address  string
	City     string
	Province string
	Country  string
}

func (s *ServiceImpl) CreateAddress(ctx context.Context, params CreateAddressParams) (model.Address, error) {
	return s.storage.CreateAddress(ctx, model.Address{
		UserID:   params.UserID,
		FullName: params.FullName,
		Phone:    params.Phone,
		Address:  params.Address,
		City:     params.City,
		Province: params.Province,
		Country:  params.Country,
	})
}

type UpdateAddressParams = storage.UpdateAddressParams

func (s *ServiceImpl) UpdateAddress(ctx context.Context, params UpdateAddressParams) (model.Address, error) {
	return s.storage.UpdateAddress(ctx, params)
}

type DeleteAddressParams struct {
	ID     int64
	UserID *int64
}

func (s *ServiceImpl) DeleteAddress(ctx context.Context, params DeleteAddressParams) error {
	return s.storage.DeleteAddress(ctx, storage.DeleteAddressParams{
		ID:     params.ID,
		UserID: params.UserID,
	})
}
