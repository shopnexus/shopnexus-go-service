package account

import (
	"context"
	"errors"
	"fmt"
	"shopnexus-go-service/config"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/service/storage"
	"shopnexus-go-service/internal/utils/ptr"
	"slices"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type ServiceImpl struct {
	storage storage.Service
}

type Service interface {
	WithTx(txStorage *storage.TxStorage) (Service, error)

	// Account
	CheckPassword(hashedPassword, password string) bool
	CreateHash(password string) (string, error)
	FindAccount(ctx context.Context, params FindAccountParams) (model.Account, error)
	Login(ctx context.Context, params LoginParams) (string, error)
	Register(ctx context.Context, account model.Account) (string, error)
	UpdateAccount(ctx context.Context, params UpdateAccountParams) (model.AccountBase, error)
	UpdateAccountUser(ctx context.Context, params UpdateAccountUserParams) (model.AccountUser, error)

	// Cart
	GetCart(ctx context.Context, userID int64) (model.Cart, error)
	AddCartItem(ctx context.Context, params AddCartItemParams) (int64, error)
	UpdateCartItem(ctx context.Context, params UpdateCartItemParams) (int64, error)
	ClearCart(ctx context.Context, userID int64) error

	// Address
	GetAddress(ctx context.Context, params GetAddressParams) (model.Address, error)
	ListAddresses(ctx context.Context, params ListAddressesParams) (model.PaginateResult[model.Address], error)
	CreateAddress(ctx context.Context, params CreateAddressParams) (model.Address, error)
	UpdateAddress(ctx context.Context, params UpdateAddressParams) (model.Address, error)
	DeleteAddress(ctx context.Context, params DeleteAddressParams) error

	// Permission
	GetPermissions(ctx context.Context, params GetPermissionsParams) ([]model.Permission, error)
	HasPermission(ctx context.Context, params HasPermissionParams) (bool, error)
}

func NewService(storage storage.Service) (Service, error) {
	return &ServiceImpl{
		storage: storage,
	}, nil
}

func (s *ServiceImpl) WithTx(txStorage *storage.TxStorage) (Service, error) {
	return NewService(txStorage)
}

func (s *ServiceImpl) GenerateAccessToken(account model.Account) (string, error) {
	tokenDuration := time.Duration(config.GetConfig().App.AccessTokenDuration * int64(time.Second))

	claims := model.Claims{
		UserID: account.GetBase().ID,
		Role:   account.GetBase().Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "shopnexus",
			Subject:   strconv.Itoa(int(account.GetBase().ID)),
			Audience:  []string{"shopnexus"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := config.GetConfig().SensitiveKeys.JWTSecret

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	// test if token is valid
	// _, err = ValidateAccessToken(signedToken)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to validate token: %w", err)
	// }

	return signedToken, nil
}

func (s *ServiceImpl) CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// CreateHash Create hash and add some salt :P
func (s *ServiceImpl) CreateHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

type FindAccountParams struct {
	Role     model.AccountType
	UserID   *int64
	Username *string
	Email    *string
	Phone    *string
	Password *string
}

func (s *ServiceImpl) FindAccount(ctx context.Context, params FindAccountParams) (account model.Account, err error) {
	if params.Username == nil && params.Email == nil && params.Phone == nil && params.UserID == nil {
		return nil, model.ErrInvalidCreds
	}

	switch params.Role {
	case model.RoleAdmin:
		account, err = s.storage.GetAccount(ctx, model.AccountAdmin{
			AccountBase: model.AccountBase{
				ID:       ptr.DerefDefault(params.UserID, 0),
				Username: ptr.DerefDefault(params.Username, ""),
			},
		})
	case model.RoleUser:
		account, err = s.storage.GetAccount(ctx, model.AccountUser{
			AccountBase: model.AccountBase{
				ID:       ptr.DerefDefault(params.UserID, 0),
				Username: ptr.DerefDefault(params.Username, ""),
			},
			Email: ptr.DerefDefault(params.Email, ""),
			Phone: ptr.DerefDefault(params.Phone, ""),
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
	Role     model.AccountType
	Password string
	Username *string
	Email    *string
	Phone    *string
}

func (s *ServiceImpl) Login(ctx context.Context, params LoginParams) (string, error) {
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

	return s.GenerateAccessToken(account)
}

func (s *ServiceImpl) Register(ctx context.Context, account model.Account) (string, error) {
	txStorage, err := s.storage.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer txStorage.Rollback(ctx)

	newAccount, err := txStorage.CreateAccount(ctx, account)
	if err != nil {
		return "", err
	}

	if newAccount.GetBase().Role == model.RoleUser {
		// Create cart
		err = txStorage.CreateCart(ctx, newAccount.GetBase().ID)
		if err != nil {
			return "", err
		}
	}

	if err = txStorage.Commit(ctx); err != nil {
		return "", err
	}

	return s.GenerateAccessToken(newAccount)
}

type GetPermissionsParams struct {
	AccountID int64
	Role      *model.AccountType
}

func (s *ServiceImpl) GetPermissions(ctx context.Context, params GetPermissionsParams) ([]model.Permission, error) {
	var role model.AccountType

	if params.Role == nil {
		account, err := s.storage.GetAccount(ctx, model.AccountBase{
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
		return s.storage.GetPermissions(ctx, storage.GetPermissionsParams{
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
	Role        *model.AccountType
	Permissions []model.Permission
}

func (s *ServiceImpl) HasPermission(ctx context.Context, params HasPermissionParams) (bool, error) {
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
	AvatarURL            *string
}

func (s *ServiceImpl) UpdateAccount(ctx context.Context, params UpdateAccountParams) (model.AccountBase, error) {
	return s.storage.UpdateAccount(ctx, storage.UpdateAccountParams{
		ID:                   params.ID,
		Username:             params.Username,
		Password:             params.Password,
		NullCustomPermission: params.NullCustomPermission,
		CustomPermission:     params.CustomPermission,
		AvatarURL:            params.AvatarURL,
	})
}

type UpdateAccountUserParams struct {
	ID               int64
	Email            *string
	Phone            *string
	Gender           *model.Gender
	FullName         *string
	DefaultAddressID *int64
}

func (s *ServiceImpl) UpdateAccountUser(ctx context.Context, params UpdateAccountUserParams) (model.AccountUser, error) {
	return s.storage.UpdateAccountUser(ctx, storage.UpdateAccountUserParams{
		ID:               params.ID,
		Email:            params.Email,
		Phone:            params.Phone,
		Gender:           params.Gender,
		FullName:         params.FullName,
		DefaultAddressID: params.DefaultAddressID,
	})
}

func ValidateAccessToken(tokenStr string) (claims model.Claims, err error) {
	secret := config.GetConfig().SensitiveKeys.JWTSecret

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return claims, err
	}

	if !token.Valid {
		return claims, errors.New("invalid token or token expired")
	}

	return claims, nil
}
