package product

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/storage/pgx/repository"
)

func (s *ProductService) GetBrand(ctx context.Context, id int64) (model.Brand, error) {
	brand, err := s.repo.GetBrand(ctx, id)
	if err != nil {
		return model.Brand{}, err
	}

	return brand, nil
}

type ListBrandsParams = repository.ListBrandsParams

func (s *ProductService) ListBrands(ctx context.Context, params ListBrandsParams) (result model.PaginateResult[model.Brand], err error) {
	total, err := s.repo.CountBrands(ctx, params)
	if err != nil {
		return result, err
	}

	brands, err := s.repo.ListBrands(ctx, params)
	if err != nil {
		return result, err
	}

	return model.PaginateResult[model.Brand]{
		Data:       brands,
		Limit:      params.Limit,
		Page:       params.Page,
		Total:      total,
		NextPage:   params.NextPage(total),
		NextCursor: nil,
	}, nil
}

type CreateBrandParams struct {
	UserID int64
	model.Brand
}

func (s *ProductService) CreateBrand(ctx context.Context, params CreateBrandParams) (model.Brand, error) {
	txRepo, err := s.repo.Begin(ctx)
	if err != nil {
		return model.Brand{}, err
	}
	defer txRepo.Rollback(ctx)

	newBrand, err := txRepo.CreateBrand(ctx, params.Brand)
	if err != nil {
		return model.Brand{}, err
	}

	if err = txRepo.Commit(ctx); err != nil {
		return model.Brand{}, err
	}

	return newBrand, nil
}

type UpdateBrandParams struct {
	RepoParams repository.UpdateBrandParams
	Resources  []string
}

func (s *ProductService) UpdateBrand(ctx context.Context, params UpdateBrandParams) error {
	txRepo, err := s.repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer txRepo.Rollback(ctx)

	if err = txRepo.UpdateBrand(ctx, params.RepoParams); err != nil {
		return err
	}

	if err = txRepo.UpdateResources(ctx, params.RepoParams.ID, model.ResourceTypeBrand, params.Resources); err != nil {
		return err
	}

	return txRepo.Commit(ctx)
}

func (s *ProductService) DeleteBrand(ctx context.Context, id int64) error {
	return s.repo.DeleteBrand(ctx, id)
}
