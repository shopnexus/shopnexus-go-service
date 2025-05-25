package storage

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/utils/ptr"

	"github.com/jackc/pgx/v5/pgtype"
)

type GetPaymentParams struct {
	ID     int64
	UserID *int64
	Status *model.Status
}

func (r *ServiceImpl) ExistsPayment(ctx context.Context, params GetPaymentParams) (bool, error) {
	return r.sqlc.ExistsPayment(ctx, sqlc.ExistsPaymentParams{
		ID:     params.ID,
		UserID: *PtrToPgtype(&pgtype.Int8{}, params.UserID),
		Status: *PtrBrandedToPgType(&sqlc.NullPaymentStatus{}, params.Status),
	})
}

// GetPayment retrieves a payment by its ID, and optionally the user ID.
// If the user ID is provided, the payment must belong to the user.
func (r *ServiceImpl) GetPayment(ctx context.Context, params GetPaymentParams) (model.Payment, error) {
	row, err := r.sqlc.GetPayment(ctx, sqlc.GetPaymentParams{
		ID:     params.ID,
		UserID: *PtrToPgtype(&pgtype.Int8{}, params.UserID),
	})
	if err != nil {
		return model.Payment{}, err
	}

	products, err := r.GetPaymentProducts(ctx, row.ID)
	if err != nil {
		return model.Payment{}, err
	}

	return model.Payment{
		ID:          row.ID,
		UserID:      row.UserID,
		Address:     row.Address,
		Method:      model.PaymentMethod(row.Method),
		Total:       row.Total,
		Status:      model.Status(row.Status),
		DateCreated: row.DateCreated.Time.UnixMilli(),
		Products:    products,
	}, nil
}

func (r *ServiceImpl) GetPaymentProducts(ctx context.Context, paymentID int64) ([]model.ProductOnPayment, error) {
	rows, err := r.sqlc.GetPaymentProducts(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	var products []model.ProductOnPayment
	for _, row := range rows {
		productSerials, err := r.GetPaymentProductSerials(ctx, row.ID)
		if err != nil {
			return nil, err
		}

		serialIDs := make([]string, len(productSerials))
		for i, serial := range productSerials {
			serialIDs[i] = serial.SerialID
		}

		products = append(products, model.ProductOnPayment{
			ID: row.ID,
			ItemQuantityBase: model.ItemQuantityBase[int64]{
				ItemID:   row.ProductID,
				Quantity: row.Quantity,
			},
			SerialIDs:  serialIDs,
			Price:      row.Price,
			TotalPrice: row.TotalPrice,
		})
	}

	return products, nil
}

func (r *ServiceImpl) GetPaymentProductSerials(ctx context.Context, productOnPaymentID int64) ([]model.ProductSerial, error) {
	rows, err := r.sqlc.GetPaymentProductSerials(ctx, productOnPaymentID)
	if err != nil {
		return nil, err
	}

	productSerials := make([]model.ProductSerial, len(rows))
	for i, row := range rows {
		productSerials[i] = model.ProductSerial{
			SerialID:    row.SerialID,
			ProductID:   row.ProductID,
			IsSold:      row.IsSold,
			IsActive:    row.IsActive,
			DateCreated: row.DateCreated.Time.UnixMilli(),
			DateUpdated: row.DateUpdated.Time.UnixMilli(),
		}
	}

	return productSerials, nil
}

type ListPaymentsParams struct {
	model.PaginationParams
	UserID          *int64
	Method          *model.PaymentMethod
	Status          *model.Status
	Address         *string
	TotalFrom       *int64
	TotalTo         *int64
	DateCreatedFrom *int64
	DateCreatedTo   *int64
}

func (r ServiceImpl) CountPayments(ctx context.Context, params ListPaymentsParams) (int64, error) {
	return r.sqlc.CountPayments(ctx, sqlc.CountPaymentsParams{
		UserID:          *PtrToPgtype(&pgtype.Int8{}, params.UserID),
		Method:          *PtrToPgtype(&sqlc.NullPaymentPaymentMethod{}, params.Method),
		Status:          *PtrToPgtype(&sqlc.NullPaymentStatus{}, params.Status),
		Address:         *PtrToPgtype(&pgtype.Text{}, params.Address),
		TotalFrom:       *PtrToPgtype(&pgtype.Int8{}, params.TotalFrom),
		TotalTo:         *PtrToPgtype(&pgtype.Int8{}, params.TotalTo),
		DateCreatedFrom: *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:   *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedTo)),
	})
}

func (r *ServiceImpl) ListPayments(ctx context.Context, params ListPaymentsParams) ([]model.Payment, error) {
	rows, err := r.sqlc.ListPayments(ctx, sqlc.ListPaymentsParams{
		Offset:          params.Offset(),
		Limit:           params.Limit,
		UserID:          *PtrToPgtype(&pgtype.Int8{}, params.UserID),
		Method:          *PtrToPgtype(&sqlc.NullPaymentPaymentMethod{}, params.Method),
		Status:          *PtrToPgtype(&sqlc.NullPaymentStatus{}, params.Status),
		Address:         *PtrToPgtype(&pgtype.Text{}, params.Address),
		TotalFrom:       *PtrToPgtype(&pgtype.Int8{}, params.TotalFrom),
		TotalTo:         *PtrToPgtype(&pgtype.Int8{}, params.TotalTo),
		DateCreatedFrom: *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:   *PtrToPgtype(&pgtype.Timestamptz{}, ptr.PtrMilisToTime(params.DateCreatedTo)),
	})
	if err != nil {
		return nil, err
	}

	var payments []model.Payment
	for _, row := range rows {
		products, err := r.GetPaymentProducts(ctx, row.ID)
		if err != nil {
			return nil, err
		}

		payments = append(payments, model.Payment{
			ID:          row.ID,
			UserID:      row.UserID,
			Address:     row.Address,
			Method:      model.PaymentMethod(row.Method),
			Total:       row.Total,
			Status:      model.Status(row.Status),
			DateCreated: row.DateCreated.Time.UnixMilli(),
			Products:    products,
		})
	}

	return payments, nil
}

func (r *ServiceImpl) CreatePayment(ctx context.Context, payment model.Payment) (model.Payment, error) {
	paymentRow, err := r.sqlc.CreatePayment(ctx, sqlc.CreatePaymentParams{
		UserID:  payment.UserID,
		Method:  sqlc.PaymentPaymentMethod(payment.Method),
		Status:  sqlc.PaymentStatus(payment.Status),
		Address: payment.Address,
		Total:   payment.Total,
	})
	if err != nil {
		return model.Payment{}, err
	}

	// 1. Create the payment products

	var createPaymentProductsArgs []sqlc.CreatePaymentProductsParams
	for _, productOnPayment := range payment.Products {
		createPaymentProductsArgs = append(createPaymentProductsArgs, sqlc.CreatePaymentProductsParams{
			PaymentID:  paymentRow.ID,
			ProductID:  productOnPayment.GetID(),
			Quantity:   productOnPayment.GetQuantity(),
			Price:      productOnPayment.Price,
			TotalPrice: productOnPayment.TotalPrice,
		})
	}

	_, err = r.sqlc.CreatePaymentProducts(ctx, createPaymentProductsArgs)
	if err != nil {
		return model.Payment{}, err
	}

	createdPaymentProducts, err := r.GetPaymentProducts(ctx, paymentRow.ID)
	if err != nil {
		return model.Payment{}, err
	}

	// 2. Create the product serial to these payment products

	var createPaymentProductSerialsArgs []sqlc.CreatePaymentProductSerialsParams
	for popIdx, productOnPayment := range createdPaymentProducts {
		// Assign serial IDs to newly created payment products
		createdPaymentProducts[popIdx].SerialIDs = payment.Products[popIdx].SerialIDs

		// Start creating serial to payment product
		for _, serialID := range payment.Products[popIdx].SerialIDs {
			createPaymentProductSerialsArgs = append(createPaymentProductSerialsArgs, sqlc.CreatePaymentProductSerialsParams{
				ProductOnPaymentID: productOnPayment.ID,
				ProductSerialID:    serialID,
			})
		}
	}

	_, err = r.sqlc.CreatePaymentProductSerials(ctx, createPaymentProductSerialsArgs)
	if err != nil {
		return model.Payment{}, err
	}

	return model.Payment{
		ID:          paymentRow.ID,
		UserID:      paymentRow.UserID,
		Address:     paymentRow.Address,
		Method:      model.PaymentMethod(paymentRow.Method),
		Total:       paymentRow.Total,
		Status:      model.Status(paymentRow.Status),
		DateCreated: paymentRow.DateCreated.Time.UnixMilli(),
		Products:    createdPaymentProducts,
	}, nil
}

type UpdatePaymentParams struct {
	ID      int64
	UserID  *int64
	Method  *model.PaymentMethod
	Status  *model.Status
	Address *string
	Total   *int64
}

func (r *ServiceImpl) UpdatePayment(ctx context.Context, params UpdatePaymentParams) error {
	err := r.sqlc.UpdatePayment(ctx, sqlc.UpdatePaymentParams{
		ID:      params.ID,
		UserID:  *PtrToPgtype(&pgtype.Int8{}, params.UserID),
		Method:  *PtrBrandedToPgType(&sqlc.NullPaymentPaymentMethod{}, params.Method),
		Status:  *PtrBrandedToPgType(&sqlc.NullPaymentStatus{}, params.Status),
		Address: *PtrToPgtype(&pgtype.Text{}, params.Address),
		Total:   *PtrToPgtype(&pgtype.Int8{}, params.Total),
	})
	if err != nil {
		return err
	}

	return nil
}

type DeletePaymentParams struct {
	ID     int64
	UserID *int64
}

func (r *ServiceImpl) DeletePayment(ctx context.Context, params DeletePaymentParams) error {
	return r.sqlc.DeletePayment(ctx, sqlc.DeletePaymentParams{
		ID:     params.ID,
		UserID: *PtrToPgtype(&pgtype.Int8{}, params.UserID),
	})
}
