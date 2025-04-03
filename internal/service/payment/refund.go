package payment

import (
	"context"
	"fmt"
	"shopnexus-go-service/internal/model"
	repository "shopnexus-go-service/internal/repository"
	"shopnexus-go-service/internal/service/account"
	"shopnexus-go-service/internal/util"
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

	// Payment must exist and is successful
	existsPayment, err := txRepo.ExistsPayment(ctx, repository.GetPaymentParams{
		ID:     params.PaymentID,
		UserID: &params.UserID,
		Status: util.ToPtr(model.StatusSuccess),
	})
	if err != nil {
		return model.Refund{}, err
	}
	if !existsPayment {
		return model.Refund{}, fmt.Errorf("payment %d not found", params.PaymentID)
	}

	// Check if there is an existing refund for the payment
	existsRefund, err := txRepo.ExistsRefund(ctx, repository.ExistsRefundParams{
		PaymentID: params.PaymentID,
		UserID:    params.UserID,
	})
	if err != nil {
		return model.Refund{}, err
	}
	if existsRefund {
		return model.Refund{}, fmt.Errorf("refund for payment %d already exists and is pending or resolved", params.PaymentID)
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

	if params.Status != nil {
		// Check if account has permission to update refund status
		if ok, err := s.accountSvc.HasPermission(ctx, account.HasPermissionParams{
			AccountID: params.UserID,
			Permissions: []model.Permission{
				model.PermissionUpdateRefund,
			},
		}); !ok {
			return fmt.Errorf("account %d has no permission to update refund status: %w", params.UserID, err)
		}
	}

	if err = txRepo.UpdateRefund(ctx, repository.UpdateRefundParams{
		ID:      params.ID,
		UserID:  &params.UserID,
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
