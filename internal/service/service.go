package service

import (
	"shopnexus-go-service/internal/repository"
	"shopnexus-go-service/internal/service/account"
	"shopnexus-go-service/internal/service/payment"
	"shopnexus-go-service/internal/service/product"
	"shopnexus-go-service/internal/service/s3"
)

type Services struct {
	Account *account.AccountService
	Payment *payment.PaymentService
	Product *product.ProductService
	S3      *s3.S3Service
}

func NewServices(repo *repository.RepositoryImpl) *Services {
	accountSvc := account.NewAccountService(repo)
	productSvc := product.NewProductService(repo, accountSvc)
	paymentSvc := payment.NewPaymentService(repo, accountSvc, productSvc)
	s3Svc := s3.NewS3Service(repo)

	return &Services{
		Account: accountSvc,
		Payment: paymentSvc,
		Product: productSvc,
		S3:      s3Svc,
	}
}
