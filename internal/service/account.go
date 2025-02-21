package service

import (
	"context"
	"shopnexus-go-service/internal/model"
	repository "shopnexus-go-service/internal/repository"
)

type AccountService struct {
	Repo *repository.Repository
}

func NewAccountService(repo *repository.Repository) *AccountService {
	return &AccountService{
		Repo: repo,
	}
}

func (s *AccountService) IsAdmin(ctx context.Context, accountID int64) (bool, error) {
	accountBase, err := s.Repo.GetAccountBase(ctx, accountID)
	if err != nil {
		return false, err
	}

	return accountBase.Role == model.RoleAdmin, nil
}
