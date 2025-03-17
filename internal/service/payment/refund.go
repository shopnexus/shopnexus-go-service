package payment

import (
	"context"
	"fmt"
	"shopnexus-go-service/internal/model"
	repository "shopnexus-go-service/internal/repository"
)

type GetRefundParams struct {
	UserID   int64
	RefundID int64
}

func (s *PaymentService) GetRefund(ctx context.Context, params GetRefundParams) (model.Refund, error) {
	refund, err := s.Repo.GetRefund(ctx, repository.GetRefundParams{
		ID:     params.RefundID,
		UserID: &params.UserID,
	})
	if err != nil {
		return model.Refund{}, err
	}

	return refund, nil
}

type ListRefundsParams = repository.ListRefundsParams

func (s *PaymentService) ListRefunds(ctx context.Context, params ListRefundsParams) (result model.PaginateResult[model.Refund], err error) {
	total, err := s.Repo.CountRefunds(ctx, params)
	if err != nil {
		return result, err
	}

	refunds, err := s.Repo.ListRefunds(ctx, params)
	if err != nil {
		return result, err
	}

	return model.PaginateResult[model.Refund]{
		Data:       refunds,
		Limit:      params.Limit,
		Page:       params.Page,
		Total:      total,
		NextPage:   params.NextPage(total),
		NextCursor: nil,
	}, nil
}

type CreateRefundParams struct {
	UserID    int64
	PaymentID int64
	Method    model.RefundMethod
	Reason    string
	Address   string
	Resources []string
}

func (s *PaymentService) CreateRefund(ctx context.Context, params CreateRefundParams) (model.Refund, error) {
	txRepo, err := s.Repo.Begin(ctx)
	if err != nil {
		return model.Refund{}, err
	}
	defer txRepo.Rollback(ctx)

	// Method drop_off must not contains address
	if params.Method == model.RefundMethodDropOff && params.Address != "" {
		return model.Refund{}, fmt.Errorf("address is not required for refund method drop_off %w", model.ErrMalformedParams)
	}

	// Method pick_up must contains address
	if params.Method == model.RefundMethodPickUp && params.Address == "" {
		return model.Refund{}, fmt.Errorf("address is required for refund method pick_up %w", model.ErrMalformedParams)
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

type UpdateRefundParams struct {
	ID      int64
	UserID  int64
	Method  *model.RefundMethod
	Status  *model.Status
	Reason  *string
	Address *string
}

func (s *PaymentService) UpdateRefund(ctx context.Context, params UpdateRefundParams) error {
	txRepo, err := s.Repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer txRepo.Rollback(ctx)

	// Method drop_off must not contains address
	if *params.Method == model.RefundMethodDropOff {
		params.Address = nil
	}

	if err = txRepo.UpdateRefund(ctx, repository.UpdateRefundParams{
		ID:      params.ID,
		Method:  params.Method,
		Status:  params.Status,
		Reason:  params.Reason,
		Address: params.Address,
	}); err != nil {
		return err
	}

	if err = txRepo.Commit(ctx); err != nil {
		return err
	}

	return nil
}
