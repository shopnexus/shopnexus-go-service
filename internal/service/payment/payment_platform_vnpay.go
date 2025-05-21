package payment

import (
	"context"
	"shopnexus-go-service/internal/service/vnpay"
)

type VnpayPlatform struct {
	svc vnpay.Service
}

func (p *VnpayPlatform) CreateOrder(ctx context.Context, params CreateOrderParams) (url string, err error) {
	return p.svc.CreateOrder(ctx, vnpay.CreateOrderParams{
		PaymentID: params.PaymentID,
		Amount:    params.Amount,
		Info:      params.Info,
	})
}
