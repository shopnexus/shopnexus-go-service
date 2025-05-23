package account

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/storage/pgx/repository"
)

type GetAddressParams struct {
	ID     int64
	UserID *int64
}

func (s *AccountService) GetAddress(ctx context.Context, params GetAddressParams) (model.Address, error) {
	return s.repo.GetAddress(ctx, repository.GetAddressParams{
		ID:     params.ID,
		UserID: params.UserID,
	})
}

type ListAddressesParams = repository.ListAddressesParams

func (s *AccountService) ListAddresses(ctx context.Context, params ListAddressesParams) (result model.PaginateResult[model.Address], err error) {
	total, err := s.repo.CountAddresses(ctx, params)
	if err != nil {
		return result, err
	}

	addresses, err := s.repo.ListAddresses(ctx, params)
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

func (s *AccountService) CreateAddress(ctx context.Context, params CreateAddressParams) (model.Address, error) {
	return s.repo.CreateAddress(ctx, model.Address{
		UserID:   params.UserID,
		FullName: params.FullName,
		Phone:    params.Phone,
		Address:  params.Address,
		City:     params.City,
		Province: params.Province,
		Country:  params.Country,
	})
}

type UpdateAddressParams = repository.UpdateAddressParams

func (s *AccountService) UpdateAddress(ctx context.Context, params UpdateAddressParams) (model.Address, error) {
	return s.repo.UpdateAddress(ctx, params)
}

type DeleteAddressParams struct {
	ID     int64
	UserID *int64
}

func (s *AccountService) DeleteAddress(ctx context.Context, params DeleteAddressParams) error {
	return s.repo.DeleteAddress(ctx, repository.DeleteAddressParams{
		ID:     params.ID,
		UserID: params.UserID,
	})
}
