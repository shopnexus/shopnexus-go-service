package payment

import (
	"context"
	"fmt"
	"shopnexus-go-service/internal/model"
	repository "shopnexus-go-service/internal/repository"
	"shopnexus-go-service/internal/service/account"
	"shopnexus-go-service/internal/service/product"
	"shopnexus-go-service/internal/util"
)

var _ PaymentServiceInterface = (*PaymentService)(nil)

type PaymentService struct {
	Repo       repository.Repository
	accountSvc *account.AccountService
	productSvc *product.ProductService
	platforms  map[Platform]PaymentPlatform
}

func NewPaymentService(
	repo repository.Repository,
	accountSvc *account.AccountService,
	productSvc *product.ProductService,
) *PaymentService {
	s := &PaymentService{
		Repo:       repo,
		accountSvc: accountSvc,
		productSvc: productSvc,
		platforms:  map[Platform]PaymentPlatform{},
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

func (s *PaymentService) WithTx(txRepo *repository.TxRepository) *PaymentService {
	return NewPaymentService(txRepo, s.accountSvc, s.productSvc)
}

type GetPaymentParams = struct {
	AccountID int64
	Role      model.Role
	PaymentID int64
}

func (s *PaymentService) GetPayment(ctx context.Context, params GetPaymentParams) (model.Payment, error) {
	repoParams := repository.GetPaymentParams{
		ID: params.PaymentID,
	}

	if params.Role == model.RoleUser {
		repoParams.UserID = &params.AccountID
	}

	return s.Repo.GetPayment(ctx, repoParams)
}

type ListPaymentsParams struct {
	model.PaginationParams
	AccountID       int64
	Role            model.Role
	Method          *model.PaymentMethod
	Status          *model.Status
	Address         *string
	TotalFrom       *int64
	TotalTo         *int64
	DateCreatedFrom *int64
	DateCreatedTo   *int64
}

func (s *PaymentService) ListPayments(ctx context.Context, params ListPaymentsParams) (result model.PaginateResult[model.Payment], err error) {
	repoParams := repository.ListPaymentsParams{
		PaginationParams: params.PaginationParams,
		Method:           params.Method,
		Status:           params.Status,
		Address:          params.Address,
		TotalFrom:        params.TotalFrom,
		TotalTo:          params.TotalTo,
		DateCreatedFrom:  params.DateCreatedFrom,
		DateCreatedTo:    params.DateCreatedTo,
	}

	// User only see their own payments
	if params.Role == model.RoleUser {
		repoParams.UserID = &params.AccountID
	}

	total, err := s.Repo.CountPayments(ctx, repoParams)
	if err != nil {
		return result, err
	}

	payments, err := s.Repo.ListPayments(ctx, repoParams)
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

	// Clear the cart
	if err = txRepo.ClearCart(ctx, params.UserID); err != nil {
		return CreatePaymentResult{}, err
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

		var (
			serialIDs  []string
			totalPrice int64
		)

		// Get available sales for the product model
		sales, err := txRepo.GetAvailableSales(ctx, repository.GetLatestSaleParams{
			ProductModelID: productModelItem.GetID(),
			BrandID:        productModel.BrandID,
			Tags:           productModel.Tags,
		})
		if err != nil {
			return CreatePaymentResult{}, err
		}

		for _, pickProduct := range pickProducts {
			serialIDs = append(serialIDs, pickProduct.SerialID)
			totalPrice += productModel.ListPrice + pickProduct.AddPrice
			totalPriceBase := totalPrice

			// Apply sales
			for _, sale := range sales {
				totalPrice -= sale.Apply(totalPriceBase)
			}
		}
		totalPayment += totalPrice

		productOnPayments = append(productOnPayments, model.ProductOnPayment{
			ItemQuantityBase: model.ItemQuantityBase[int64]{
				ItemID:   productModelItem.GetID(),
				Quantity: productModelItem.GetQuantity(),
			},
			SerialIDs:  serialIDs,
			Price:      productModel.ListPrice,
			TotalPrice: totalPrice,
		})
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

	// TODO: move this update product sold to cron job check success payment (because currently we don't know if payment is success or not)
	if err = s.productSvc.WithTx(txRepo).UpdateProductSold(ctx, product.UpdateProductSoldParams{
		IDs: func() []int64 {
			ids := make([]int64, 0, len(productOnPayments))
			for _, pop := range productOnPayments {
				ids = append(ids, pop.ItemID)
			}
			return ids
		}(),
		Amount: 1,
	}); err != nil {
		return CreatePaymentResult{}, err
	}

	// Rollback if purchase failed
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
	ID        int64
	AccountID int64
	Role      model.Role
	Method    *model.PaymentMethod
	Address   *string
	Status    *model.Status
}

func (s *PaymentService) UpdatePayment(ctx context.Context, params UpdatePaymentParams) error {
	txRepo, err := s.Repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer txRepo.Rollback(ctx)

	getPaymentParams := repository.GetPaymentParams{
		ID:     params.ID,
		Status: util.ToPtr(model.StatusPending),
	}

	// User only see their own payments
	if params.Role == model.RoleUser {
		getPaymentParams.UserID = &params.AccountID
	}

	// Payment must be pending
	payment, err := txRepo.GetPayment(ctx, getPaymentParams)
	if err != nil {
		return err
	}

	// If payment method is cash, address is required
	if (params.Method == nil && payment.Method == model.PaymentMethodCash || params.Method != nil && *params.Method == model.PaymentMethodCash) &&
		(params.Address == nil && payment.Address == "" || params.Address != nil && *params.Address == "") {
		return fmt.Errorf("address is required for payment method %s", *params.Method)
	}

	// If params.Status is not nil, check if account has permission to update status
	if params.Status != nil {
		if ok, err := s.accountSvc.HasPermission(ctx, account.HasPermissionParams{
			AccountID: params.AccountID,
			Permissions: []model.Permission{
				model.PermissionUpdatePayment,
			},
		}); err != nil {
			return err
		} else if !ok {
			return fmt.Errorf("account %d does not have permission to update payment status", params.AccountID)
		}
	}

	if err = txRepo.UpdatePayment(ctx, repository.UpdatePaymentParams{
		ID:      params.ID,
		Method:  params.Method,
		Address: params.Address,
		Status:  params.Status,
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
		UserID: &params.UserID,
		Status: util.ToPtr(model.StatusCancelled),
	}); err != nil {
		return err
	}

	if err = txRepo.Commit(ctx); err != nil {
		return err
	}

	return nil
}
