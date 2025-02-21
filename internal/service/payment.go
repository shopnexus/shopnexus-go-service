package service

import (
	"context"
	"shopnexus-go-service/gen/pb"
	"shopnexus-go-service/internal/model"
)

type PaymentService struct {
	pb.UnimplementedPaymentServer
}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
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
