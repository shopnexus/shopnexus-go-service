package service

import (
	"context"
	"fmt"
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
}

type CreatePaymentResult struct {
	Payment model.Payment
	Url     string
}

func (s *PaymentService) CreatePayment(ctx context.Context, params CreatePaymentParams) (CreatePaymentResult, error) {
	txRepo, err := s.Repo.Begin(ctx)
	if err != nil {
		return CreatePaymentResult{}, err
	}
	defer txRepo.Rollback(ctx)

	cart, err := txRepo.GetCart(ctx, params.UserID)
	if err != nil {
		return CreatePaymentResult{}, err
	}

	if len(cart.ProductModels) == 0 {
		return CreatePaymentResult{}, fmt.Errorf("cart is empty")
	}

	var (
		productOnPayments []model.ProductOnPayment
		totalPayment      int64
	)

	for _, productModelItem := range cart.ProductModels {
		// Get product model details
		productModel, err := txRepo.GetProductModel(ctx, productModelItem.GetID())
		if err != nil {
			return CreatePaymentResult{}, err
		}

		// Get available products from that product model
		pickProducts, err := txRepo.GetAvailableProducts(
			ctx,
			productModelItem.GetID(),
			productModelItem.GetQuantity(),
		)
		if err != nil {
			return CreatePaymentResult{}, err
		}

		for _, pickProduct := range pickProducts {
			// TODO: add discount logic to Price or TotalPrice
			totalPrice := productModel.ListPrice * productModelItem.GetQuantity()
			totalPayment += totalPrice

			productOnPayments = append(productOnPayments, model.ProductOnPayment{
				ItemQuantityBase: model.ItemQuantityBase[string]{
					ItemID:   pickProduct.SerialID,
					Quantity: productModelItem.GetQuantity(),
				},
				Price:      productModel.ListPrice,
				TotalPrice: totalPrice,
			})
		}
	}

	newPayment, err := txRepo.CreatePayment(ctx, model.Payment{
		UserID:        params.UserID,
		Address:       params.Address,
		PaymentMethod: params.PaymentMethod,
		Total:         totalPayment,
		Status:        model.StatusPending,
		Products:      productOnPayments,
	})
	if err != nil {
		return CreatePaymentResult{}, err
	}

	if err = txRepo.Commit(ctx); err != nil {
		return CreatePaymentResult{}, err
	}

	// Payment URL generation
	url := "https://payment.com/" + fmt.Sprint(newPayment.ID)

	return CreatePaymentResult{Payment: newPayment, Url: url}, nil
}
