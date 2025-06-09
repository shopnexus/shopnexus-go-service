package storage

import (
	"context"
	"fmt"
	"shopnexus-go-service/gen/sqlc"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/utils/slices"

	"github.com/jackc/pgx/v5/pgtype"
)

func (r *ServiceImpl) GetAccountBase(ctx context.Context, accountID int64) (model.AccountBase, error) {
	baseRow, err := r.sqlc.GetAccountBase(ctx, accountID)
	if err != nil {
		return model.AccountBase{}, err
	}

	return model.AccountBase{
		ID:       baseRow.ID,
		Username: baseRow.Username,
		Password: baseRow.Password,
		Type:     model.AccountType(baseRow.Type),
	}, nil
}

type GetAccountUserParams struct {
	AccountID *int64
	Username  *string
	Email     *string
	Phone     *string
}

func (r *ServiceImpl) GetAccountUser(ctx context.Context, params GetAccountUserParams) (model.AccountUser, error) {
	userRow, err := r.sqlc.GetAccountUser(ctx, sqlc.GetAccountUserParams{
		ID:       *PtrToPgtype(&pgtype.Int8{}, params.AccountID),
		Username: *PtrToPgtype(&pgtype.Text{}, params.Username),
		Email:    *PtrToPgtype(&pgtype.Text{}, params.Email),
		Phone:    *PtrToPgtype(&pgtype.Text{}, params.Phone),
	})
	if err != nil {
		return model.AccountUser{}, err
	}

	return model.AccountUser{
		AccountBase: model.AccountBase{
			ID:       userRow.ID,
			Username: userRow.Username,
			Password: userRow.Password,
			Type:     model.AccountTypeUser,
		},
		Email:            userRow.Email,
		Phone:            userRow.Phone,
		Gender:           model.Gender(userRow.Gender),
		FullName:         userRow.FullName,
		DefaultAddressID: PgtypeToPtr[int64](userRow.DefaultAddressID),
		AvatarURL:        PgtypeToPtr[string](userRow.AvatarUrl),
	}, nil
}

type GetAccountAdminParams struct {
	AccountID *int64
	Username  *string
}

func (r *ServiceImpl) GetAccountAdmin(ctx context.Context, params GetAccountAdminParams) (model.AccountAdmin, error) {
	adminRow, err := r.sqlc.GetAccountAdmin(ctx, sqlc.GetAccountAdminParams{
		ID:       *PtrToPgtype(&pgtype.Int8{}, params.AccountID),
		Username: *PtrToPgtype(&pgtype.Text{}, params.Username),
	})
	if err != nil {
		return model.AccountAdmin{}, err
	}

	return model.AccountAdmin{
		AccountBase: model.AccountBase{
			ID:       adminRow.ID,
			Username: adminRow.Username,
			Password: adminRow.Password,
			Type:     model.AccountTypeAdmin,
		},
	}, nil
}

func (r *ServiceImpl) GetAccount(ctx context.Context, find model.Account) (account model.Account, err error) {
	switch find := find.(type) {
	case model.AccountBase:
		// Search by base account info (id && type)
		accountType := find.Type

		if accountType == "" {
			accountBase, err := r.GetAccountBase(ctx, find.ID)
			if err != nil {
				return nil, err
			}

			accountType = accountBase.Type
		}

		switch accountType {
		case model.AccountTypeAdmin:
			admin, err := r.GetAccountAdmin(ctx, GetAccountAdminParams{
				AccountID: &find.ID,
				Username:  &find.Username,
			})
			if err != nil {
				return nil, err
			}
			account = admin
		case model.AccountTypeUser:
			user, err := r.GetAccountUser(ctx, GetAccountUserParams{
				AccountID: &find.ID,
			})
			if err != nil {
				return nil, err
			}
			account = user
		default:
			return nil, fmt.Errorf("unknown account role: %s", accountType)
		}

	case model.AccountUser:
		// Search by user info
		account, err = r.GetAccountUser(ctx, GetAccountUserParams{
			AccountID: &find.ID,
			Username:  &find.Username,
			Email:     &find.Email,
			Phone:     &find.Phone,
		})

	case model.AccountAdmin:
		// Search by admin info
		account, err = r.GetAccountAdmin(ctx, GetAccountAdminParams{
			AccountID: &find.ID,
			Username:  &find.Username,
		})
	default:
		return nil, fmt.Errorf("unknown account type: %T", find)
	}
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *ServiceImpl) CreateAccount(ctx context.Context, account model.Account) (model.Account, error) {
	switch account := account.(type) {
	case model.AccountUser:
		id, err := r.sqlc.CreateAccountUser(ctx, sqlc.CreateAccountUserParams{
			Username: account.Username,
			Password: account.Password,
			Email:    account.Email,
			Phone:    account.Phone,
			Gender:   sqlc.AccountGender(account.Gender),
			FullName: account.FullName,
		})
		if err != nil {
			return nil, err
		}

		return model.AccountUser{
			AccountBase: model.AccountBase{
				ID:       id,
				Username: account.Username,
				Password: account.Password,
				Type:     model.AccountTypeUser,
			},
			Email:    account.Email,
			Phone:    account.Phone,
			Gender:   account.Gender,
			FullName: account.FullName,
		}, nil
	case model.AccountAdmin:
		id, err := r.sqlc.CreateAccountAdmin(ctx, sqlc.CreateAccountAdminParams{
			Username: account.Username,
			Password: account.Password,
		})
		if err != nil {
			return nil, err
		}

		return model.AccountAdmin{
			AccountBase: model.AccountBase{
				ID:       id,
				Username: account.Username,
				Password: account.Password,
				Type:     model.AccountTypeAdmin,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown account type: %T", account)
	}
}

type UpdateAccountParams struct {
	ID                   int64
	Username             *string
	Password             *string
	NullCustomPermission bool
	CustomPermission     *string
	AvatarURL            *string
}

func (r *ServiceImpl) UpdateAccount(ctx context.Context, params UpdateAccountParams) (model.AccountBase, error) {
	row, err := r.sqlc.UpdateAccount(ctx, sqlc.UpdateAccountParams{
		ID:       params.ID,
		Username: *PtrToPgtype(&pgtype.Text{}, params.Username),
		Password: *PtrToPgtype(&pgtype.Text{}, params.Password),
	})
	if err != nil {
		return model.AccountBase{}, err
	}

	return model.AccountBase{
		ID:       row.ID,
		Type:     model.AccountType(row.Type),
		Username: row.Username,
		Password: row.Password,
	}, nil
}

type UpdateAccountUserParams struct {
	ID                   int64
	Email                *string
	Phone                *string
	Gender               *model.Gender
	FullName             *string
	DefaultAddressID     *int64
	NullDefaultAddressID bool
}

func (r *ServiceImpl) UpdateAccountUser(ctx context.Context, params UpdateAccountUserParams) (model.AccountUser, error) {
	row, err := r.sqlc.UpdateAccountUser(ctx, sqlc.UpdateAccountUserParams{
		ID:                   params.ID,
		Email:                *PtrToPgtype(&pgtype.Text{}, params.Email),
		Phone:                *PtrToPgtype(&pgtype.Text{}, params.Phone),
		Gender:               *PtrBrandedToPgType(&sqlc.NullAccountGender{}, params.Gender),
		FullName:             *PtrToPgtype(&pgtype.Text{}, params.FullName),
		DefaultAddressID:     *PtrToPgtype(&pgtype.Int8{}, params.DefaultAddressID),
		NullDefaultAddressID: params.NullDefaultAddressID,
	})
	if err != nil {
		return model.AccountUser{}, err
	}

	return model.AccountUser{
		Email:            row.Email,
		Phone:            row.Phone,
		Gender:           model.Gender(row.Gender),
		FullName:         row.FullName,
		DefaultAddressID: PgtypeToPtr[int64](row.DefaultAddressID),
	}, nil
}

type GetPermissionsParams struct {
	AccountID int64
	Role      model.AccountType
}

func (r *ServiceImpl) GetPermissions(ctx context.Context, params GetPermissionsParams) ([]model.Permission, error) {
	permissions, err := r.sqlc.GetAdminPermissions(ctx, params.AccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	return slices.UnsafeConvertSlice[string, model.Permission](permissions), nil
}
