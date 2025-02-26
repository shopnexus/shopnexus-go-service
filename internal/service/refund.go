package service

import (
	"context"
	"shopnexus-go-service/internal/model"
	repository "shopnexus-go-service/internal/repository"
)

type RefundService struct {
	Repo *repository.Repository
}

func NewRefundService(repo *repository.Repository) *RefundService {
	return &RefundService{
		Repo: repo,
	}
}

type GetRefundParams struct {
	ID     int64
	UserID int64
}

func (s *RefundService) Get(ctx context.Context, params GetRefundParams) (model.Refund, error) {
	refund, err := s.Repo.GetRefund(ctx, repository.GetRefundParams{
		ID:     params.ID,
		UserID: &params.UserID,
	})
	if err != nil {
		return model.Refund{}, err
	}

	return refund, nil
}

type CreateRefundParams struct {
	UserID    int64
	PaymentID int64
	Method    model.RefundMethod
	Reason    string
	Address   *string
	Resources []string
}

func (s *RefundService) Create(ctx context.Context, params CreateRefundParams) (model.Refund, error) {
	txRepo, err := s.Repo.Begin(ctx)
	if err != nil {
		return model.Refund{}, err
	}
	defer txRepo.Rollback(ctx)

	// Method drop_off must not contains address
	if params.Method == model.RefundMethodDropOff && params.Address != nil {
		return model.Refund{}, model.ErrMalformedParams
	}

	// Method pick_up must contains address
	if params.Method == model.RefundMethodPickUp && params.Address == nil {
		return model.Refund{}, model.ErrMalformedParams
	}

	refund, err := txRepo.CreateRefund(ctx, model.Refund{
		PaymentID: params.PaymentID,
		Method:    params.Method,
		Status:    model.StatusPending,
		Reason:    params.Reason,
		Address:   params.Address,
	})
	if err != nil {
		return model.Refund{}, err
	}

	if err = txRepo.Commit(ctx); err != nil {
		return model.Refund{}, err
	}

	return refund, nil
}
