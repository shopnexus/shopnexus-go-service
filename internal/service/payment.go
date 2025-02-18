package service

import (
	"context"
	"shopnexus-go-service/gen/pb"
)

type PaymentService struct {
	pb.UnimplementedPaymentServer
}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (s *PaymentService) Create(ctx context.Context, params *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	return nil, nil
}
