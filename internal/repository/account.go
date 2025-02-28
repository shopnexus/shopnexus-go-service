package repository

import (
	"context"
	"fmt"
	"shopnexus-go-service/gen/sqlc"
	pgxutil "shopnexus-go-service/internal/db/pgx"
	"shopnexus-go-service/internal/model"

	"github.com/jackc/pgx/v5/pgtype"
)

func (r *Repository) GetAccountBase(ctx context.Context, accountID int64) (model.AccountBase, error) {
	baseRow, err := r.sqlc.GetAccountBase(ctx, accountID)
	if err != nil {
		return model.AccountBase{}, err
	}

	return model.AccountBase{
		ID:       baseRow.ID,
		Username: baseRow.Username,
		Password: baseRow.Password,
		Role:     model.Role(baseRow.Role),
	}, nil
}

type GetAccountUserParams struct {
	AccountID *int64
	Username  *string
	Email     *string
	Phone     *string
}

func (r *Repository) GetAccountUser(ctx context.Context, params GetAccountUserParams) (model.AccountUser, error) {
	userRow, err := r.sqlc.GetAccountUser(ctx, sqlc.GetAccountUserParams{
		ID:       *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AccountID),
		Username: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Username),
		Email:    *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Email),
		Phone:    *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Phone),
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
		},
		Email:            userRow.Email,
		Phone:            userRow.Phone,
		Gender:           model.Gender(userRow.Gender),
		FullName:         userRow.FullName,
		DefaultAddressID: pgxutil.PgtypeToPtr[int64](userRow.DefaultAddressID),
	}, nil
}

type GetAccountAdminParams struct {
	AccountID *int64
	Username  *string
}

func (r *Repository) GetAccountAdmin(ctx context.Context, params GetAccountAdminParams) (model.AccountAdmin, error) {
	adminRow, err := r.sqlc.GetAccountAdmin(ctx, sqlc.GetAccountAdminParams{
		ID:       *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.AccountID),
		Username: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Username),
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
		},
	}, nil
}

func (r *Repository) GetAccount(ctx context.Context, find model.Account) (account model.Account, err error) {
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

func (r *Repository) CreateAccount(ctx context.Context, account model.Account) (model.Account, error) {
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
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown account type: %T", account)
	}
}

func (r *Repository) GetCart(ctx context.Context, cartID int64) (model.Cart, error) {
	cartRow, err := r.sqlc.GetCart(ctx, cartID)
	if err != nil {
		return model.Cart{}, err
	}

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
		ID:            cartRow,
		ProductModels: items,
	}, nil
}

func (r *Repository) CreateCart(ctx context.Context, userID int64) error {
	return r.sqlc.CreateCart(ctx, userID)
}

type AddCartItemParams struct {
	CartID         int64
	ProductModelID int64
	Quantity       int64
}

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

type RemoveCartItemParams struct {
	CartID         int64
	ProductModelID int64
}

func (r *Repository) RemoveCartItem(ctx context.Context, params RemoveCartItemParams) error {
	return r.sqlc.RemoveCartItem(ctx, sqlc.RemoveCartItemParams{
		CartID:         params.CartID,
		ProductModelID: params.ProductModelID,
	})
}

func (r *Repository) ClearCart(ctx context.Context, cartID int64) error {
	return r.sqlc.ClearCart(ctx, cartID)
}
