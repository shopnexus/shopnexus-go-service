package repository

import (
	"context"
	"shopnexus-go-service/internal/model"
)

func (r *Repository) CreatePayment(ctx context.Context, payment model.Payment) (model.Payment, error) {
	return model.Payment{}, nil
}
