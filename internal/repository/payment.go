package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	"shopnexus-go-service/internal/model"
)

func (r *Repository) CreatePayment(ctx context.Context, payment model.Payment) (model.Payment, error) {
	row, err := r.sqlc.CreatePayment(ctx, sqlc.CreatePaymentParams{
		UserID:        payment.UserID,
		Address:       payment.Address,
		PaymentMethod: sqlc.PaymentPaymentMethod(payment.PaymentMethod),
		Total:         payment.Total,
		Status:        sqlc.PaymentStatus(payment.Status),
	})
	if err != nil {
		return model.Payment{}, err
	}

	var createArgs []sqlc.CreatePaymentProductsParams
	for _, product := range payment.Products {
		createArgs = append(createArgs, sqlc.CreatePaymentProductsParams{
			PaymentID:       payment.ID,
			ProductSerialID: product.ItemID,
			Quantity:        product.Quantity,
			Price:           product.Price,
			TotalPrice:      product.TotalPrice,
		})
	}

	_, err = r.sqlc.CreatePaymentProducts(ctx, createArgs)
	if err != nil {
		return model.Payment{}, err
	}

	return model.Payment{
		ID:            row.ID,
		UserID:        row.UserID,
		Address:       row.Address,
		PaymentMethod: model.PaymentMethod(row.PaymentMethod),
		Total:         row.Total,
		Status:        model.PaymentStatus(row.Status),
		DateCreated:   row.DateCreated.Time.UnixMilli(),
		Products:      payment.Products,
	}, nil
}
