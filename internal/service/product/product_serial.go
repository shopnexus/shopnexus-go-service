package product

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/repository"
	"shopnexus-go-service/internal/service/account"
)

func (s *ProductService) GetProductSerial(ctx context.Context, serialID string) (model.ProductSerial, error) {
	return s.repo.GetProductSerial(ctx, serialID)
}

type ListProductSerialsParams = repository.ListProductSerialsParams

func (s *ProductService) ListProductSerials(ctx context.Context, params ListProductSerialsParams) (result model.PaginateResult[model.ProductSerial], err error) {
	total, err := s.repo.CountProductSerials(ctx, params)
	if err != nil {
		return result, err
	}

	serials, err := s.repo.ListProductSerials(ctx, params)
	if err != nil {
		return result, err
	}

	return model.PaginateResult[model.ProductSerial]{
		Data:       serials,
		Limit:      params.Limit,
		Page:       params.Page,
		Total:      total,
		NextPage:   params.NextPage(total),
		NextCursor: nil,
	}, nil
}

func (s *ProductService) CreateProductSerial(ctx context.Context, serial model.ProductSerial) (model.ProductSerial, error) {
	return s.repo.CreateProductSerial(ctx, serial)
}

type UpdateProductSerialParams = repository.UpdateProductSerialParams

func (s *ProductService) UpdateProductSerial(ctx context.Context, params UpdateProductSerialParams) error {
	return s.repo.UpdateProductSerial(ctx, params)
}

type DeleteProductSerialPParams struct {
	AccountID int64
	Role      model.Role
	SerialID  string
}

func (s *ProductService) DeleteProductSerial(ctx context.Context, params DeleteProductSerialPParams) error {
	hasPermission, err := s.accountSvc.HasPermission(ctx, account.HasPermissionParams{
		AccountID:   params.AccountID,
		Role:        &params.Role,
		Permissions: []model.Permission{model.PermissionDeleteProductSerial},
	})
	if err != nil {
		return err
	}

	if !hasPermission {
		return model.ErrPermissionDenied
	}

	return s.repo.DeleteProductSerial(ctx, params.SerialID)
}

func (s *ProductService) MarkProductSerialsAsSold(ctx context.Context, serialIDs []string) error {
	return s.repo.MarkProductSerialsAsSold(ctx, serialIDs)
}
