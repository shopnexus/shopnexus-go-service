package account

import (
	"context"
	"fmt"
	"shopnexus-go-service/internal/model"
	repository "shopnexus-go-service/internal/repository"
	"shopnexus-go-service/internal/util"

	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	repo *repository.Repository
}

type AccountServiceInterface interface {
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
}

func NewAccountService(repo *repository.Repository) *AccountService {
	return &AccountService{
		repo: repo,
	}
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
	Username *string
	Email    *string
	Phone    *string
	Role     model.Role
	Password string
}

func (s *AccountService) FindAccount(ctx context.Context, params FindAccountParams) (account model.Account, err error) {
	if params.Username == nil && params.Email == nil && params.Phone == nil {
		return nil, model.ErrInvalidCreds
	}

	switch params.Role {
	case model.RoleAdmin:
		account, err = s.repo.GetAccount(ctx, model.AccountAdmin{
			AccountBase: model.AccountBase{
				Username: util.DerefDefault(params.Username, ""),
			},
		})
	case model.RoleUser:
		account, err = s.repo.GetAccount(ctx, model.AccountUser{
			AccountBase: model.AccountBase{
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

	// Check hash password
	if ok := s.CheckPassword(account.GetBase().Password, params.Password); !ok {
		return nil, model.ErrWrongPassword
	}

	return account, nil
}

type LoginParams = FindAccountParams

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
			Password: params.Password,
			Role:     model.RoleAdmin,
		})
	case model.RoleUser:
		account, err = s.FindAccount(ctx, FindAccountParams{
			Username: params.Username,
			Email:    params.Email,
			Phone:    params.Phone,
			Password: params.Password,
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
	Role      model.Role
}

func (s *AccountService) GetPermissions(ctx context.Context, params GetPermissionsParams) ([]model.Permission, error) {
	switch params.Role {
	case model.RoleAdmin:
		return model.GetAllPermissions(), nil
	case model.RoleStaff:
		return s.repo.GetPermissions(ctx, repository.GetPermissionsParams{
			AccountID: params.AccountID,
			Role:      params.Role,
		})
	case model.RoleUser:
		return nil, fmt.Errorf("no permissions")
	default:
		return nil, fmt.Errorf("unknown role: %s", params.Role)
	}
}
