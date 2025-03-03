package payment

import (
	"context"
	"fmt"
	"shopnexus-go-service/internal/model"
	repository "shopnexus-go-service/internal/repository"
)

type PaymentService struct {
	Repo      *repository.Repository
	platforms map[Platform]PaymentPlatform
}

func NewPaymentService(repo *repository.Repository) *PaymentService {
	s := &PaymentService{
		Repo:      repo,
		platforms: map[Platform]PaymentPlatform{},
	}

	// Init payment platforms
	vnpay := &VnpayPlatform{}
	s.platforms[PlatformVNPAY] = vnpay

	return s
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

type PaymentServiceInterface interface {
	CreatePayment(ctx context.Context, params CreatePaymentParams) (CreatePaymentResult, error)
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

	// Create payment url
	var pp PaymentPlatform

	switch params.PaymentMethod {
	case model.PaymentMethodVnpay:
		pp, err = s.GetPlatform(PlatformVNPAY)
		if err != nil {
			return CreatePaymentResult{}, err
		}
	case model.PaymentMethodMomo:
		// TODO: support momo payment
		return CreatePaymentResult{}, fmt.Errorf("payment method momo not yet supported")
		// pp, err = s.GetPlatform(PlatformMOMO)
		// if err != nil {
		// 	return CreatePaymentResult{}, err
		// }
	case model.PaymentMethodCash:
		// Do nothing
		// TODO: add logic for cash payment
		return CreatePaymentResult{}, fmt.Errorf("payment method cash not yet supported")
	default:
		return CreatePaymentResult{}, fmt.Errorf("payment method %s not supported", params.PaymentMethod)
	}

	url, err := pp.CreateOrder(ctx, CreateOrderParams{
		PaymentID: newPayment.ID,
		Info:      fmt.Sprintf("Payment for order %d", newPayment.ID),
		Amount:    newPayment.Total,
	})
	if err != nil {
		return CreatePaymentResult{}, err
	}

	if err = txRepo.Commit(ctx); err != nil {
		return CreatePaymentResult{}, err
	}

	return CreatePaymentResult{Payment: newPayment, Url: url}, nil
}

func (s *PaymentService) GetPlatform(platform Platform) (PaymentPlatform, error) {
	pp, ok := s.platforms[platform]
	if !ok {
		return nil, fmt.Errorf("platform %s not found", platform)
	}
	return pp, nil
}
