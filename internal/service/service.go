package service

import "shopnexus-go-service/internal/repository"

type Service struct {
	Account *AccountService
	Cart    *CartService
	Payment *PaymentService
	Product *ProductService
	Refund  *RefundService
	S3      *S3Service
}

func NewServices(repo *repository.Repository) *Service {
	account := NewAccountService(repo)
	cart := NewCartService(repo)
	payment := NewPaymentService(repo)
	product := NewProductService(repo)
	refund := NewRefundService(repo)
	s3 := NewS3Service(repo)

	return &Service{
		Account: account,
		Cart:    cart,
		Payment: payment,
		Product: product,
		Refund:  refund,
		S3:      s3,
	}
}
