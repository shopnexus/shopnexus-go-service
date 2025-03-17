package payment

import (
	"context"
	"fmt"
	"shopnexus-go-service/internal/model"
	repository "shopnexus-go-service/internal/repository"
	"shopnexus-go-service/internal/util"
)

var _ PaymentServiceInterface = (*PaymentService)(nil)

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

type PaymentServiceInterface interface {
	// Payment
	GetPayment(ctx context.Context, params GetPaymentParams) (model.Payment, error)
	ListPayments(ctx context.Context, params ListPaymentsParams) (model.PaginateResult[model.Payment], error)
	CreatePayment(ctx context.Context, params CreatePaymentParams) (CreatePaymentResult, error)
	UpdatePayment(ctx context.Context, params UpdatePaymentParams) error
	CancelPayment(ctx context.Context, params CancelPaymentParams) error

	// Refund
	GetRefund(ctx context.Context, params GetRefundParams) (model.Refund, error)
	ListRefunds(ctx context.Context, params ListRefundsParams) (model.PaginateResult[model.Refund], error)
	CreateRefund(ctx context.Context, params CreateRefundParams) (model.Refund, error)
	UpdateRefund(ctx context.Context, params UpdateRefundParams) error
	CancelRefund(ctx context.Context, params CancelRefundParams) error
}

type GetPaymentParams = struct {
	UserID    int64
	PaymentID int64
}

func (s *PaymentService) GetPayment(ctx context.Context, params GetPaymentParams) (model.Payment, error) {
	return s.Repo.GetPayment(ctx, repository.GetPaymentParams{
		ID:     params.PaymentID,
		UserID: &params.UserID,
	})
}

type ListPaymentsParams = repository.ListPaymentsParams

func (s *PaymentService) ListPayments(ctx context.Context, params ListPaymentsParams) (result model.PaginateResult[model.Payment], err error) {
	total, err := s.Repo.CountPayments(ctx, params)
	if err != nil {
		return result, err
	}

	payments, err := s.Repo.ListPayments(ctx, params)
	if err != nil {
		return result, err
	}

	return model.PaginateResult[model.Payment]{
		Data:       payments,
		Limit:      params.Limit,
		Page:       params.Page,
		Total:      total,
		NextPage:   params.NextPage(total),
		NextCursor: nil,
	}, nil
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

	// Get user cart
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

	// Calculate total payment
	// Iterate through each product model in the cart
	for _, productModelItem := range cart.ProductModels {
		// Get product model details
		productModel, err := txRepo.GetProductModel(ctx, productModelItem.GetID())
		if err != nil {
			return CreatePaymentResult{}, err
		}

		// Get any available products from that product model
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

	// Create payment
	newPayment, err := txRepo.CreatePayment(ctx, model.Payment{
		UserID:   params.UserID,
		Address:  params.Address,
		Method:   params.PaymentMethod,
		Total:    totalPayment,
		Status:   model.StatusPending,
		Products: productOnPayments,
	})
	if err != nil {
		return CreatePaymentResult{}, err
	}

	// Create payment url
	var pp PaymentPlatform

	switch params.PaymentMethod {
	case model.PaymentMethodVnpay:
		pp, err = s.getPlatform(PlatformVNPAY)
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

func (s *PaymentService) getPlatform(platform Platform) (PaymentPlatform, error) {
	pp, ok := s.platforms[platform]
	if !ok {
		return nil, fmt.Errorf("platform %s not found", platform)
	}
	return pp, nil
}

type UpdatePaymentParams struct {
	ID      int64
	UserID  int64
	Method  *model.PaymentMethod
	Address *string
}

func (s *PaymentService) UpdatePayment(ctx context.Context, params UpdatePaymentParams) error {
	txRepo, err := s.Repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer txRepo.Rollback(ctx)

	exists, err := txRepo.ExistsPayment(ctx, repository.GetPaymentParams{
		ID:     params.ID,
		UserID: &params.UserID,
	})
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("payment %d not found", params.ID)
	}

	payment, err := txRepo.GetPayment(ctx, repository.GetPaymentParams{
		ID:     params.ID,
		UserID: &params.UserID,
	})
	if err != nil {
		return err
	}

	// If payment method is cash, address is required
	if (params.Method == nil && payment.Method == model.PaymentMethodCash || params.Method != nil && *params.Method == model.PaymentMethodCash) &&
		(params.Address == nil && payment.Address == "" || params.Address != nil && *params.Address == "") {
		return fmt.Errorf("address is required for payment method %s", *params.Method)
	}

	if err = txRepo.UpdatePayment(ctx, repository.UpdatePaymentParams{
		ID:      params.ID,
		Method:  params.Method,
		Address: params.Address,
	}); err != nil {
		return err
	}

	if err = txRepo.Commit(ctx); err != nil {
		return err
	}

	return nil
}

type CancelPaymentParams = struct {
	UserID    int64
	PaymentID int64
}

func (s *PaymentService) CancelPayment(ctx context.Context, params CancelPaymentParams) error {
	txRepo, err := s.Repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer txRepo.Rollback(ctx)

	payment, err := txRepo.GetPayment(ctx, repository.GetPaymentParams{
		ID:     params.PaymentID,
		UserID: &params.UserID,
	})
	if err != nil {
		return err
	}

	// No need to check ownership as we already check it in GetPayment
	// if payment.UserID != *params.UserID {
	// 	return fmt.Errorf("payment %d not belong to user %d", params.PaymentID, params.UserID)
	// }

	if payment.Status != model.StatusPending {
		return fmt.Errorf("payment %d cannot be cancelled", params.PaymentID)
	}

	if err = txRepo.UpdatePayment(ctx, repository.UpdatePaymentParams{
		ID:     params.PaymentID,
		Status: util.ToPtr(model.StatusCancelled),
	}); err != nil {
		return err
	}

	if err = txRepo.Commit(ctx); err != nil {
		return err
	}

	return nil
}

type CancelRefundParams = struct {
	UserID   int64
	RefundID int64
}

func (s *PaymentService) CancelRefund(ctx context.Context, params CancelRefundParams) error {
	txRepo, err := s.Repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer txRepo.Rollback(ctx)

	refund, err := txRepo.GetRefund(ctx, repository.GetRefundParams{
		ID:     params.RefundID,
		UserID: &params.UserID,
	})
	if err != nil {
		return err
	}

	if refund.Status != model.StatusPending {
		return fmt.Errorf("refund %d cannot be cancelled", params.RefundID)
	}

	if err = txRepo.UpdateRefund(ctx, repository.UpdateRefundParams{
		ID:     params.RefundID,
		Status: util.ToPtr(model.StatusCancelled),
	}); err != nil {
		return err
	}

	if err = txRepo.Commit(ctx); err != nil {
		return err
	}

	return nil
}
