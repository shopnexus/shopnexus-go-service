package vnpay

import (
	"context"
	"shopnexus-go-service/config"
	vnpayclient "shopnexus-go-service/internal/client/vnpay"
)

type ServiceImpl struct {
	client vnpayclient.Client
}

type Service interface {
	VerifyPayment(ctx context.Context, ipn map[string]any) error
}

func NewService() (Service, error) {
	return &ServiceImpl{
		client: vnpayclient.NewClient(vnpayclient.ClientOptions{
			TmnCode:    config.GetConfig().Vnpay.TmnCode,
			HashSecret: config.GetConfig().Vnpay.HashSecret,
		}),
	}, nil
}

func (s *ServiceImpl) VerifyPayment(ctx context.Context, ipn map[string]any) error {
	return s.client.VerifyPayment(ctx, ipn)
}
