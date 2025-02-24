package service

import (
	"context"
	"shopnexus-go-service/internal/model"
	repository "shopnexus-go-service/internal/repository"
	"shopnexus-go-service/internal/util"

	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	repo *repository.Repository
}

func NewAccountService(repo *repository.Repository) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

func (s *AccountService) IsAdmin(ctx context.Context, accountID int64) (bool, error) {
	accountBase, err := s.repo.GetAccountBase(ctx, accountID)
	if err != nil {
		return false, err
	}

	return accountBase.Role == model.RoleAdmin, nil
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
	Password string
}

func (s *AccountService) FindAccount(ctx context.Context, params FindAccountParams) (model.Account, error) {
	if params.Username == nil && params.Email == nil && params.Phone == nil {
		return nil, model.ErrInvalidCreds
	}

	account, err := s.repo.GetAccount(ctx, model.AccountUser{
		AccountBase: model.AccountBase{
			Username: util.DerefDefault(params.Username, ""),
		},
		Email: util.DerefDefault(params.Email, ""),
		Phone: util.DerefDefault(params.Phone, ""),
	})
	if err != nil {
		return nil, err
	}

	// Check hash password
	if ok := s.CheckPassword(account.GetBase().Password, params.Password); !ok {
		return nil, model.ErrWrongPassword
	}

	return account, nil
}

type LoginUserParams struct {
	Username *string
	Email    *string
	Phone    *string
	Role     model.Role
	Password string
}

func (s *AccountService) Login(ctx context.Context, params LoginUserParams) (string, error) {
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
		})
	case model.RoleUser:
		account, err = s.FindAccount(ctx, FindAccountParams{
			Username: params.Username,
			Email:    params.Email,
			Phone:    params.Phone,
			Password: params.Password,
		})
	}

	if err != nil {
		return "", err
	}

	return util.GenerateAccessToken(account)
}

func (s *AccountService) Register(ctx context.Context, account model.Account) (string, error) {
	newAccount, err := s.repo.CreateAccount(ctx, account)
	if err != nil {
		return "", err
	}

	return util.GenerateAccessToken(newAccount)
}
