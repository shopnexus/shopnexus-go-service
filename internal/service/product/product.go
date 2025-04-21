package product

import (
	"context"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/repository"
	"shopnexus-go-service/internal/service/account"
)

type ProductService struct {
	repo       repository.Repository
	accountSvc account.AccountServiceInterface
}

var _ ProductServiceInterface = (*ProductService)(nil)

func NewProductService(repo repository.Repository, accountSvc account.AccountServiceInterface) *ProductService {
	return &ProductService{
		repo:       repo,
		accountSvc: accountSvc,
	}
}

func (s *ProductService) WithTx(txRepo repository.Repository) *ProductService {
	//TODO: Use WithTX to all injected service
	return NewProductService(txRepo, s.accountSvc.WithTx(txRepo))
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (model.Product, error) {
	return s.repo.GetProduct(ctx, id)
}

type ListProductsParams = repository.ListProductsParams

func (s *ProductService) ListProducts(ctx context.Context, params ListProductsParams) (result model.PaginateResult[model.Product], err error) {
	total, err := s.repo.CountProducts(ctx, params)
	if err != nil {
		return result, err
	}

	products, err := s.repo.ListProducts(ctx, params)
	if err != nil {
		return result, err
	}

	return model.PaginateResult[model.Product]{
		Data:       products,
		Limit:      params.Limit,
		Page:       params.Page,
		Total:      total,
		NextPage:   params.NextPage(total),
		NextCursor: nil,
	}, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	newProduct, err := s.repo.CreateProduct(ctx, product)
	if err != nil {
		return model.Product{}, err
	}

	return newProduct, nil
}

type UpdateProductParams = struct {
	// TODO: sửa lại ko xài RepoParams, phải tự ghi ra hết
	ID             int64
	ProductModelID *int64
	Quantity       *int64
	Sold           *int64
	AddPrice       *int64
	CanCombine     *bool
	IsActive       *bool
	Metadata       *[]byte
	Resources      *[]string
}

func (s *ProductService) UpdateProduct(ctx context.Context, params UpdateProductParams) error {
	txRepo, err := s.repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer txRepo.Rollback(ctx)

	if err = s.repo.UpdateProduct(ctx, repository.UpdateProductParams{
		ID:             params.ID,
		ProductModelID: params.ProductModelID,
		Quantity:       params.Quantity,
		Sold:           params.Sold,
		AddPrice:       params.AddPrice,
		CanCombine:     params.CanCombine,
		IsActive:       params.IsActive,
		Metadata:       params.Metadata,
		Resources:      params.Resources,
	}); err != nil {
		return err
	}

	if err = txRepo.Commit(ctx); err != nil {
		return err
	}

	return nil
}

type UpdateProductSoldParams = struct {
	IDs    []int64
	Amount int64
}

func (s *ProductService) UpdateProductSold(ctx context.Context, params UpdateProductSoldParams) error {
	txRepo, err := s.repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer txRepo.Rollback(ctx)

	if err = s.repo.UpdateProductSold(ctx, params.IDs, params.Amount); err != nil {
		return err
	}

	if err = txRepo.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.DeleteProduct(ctx, id)
}

type ProductServiceInterface interface {
	GetBrand(ctx context.Context, id int64) (model.Brand, error)
	ListBrands(ctx context.Context, params ListBrandsParams) (model.PaginateResult[model.Brand], error)
	CreateBrand(ctx context.Context, params CreateBrandParams) (model.Brand, error)
	UpdateBrand(ctx context.Context, params UpdateBrandParams) error
	DeleteBrand(ctx context.Context, id int64) error

	GetProductModel(ctx context.Context, id int64) (model.ProductModel, error)
	ListProductModels(ctx context.Context, params ListProductModelsParams) (model.PaginateResult[model.ProductModel], error)
	CreateProductModel(ctx context.Context, params CreateProductModelParams) (model.ProductModel, error)
	UpdateProductModel(ctx context.Context, params UpdateProductModelParams) error
	DeleteProductModel(ctx context.Context, id int64) error
	ListProductTypes(ctx context.Context, params ListProductTypesParams) ([]model.ProductType, error)

	GetProduct(ctx context.Context, id int64) (model.Product, error)
	ListProducts(ctx context.Context, params ListProductsParams) (model.PaginateResult[model.Product], error)
	CreateProduct(ctx context.Context, product model.Product) (model.Product, error)
	UpdateProduct(ctx context.Context, params UpdateProductParams) error
	DeleteProduct(ctx context.Context, id int64) error

	GetProductSerial(ctx context.Context, serialID string) (model.ProductSerial, error)
	ListProductSerials(ctx context.Context, params ListProductSerialsParams) (model.PaginateResult[model.ProductSerial], error)
	CreateProductSerial(ctx context.Context, serial model.ProductSerial) (model.ProductSerial, error)
	UpdateProductSerial(ctx context.Context, params UpdateProductSerialParams) error
	DeleteProductSerial(ctx context.Context, params DeleteProductSerialPParams) error
	MarkProductSerialsAsSold(ctx context.Context, serialIDs []string) error

	GetSale(ctx context.Context, id int64) (model.Sale, error)
	ListSales(ctx context.Context, params ListSalesParams) (model.PaginateResult[model.Sale], error)
	CreateSale(ctx context.Context, params CreateSaleParams) (model.Sale, error)
	UpdateSale(ctx context.Context, params UpdateSaleParams) error
	DeleteSale(ctx context.Context, id int64) error

	GetTag(ctx context.Context, tag string) (TagResponse, error)
	ListTags(ctx context.Context, params ListTagsParams) (model.PaginateResult[TagResponse], error)
	CreateTag(ctx context.Context, tag model.Tag) error
	UpdateTag(ctx context.Context, params UpdateTagParams) error
	DeleteTag(ctx context.Context, tag string) error
}
