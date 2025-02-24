package service

import "shopnexus-go-service/internal/repository"

type Service struct {
	Account *AccountService
	Cart    *CartService
	Payment *PaymentService
	Product *ProductService
}

func NewServices(repo *repository.Repository) *Service {
	account := NewAccountService(repo)
	cart := NewCartService(repo)
	payment := NewPaymentService(repo)
	product := NewProductService(repo)

	return &Service{
		Account: account,
		Cart:    cart,
		Payment: payment,
		Product: product,
	}
}
