package storage

import (
	"context"
	"fmt"
	"shopnexus-go-service/gen/sqlc"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/utils/bytes"

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
		Role:     model.Role(baseRow.Role),
		Avatar:   PgtypeToPtr[string](baseRow.AvatarUrl),
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
			Role:     model.Role(userRow.Role),
			Avatar:   PgtypeToPtr[string](userRow.AvatarUrl),
		},
		Email:            userRow.Email,
		Phone:            userRow.Phone,
		Gender:           model.Gender(userRow.Gender),
		FullName:         userRow.FullName,
		DefaultAddressID: PgtypeToPtr[int64](userRow.DefaultAddressID),
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
			Role:     model.Role(adminRow.Role),
			Avatar:   PgtypeToPtr[string](adminRow.AvatarUrl),
		},
	}, nil
}

func (r *ServiceImpl) GetAccount(ctx context.Context, find model.Account) (account model.Account, err error) {
	switch find := find.(type) {
	case model.AccountBase:
		// Search by base account info (id && type)
		accountType := find.Role

		if accountType == "" {
			accountBase, err := r.GetAccountBase(ctx, find.ID)
			if err != nil {
				return nil, err
			}

			accountType = accountBase.Role
		}

		switch accountType {
		case model.RoleAdmin:
			admin, err := r.GetAccountAdmin(ctx, GetAccountAdminParams{
				AccountID: &find.ID,
				Username:  &find.Username,
			})
			if err != nil {
				return nil, err
			}
			account = admin
		case model.RoleUser:
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
				Role:     model.RoleUser,
				Avatar:   account.Avatar,
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
				Role:     model.RoleAdmin,
				Avatar:   account.Avatar,
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
		ID:                   params.ID,
		Username:             *PtrToPgtype(&pgtype.Text{}, params.Username),
		Password:             *PtrToPgtype(&pgtype.Text{}, params.Password),
		NullCustomPermission: params.NullCustomPermission,
		CustomPermission:     *PtrToPgtype(&pgtype.Bits{}, params.CustomPermission),
		AvatarUrl:            *PtrToPgtype(&pgtype.Text{}, params.AvatarURL),
	})
	if err != nil {
		return model.AccountBase{}, err
	}

	return model.AccountBase{
		ID:       row.ID,
		Role:     model.Role(row.Role),
		Username: row.Username,
		Password: row.Password,
		Avatar:   PgtypeToPtr[string](row.AvatarUrl),
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
	Role      model.Role
}

func (r *ServiceImpl) GetPermissions(ctx context.Context, params GetPermissionsParams) ([]model.Permission, error) {
	permissionBits, err := r.sqlc.GetRolePermissions(ctx, string(params.Role))
	if err != nil {
		return nil, err
	}

	if !permissionBits.Valid {
		return nil, fmt.Errorf("role %s does not have any permissions", params.Role)
	}

	// Get custom permissions
	customPermissions, err := r.sqlc.GetCustomPermissions(ctx, params.AccountID)
	if err != nil {
		return nil, err
	}
	if customPermissions.Valid {
		// Merge custom permissions with role permissions
		permissionBits.Bytes = bytes.OrByteSlices(permissionBits.Bytes, customPermissions.Bytes)
	}

	// Convert bit array to permissions slice
	var permissions []model.Permission

	for i := permissionBits.Len - 1; i >= 0; i-- {
		byteIndex := i / 8                                            // Find which byte contains the bit
		bitIndex := 7 - (i % 8)                                       // Find the position within the byte
		bitValue := (permissionBits.Bytes[byteIndex] >> bitIndex) & 1 // Extract bit value
		if bitValue == 1 {
			permissions = append(permissions, model.Permission(permissionBits.Len-i))
		}
	}

	return permissions, nil
}
