package service

import (
	"context"
	"shopnexus-go-service/internal/model"
	repository "shopnexus-go-service/internal/repository"
)

type PaymentService struct {
	Repo *repository.Repository
}

func NewPaymentService(repo *repository.Repository) *PaymentService {
	return &PaymentService{
		Repo: repo,
	}
}

type CreatePaymentParams struct {
	UserID        int64
	Address       string
	PaymentMethod model.PaymentMethod
	Total         float64
}

func (s *PaymentService) CreatePayment(ctx context.Context, params CreatePaymentParams) (model.Payment, error) {
	return model.Payment{}, nil
}
