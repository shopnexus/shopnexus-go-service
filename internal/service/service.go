package service

import (
	"shopnexus-go-service/internal/service/account"
	"shopnexus-go-service/internal/service/payment"
	"shopnexus-go-service/internal/service/product"
	"shopnexus-go-service/internal/service/s3"
	"shopnexus-go-service/internal/service/storage"
	"shopnexus-go-service/internal/service/vnpay"
)

type Services struct {
	Account account.Service
	Payment payment.Service
	Product product.Service
	S3      s3.Service
	Vnpay   vnpay.Service
}

func NewServices() (*Services, error) {
	storageSvc, err := storage.NewService()
	if err != nil {
		return nil, err
	}

	accountSvc, err := account.NewService(storageSvc)
	if err != nil {
		return nil, err
	}
	productSvc, err := product.NewService(storageSvc, accountSvc)
	if err != nil {
		return nil, err
	}
	paymentSvc, err := payment.NewService(storageSvc, accountSvc, productSvc)
	if err != nil {
		return nil, err
	}
	s3Svc, err := s3.NewService()
	if err != nil {
		return nil, err
	}
	vnpaySvc, err := vnpay.NewService()
	if err != nil {
		return nil, err
	}

	return &Services{
		Account: accountSvc,
		Payment: paymentSvc,
		Product: productSvc,
		S3:      s3Svc,
		Vnpay:   vnpaySvc,
	}, nil
}
