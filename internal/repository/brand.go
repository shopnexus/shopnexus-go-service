package repository

import (
	"context"
	"shopnexus-go-service/gen/sqlc"
	pgxutil "shopnexus-go-service/internal/db/pgx"
	"shopnexus-go-service/internal/model"

	"github.com/jackc/pgx/v5/pgtype"
)

func (r *RepositoryImpl) GetBrand(ctx context.Context, id int64) (model.Brand, error) {
	brand, err := r.sqlc.GetBrand(ctx, id)
	if err != nil {
		return model.Brand{}, err
	}

	resources, err := r.GetResources(ctx, brand.ID, model.ResourceTypeBrand)
	if err != nil {
		return model.Brand{}, err
	}

	return model.Brand{
		ID:          brand.ID,
		Name:        brand.Name,
		Description: brand.Description,
		Resources:   resources,
	}, nil
}

type ListBrandsParams struct {
	model.PaginationParams
	Name        *string
	Description *string
}

func (r *RepositoryImpl) CountBrands(ctx context.Context, params ListBrandsParams) (int64, error) {
	return r.sqlc.CountBrands(ctx, sqlc.CountBrandsParams{
		Name:        *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
		Description: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Description),
	})
}

func (r *RepositoryImpl) ListBrands(ctx context.Context, params ListBrandsParams) ([]model.Brand, error) {
	brands, err := r.sqlc.ListBrands(ctx, sqlc.ListBrandsParams{
		Offset:      params.Offset(),
		Limit:       params.Limit,
		Name:        *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Name),
		Description: *pgxutil.PtrToPgtype(&pgtype.Text{}, params.Description),
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.Brand, len(brands))
	for i, brand := range brands {
		resources, err := r.GetResources(ctx, brand.ID, model.ResourceTypeBrand)
		if err != nil {
			return nil, err
		}

		result[i] = model.Brand{
			ID:          brand.ID,
			Name:        brand.Name,
			Description: brand.Description,
			Resources:   resources,
		}
	}
	return result, nil
}

func (r *RepositoryImpl) CreateBrand(ctx context.Context, brand model.Brand) (model.Brand, error) {
	row, err := r.sqlc.CreateBrand(ctx, sqlc.CreateBrandParams{
		Name:        brand.Name,
		Description: brand.Description,
	})
	if err != nil {
		return model.Brand{}, err
	}

	if err = r.AddResources(ctx, row.ID, model.ResourceTypeBrand, brand.Resources); err != nil {
		return model.Brand{}, err
	}

	return model.Brand{
		ID:          row.ID,
		Name:        brand.Name,
		Description: brand.Description,
		Resources:   brand.Resources,
	}, nil
}

type UpdateBrandParams struct {
	ID          int64
	Name        *string
	Description *string
}

func (r *RepositoryImpl) UpdateBrand(ctx context.Context, params UpdateBrandParams) error {
	return r.sqlc.UpdateBrand(ctx, sqlc.UpdateBrandParams{
		ID:          params.ID,
		Name:        *pgxutil.PtrToPgtype(&pgtype.Text{}, &params.Name),
		Description: *pgxutil.PtrToPgtype(&pgtype.Text{}, &params.Description),
	})
}

func (r *RepositoryImpl) DeleteBrand(ctx context.Context, id int64) error {
	return r.sqlc.DeleteBrand(ctx, id)
}
