package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	pgxutil "shopnexus-go-service/internal/db/pgx"
	"shopnexus-go-service/internal/model"

	"github.com/jackc/pgx/v5/pgtype"
)

type GetAddressParams struct {
	ID     int64
	UserID *int64
}

func (r *RepositoryImpl) GetAddress(ctx context.Context, params GetAddressParams) (model.Address, error) {
	row, err := r.sqlc.GetAddress(ctx, sqlc.GetAddressParams{
		ID:     params.ID,
		UserID: *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UserID),
	})
	if err != nil {
		return model.Address{}, err
	}

	return model.Address{
		ID:          row.ID,
		UserID:      row.UserID,
		FullName:    row.FullName,
		Phone:       row.Phone,
		Address:     row.Address,
		City:        row.City,
		Province:    row.Province,
		Country:     row.Country,
		DateCreated: row.DateCreated.Time.UnixMilli(),
		DateUpdated: row.DateUpdated.Time.UnixMilli(),
	}, nil
}

type ListAddressesParams struct {
	model.PaginationParams
	UserID   *int64
	FullName *string
	Phone    *string
	Address  *string
	City     *string
	Province *string
	Country  *string
}

func (r *RepositoryImpl) CountAddresses(ctx context.Context, params ListAddressesParams) (int64, error) {
	return r.sqlc.CountAddresses(ctx, sqlc.CountAddressesParams{
		UserID:   *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UserID),
		FullName: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.FullName),
		Phone:    *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Phone),
		Address:  *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Address),
		City:     *pgxutil.PtrToPgtype(&pgtype.Text{}, params.City),
		Province: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Province),
		Country:  *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Country),
	})
}

func (r *RepositoryImpl) ListAddresses(ctx context.Context, params ListAddressesParams) ([]model.Address, error) {
	addresses, err := r.sqlc.ListAddresses(ctx, sqlc.ListAddressesParams{
		Limit:    params.Limit,
		Offset:   params.Offset(),
		UserID:   *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UserID),
		FullName: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.FullName),
		Phone:    *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Phone),
		Address:  *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Address),
		City:     *pgxutil.PtrToPgtype(&pgtype.Text{}, params.City),
		Province: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Province),
		Country:  *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Country),
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.Address, len(addresses))
	for i, address := range addresses {
		result[i] = model.Address{
			ID:          address.ID,
			UserID:      address.UserID,
			FullName:    address.FullName,
			Phone:       address.Phone,
			Address:     address.Address,
			City:        address.City,
			Province:    address.Province,
			Country:     address.Country,
			DateCreated: address.DateCreated.Time.UnixMilli(),
			DateUpdated: address.DateUpdated.Time.UnixMilli(),
		}
	}

	return result, nil
}

func (r *RepositoryImpl) CreateAddress(ctx context.Context, address model.Address) (model.Address, error) {
	row, err := r.sqlc.CreateAddress(ctx, sqlc.CreateAddressParams{
		UserID:   address.UserID,
		FullName: address.FullName,
		Phone:    address.Phone,
		Address:  address.Address,
		City:     address.City,
		Province: address.Province,
		Country:  address.Country,
	})
	if err != nil {
		return model.Address{}, err
	}

	return model.Address{
		ID:          row.ID,
		UserID:      row.UserID,
		FullName:    row.FullName,
		Phone:       row.Phone,
		Address:     row.Address,
		City:        row.City,
		Province:    row.Province,
		Country:     row.Country,
		DateCreated: row.DateCreated.Time.UnixMilli(),
		DateUpdated: row.DateUpdated.Time.UnixMilli(),
	}, nil
}

type UpdateAddressParams struct {
	ID       int64
	UserID   *int64
	FullName *string
	Phone    *string
	Address  *string
	City     *string
	Province *string
	Country  *string
}

func (r *RepositoryImpl) UpdateAddress(ctx context.Context, params UpdateAddressParams) (model.Address, error) {
	row, err := r.sqlc.UpdateAddress(ctx, sqlc.UpdateAddressParams{
		ID:       params.ID,
		UserID:   *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UserID),
		FullName: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.FullName),
		Phone:    *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Phone),
		Address:  *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Address),
		City:     *pgxutil.PtrToPgtype(&pgtype.Text{}, params.City),
		Province: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Province),
		Country:  *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Country),
	})
	if err != nil {
		return model.Address{}, err
	}

	return model.Address{
		ID:          row.ID,
		UserID:      row.UserID,
		FullName:    row.FullName,
		Phone:       row.Phone,
		Address:     row.Address,
		City:        row.City,
		Province:    row.Province,
		Country:     row.Country,
		DateCreated: row.DateCreated.Time.UnixMilli(),
		DateUpdated: row.DateUpdated.Time.UnixMilli(),
	}, nil
}

type DeleteAddressParams struct {
	ID     int64
	UserID *int64
}

func (r *RepositoryImpl) DeleteAddress(ctx context.Context, params DeleteAddressParams) error {
	_, err := r.sqlc.DeleteAddress(ctx, sqlc.DeleteAddressParams{
		ID:     params.ID,
		UserID: *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UserID),
	})
	return err
}
