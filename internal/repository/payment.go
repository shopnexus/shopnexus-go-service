package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	pgxutil "shopnexus-go-service/internal/db/pgx"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/util"

	"github.com/jackc/pgx/v5/pgtype"
)

type GetPaymentParams struct {
	ID     int64
	UserID *int64
	Status *model.Status
}

func (r *Repository) ExistsPayment(ctx context.Context, params GetPaymentParams) (bool, error) {
	return r.sqlc.ExistsPayment(ctx, sqlc.ExistsPaymentParams{
		ID:     params.ID,
		UserID: *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UserID),
		Status: *pgxutil.PtrBrandedToPgType(&sqlc.NullPaymentStatus{}, params.Status),
	})
}

// GetPayment retrieves a payment by its ID, and optionally the user ID.
// If the user ID is provided, the payment must belong to the user.
func (r *Repository) GetPayment(ctx context.Context, params GetPaymentParams) (model.Payment, error) {
	row, err := r.sqlc.GetPayment(ctx, sqlc.GetPaymentParams{
		ID:     params.ID,
		UserID: *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UserID),
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

func (r *Repository) GetPaymentProducts(ctx context.Context, paymentID int64) ([]model.ProductOnPayment, error) {
	rows, err := r.sqlc.GetPaymentProducts(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	productByModel := map[int64]model.ProductOnPayment{}
	for _, row := range rows {
		if product, exists := productByModel[row.ProductModelID]; exists {
			product.Quantity += row.Quantity
			product.SerialIDs = append(product.SerialIDs, row.ProductSerialID)
			product.TotalPrice += row.TotalPrice
			productByModel[row.ProductModelID] = product
		} else {
			productByModel[row.ProductModelID] = model.ProductOnPayment{
				ItemQuantityBase: model.ItemQuantityBase[int64]{
					ItemID:   row.ProductModelID,
					Quantity: row.Quantity,
				},
				SerialIDs:  []string{row.ProductSerialID},
				Price:      row.Price,
				TotalPrice: row.TotalPrice,
			}
		}
	}

	var products []model.ProductOnPayment
	for _, product := range productByModel {
		products = append(products, product)
	}

	return products, nil
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

func (r Repository) CountPayments(ctx context.Context, params ListPaymentsParams) (int64, error) {
	return r.sqlc.CountPayments(ctx, sqlc.CountPaymentsParams{
		UserID:          *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UserID),
		Method:          *pgxutil.PtrToPgtype(&sqlc.NullPaymentPaymentMethod{}, params.Method),
		Status:          *pgxutil.PtrToPgtype(&sqlc.NullPaymentStatus{}, params.Status),
		Address:         *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Address),
		TotalFrom:       *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.TotalFrom),
		TotalTo:         *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.TotalTo),
		DateCreatedFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedTo)),
	})
}

func (r *Repository) ListPayments(ctx context.Context, params ListPaymentsParams) ([]model.Payment, error) {
	rows, err := r.sqlc.ListPayments(ctx, sqlc.ListPaymentsParams{
		Offset:          params.Offset(),
		Limit:           params.Limit,
		UserID:          *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.UserID),
		Method:          *pgxutil.PtrToPgtype(&sqlc.NullPaymentPaymentMethod{}, params.Method),
		Status:          *pgxutil.PtrToPgtype(&sqlc.NullPaymentStatus{}, params.Status),
		Address:         *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Address),
		TotalFrom:       *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.TotalFrom),
		TotalTo:         *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.TotalTo),
		DateCreatedFrom: *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedFrom)),
		DateCreatedTo:   *pgxutil.PtrToPgtype(&pgtype.Timestamptz{}, util.PtrMilisToTime(params.DateCreatedTo)),
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

func (r *Repository) CreatePayment(ctx context.Context, payment model.Payment) (model.Payment, error) {
	row, err := r.sqlc.CreatePayment(ctx, sqlc.CreatePaymentParams{
		UserID:  payment.UserID,
		Method:  sqlc.PaymentPaymentMethod(payment.Method),
		Status:  sqlc.PaymentStatus(payment.Status),
		Address: payment.Address,
		Total:   payment.Total,
	})
	if err != nil {
		return model.Payment{}, err
	}

	var createArgs []sqlc.CreatePaymentProductsParams
	for _, product := range payment.Products {
		for _, serialID := range product.SerialIDs {
			createArgs = append(createArgs, sqlc.CreatePaymentProductsParams{
				PaymentID:       row.ID,
				ProductSerialID: serialID,
				Quantity:        product.Quantity,
				Price:           product.Price,
				TotalPrice:      product.TotalPrice,
			})
		}
	}

	_, err = r.sqlc.CreatePaymentProducts(ctx, createArgs)
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
		Products:    payment.Products,
	}, nil
}

type UpdatePaymentParams struct {
	ID      int64
	Method  *model.PaymentMethod
	Status  *model.Status
	Address *string
	Total   *int64
}

func (r *Repository) UpdatePayment(ctx context.Context, params UpdatePaymentParams) error {
	err := r.sqlc.UpdatePayment(ctx, sqlc.UpdatePaymentParams{
		ID:      params.ID,
		Method:  *pgxutil.PtrToPgtype(&sqlc.NullPaymentPaymentMethod{}, params.Method),
		Status:  *pgxutil.PtrToPgtype(&sqlc.NullPaymentStatus{}, params.Status),
		Address: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Address),
		Total:   *pgxutil.PtrToPgtype(&pgtype.Int8{}, params.Total),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeletePayment(ctx context.Context, paymentID int64) error {
	return r.sqlc.DeletePayment(ctx, paymentID)
}
