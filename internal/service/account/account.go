package account

import (
	"context"
	"fmt"
	"shopnexus-go-service/internal/model"
	repository "shopnexus-go-service/internal/repository"
	"shopnexus-go-service/internal/util"
	"slices"

	"golang.org/x/crypto/bcrypt"
)

var _ AccountServiceInterface = (*AccountService)(nil)

type AccountService struct {
	repo repository.Repository
}

type AccountServiceInterface interface {
	WithTx(txRepo repository.Repository) *AccountService

	// Account
	CheckPassword(hashedPassword, password string) bool
	CreateHash(password string) (string, error)
	FindAccount(ctx context.Context, params FindAccountParams) (model.Account, error)
	Login(ctx context.Context, params LoginParams) (string, error)
	Register(ctx context.Context, account model.Account) (string, error)

	// Cart
	GetCart(ctx context.Context, userID int64) (model.Cart, error)
	AddCartItem(ctx context.Context, params AddCartItemParams) (int64, error)
	UpdateCartItem(ctx context.Context, params UpdateCartItemParams) (int64, error)
	ClearCart(ctx context.Context, userID int64) error

	// Address
	GetAddress(ctx context.Context, params GetAddressParams) (model.Address, error)
	ListAddresses(ctx context.Context, params ListAddressesParams) (model.PaginateResult[model.Address], error)
	CreateAddress(ctx context.Context, params model.Address) (model.Address, error)
	UpdateAddress(ctx context.Context, params UpdateAddressParams) (model.Address, error)
	DeleteAddress(ctx context.Context, params DeleteAddressParams) error

	// Permission
	GetPermissions(ctx context.Context, params GetPermissionsParams) ([]model.Permission, error)
	HasPermission(ctx context.Context, params HasPermissionParams) (bool, error)
}

func NewAccountService(repo repository.Repository) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

func (s *AccountService) WithTx(txRepo repository.Repository) *AccountService {
	return NewAccountService(txRepo)
}

func (s *AccountService) CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// CreateHash Create hash and add some salt :P
func (s *AccountService) CreateHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

type FindAccountParams struct {
	Role     model.Role
	UserID   *int64
	Username *string
	Email    *string
	Phone    *string
	Password *string
}

func (s *AccountService) FindAccount(ctx context.Context, params FindAccountParams) (account model.Account, err error) {
	if params.Username == nil && params.Email == nil && params.Phone == nil && params.UserID == nil {
		return nil, model.ErrInvalidCreds
	}

	switch params.Role {
	case model.RoleAdmin:
		account, err = s.repo.GetAccount(ctx, model.AccountAdmin{
			AccountBase: model.AccountBase{
				ID:       util.DerefDefault(params.UserID, 0),
				Username: util.DerefDefault(params.Username, ""),
			},
		})
	case model.RoleUser:
		account, err = s.repo.GetAccount(ctx, model.AccountUser{
			AccountBase: model.AccountBase{
				ID:       util.DerefDefault(params.UserID, 0),
				Username: util.DerefDefault(params.Username, ""),
			},
			Email: util.DerefDefault(params.Email, ""),
			Phone: util.DerefDefault(params.Phone, ""),
		})
	default:
		return nil, fmt.Errorf("unknown account role: %s", params.Role)
	}
	if err != nil {
		return nil, err
	}

	if params.Password != nil {
		// Check hash password
		if ok := s.CheckPassword(account.GetBase().Password, *params.Password); !ok {
			return nil, model.ErrWrongPassword
		}
	}

	return account, nil
}

type LoginParams struct {
	Role     model.Role
	Password string
	Username *string
	Email    *string
	Phone    *string
}

func (s *AccountService) Login(ctx context.Context, params LoginParams) (string, error) {
	if params.Username == nil && params.Email == nil && params.Phone == nil {
		return "", model.ErrInvalidCreds
	}

	var (
		err     error
		account model.Account
	)

	switch params.Role {
	case model.RoleAdmin:
		account, err = s.FindAccount(ctx, FindAccountParams{
			Username: params.Username,
			Password: &params.Password,
			Role:     model.RoleAdmin,
		})
	case model.RoleUser:
		account, err = s.FindAccount(ctx, FindAccountParams{
			Username: params.Username,
			Email:    params.Email,
			Phone:    params.Phone,
			Password: &params.Password,
			Role:     model.RoleUser,
		})
	}

	if err != nil {
		return "", err
	}

	return util.GenerateAccessToken(account)
}

func (s *AccountService) Register(ctx context.Context, account model.Account) (string, error) {
	txRepo, err := s.repo.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer txRepo.Rollback(ctx)

	newAccount, err := txRepo.CreateAccount(ctx, account)
	if err != nil {
		return "", err
	}

	if newAccount.GetBase().Role == model.RoleUser {
		// Create cart
		err = txRepo.CreateCart(ctx, newAccount.GetBase().ID)
		if err != nil {
			return "", err
		}
	}

	if err = txRepo.Commit(ctx); err != nil {
		return "", err
	}

	return util.GenerateAccessToken(newAccount)
}

type GetPermissionsParams struct {
	AccountID int64
	Role      *model.Role
}

func (s *AccountService) GetPermissions(ctx context.Context, params GetPermissionsParams) ([]model.Permission, error) {
	var role model.Role

	if params.Role == nil {
		account, err := s.repo.GetAccount(ctx, model.AccountBase{
			ID: params.AccountID,
		})
		if err != nil {
			return nil, err
		}
		role = account.GetBase().Role
	} else {
		role = *params.Role
	}

	switch role {
	case model.RoleAdmin:
		return model.GetAllPermissions(), nil
	case model.RoleStaff:
		return s.repo.GetPermissions(ctx, repository.GetPermissionsParams{
			AccountID: params.AccountID,
			Role:      role,
		})
	case model.RoleUser:
		return nil, model.ErrPermissionDenied
	default:
		return nil, fmt.Errorf("unknown role: %s", role)
	}
}

type HasPermissionParams struct {
	AccountID   int64
	Role        *model.Role
	Permissions []model.Permission
}

func (s *AccountService) HasPermission(ctx context.Context, params HasPermissionParams) (bool, error) {
	permissions, err := s.GetPermissions(ctx, GetPermissionsParams{
		AccountID: params.AccountID,
		Role:      params.Role,
	})
	if err != nil {
		return false, err
	}

	if len(permissions) == 0 {
		return false, model.ErrPermissionDenied
	}

	for _, permission := range params.Permissions {
		if !slices.Contains(permissions, permission) {
			return false, nil
		}
	}

	return true, nil
}

type UpdateAccountParams struct {
	ID                   int64
	Username             *string
	Password             *string
	NullCustomPermission bool
	CustomPermission     *string
}

func (s *AccountService) UpdateAccount(ctx context.Context, params UpdateAccountParams) (model.AccountBase, error) {
	return s.repo.UpdateAccount(ctx, repository.UpdateAccountParams{
		ID:                   params.ID,
		Username:             params.Username,
		Password:             params.Password,
		NullCustomPermission: params.NullCustomPermission,
		CustomPermission:     params.CustomPermission,
	})
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

func (s *AccountService) UpdateAccountUser(ctx context.Context, params UpdateAccountUserParams) (model.AccountUser, error) {
	return s.repo.UpdateAccountUser(ctx, repository.UpdateAccountUserParams{
		ID:                   params.ID,
		Email:                params.Email,
		Phone:                params.Phone,
		Gender:               params.Gender,
		FullName:             params.FullName,
		DefaultAddressID:     params.DefaultAddressID,
		NullDefaultAddressID: params.NullDefaultAddressID,
	})
}
