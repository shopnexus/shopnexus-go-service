package product

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/repository"
)

type ListSalesParams struct {
	model.PaginationParams
	Tag             *string
	ProductModelID  *int64
	BrandID         *int64
	DateStartedFrom *int64
	DateStartedTo   *int64
	DateEndedFrom   *int64
	DateEndedTo     *int64
	IsActive        *bool
}

func (s *ProductService) GetSale(ctx context.Context, id int64) (model.Sale, error) {
	return s.repo.GetSale(ctx, id)
}

func (s *ProductService) ListSales(ctx context.Context, params ListSalesParams) (model.PaginateResult[model.Sale], error) {
	count, err := s.repo.CountSales(ctx, repository.ListSalesParams{
		PaginationParams: params.PaginationParams,
		Tag:              params.Tag,
		ProductModelID:   params.ProductModelID,
		BrandID:          params.BrandID,
		DateStartedFrom:  params.DateStartedFrom,
		DateStartedTo:    params.DateStartedTo,
		DateEndedFrom:    params.DateEndedFrom,
		DateEndedTo:      params.DateEndedTo,
		IsActive:         params.IsActive,
	})
	if err != nil {
		return model.PaginateResult[model.Sale]{}, err
	}

	data, err := s.repo.ListSales(ctx, repository.ListSalesParams{
		PaginationParams: params.PaginationParams,
		Tag:              params.Tag,
		ProductModelID:   params.ProductModelID,
		BrandID:          params.BrandID,
		DateStartedFrom:  params.DateStartedFrom,
		DateStartedTo:    params.DateStartedTo,
		DateEndedFrom:    params.DateEndedFrom,
		DateEndedTo:      params.DateEndedTo,
		IsActive:         params.IsActive,
	})
	if err != nil {
		return model.PaginateResult[model.Sale]{}, err
	}

	return model.PaginateResult[model.Sale]{
		Data:  data,
		Total: count,
		Page:  params.Page,
		Limit: params.Limit,
	}, nil
}

type CreateSaleParams struct {
	Sale model.Sale
}

func (s *ProductService) CreateSale(ctx context.Context, params CreateSaleParams) (model.Sale, error) {
	return s.repo.CreateSale(ctx, params.Sale)
}

type UpdateSaleParams struct {
	ID              int64
	Tag             *string
	ProductModelID  *int64
	BrandID         *int64
	DateStarted     *int64
	DateEnded       *int64
	Quantity        *int64
	Used            *int64
	IsActive        *bool
	DiscountPercent *int32
	DiscountPrice   *int64
}

func (s *ProductService) UpdateSale(ctx context.Context, params UpdateSaleParams) error {
	return s.repo.UpdateSale(ctx, repository.UpdateSaleParams{
		ID:              params.ID,
		Tag:             params.Tag,
		ProductModelID:  params.ProductModelID,
		BrandID:         params.BrandID,
		DateStarted:     params.DateStarted,
		DateEnded:       params.DateEnded,
		Quantity:        params.Quantity,
		Used:            params.Used,
		IsActive:        params.IsActive,
		DiscountPercent: params.DiscountPercent,
		DiscountPrice:   params.DiscountPrice,
	})
}

func (s *ProductService) DeleteSale(ctx context.Context, id int64) error {
	return s.repo.DeleteSale(ctx, id)
}
