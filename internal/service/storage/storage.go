package storage

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"shopnexus-go-service/config"
	"shopnexus-go-service/gen/sqlc"
	"shopnexus-go-service/internal/client/pgxpool"
	"shopnexus-go-service/internal/model"

	"github.com/jackc/pgx/v5"
)

type Service interface {
	// Transaction methods
	Begin(ctx context.Context) (*TxStorage, error)
	// Commit(ctx context.Context) error
	// Rollback(ctx context.Context) error

	// Account methods
	GetAccountBase(ctx context.Context, accountID int64) (model.AccountBase, error)
	GetAccountUser(ctx context.Context, params GetAccountUserParams) (model.AccountUser, error)
	GetAccountAdmin(ctx context.Context, params GetAccountAdminParams) (model.AccountAdmin, error)
	GetAccount(ctx context.Context, find model.Account) (model.Account, error)
	CreateAccount(ctx context.Context, account model.Account) (model.Account, error)
	UpdateAccount(ctx context.Context, params UpdateAccountParams) (model.AccountBase, error)
	UpdateAccountUser(ctx context.Context, params UpdateAccountUserParams) (model.AccountUser, error)
	GetPermissions(ctx context.Context, params GetPermissionsParams) ([]model.Permission, error)

	// Address methods
	GetAddress(ctx context.Context, params GetAddressParams) (model.Address, error)
	CountAddresses(ctx context.Context, params ListAddressesParams) (int64, error)
	ListAddresses(ctx context.Context, params ListAddressesParams) ([]model.Address, error)
	CreateAddress(ctx context.Context, address model.Address) (model.Address, error)
	UpdateAddress(ctx context.Context, params UpdateAddressParams) (model.Address, error)
	DeleteAddress(ctx context.Context, params DeleteAddressParams) error

	// Product methods
	GetProduct(ctx context.Context, id int64) (model.Product, error)
	GetAvailableProducts(ctx context.Context, productID, amount int64) ([]model.ProductSerial, error)
	CountProducts(ctx context.Context, params ListProductsParams) (int64, error)
	ListProducts(ctx context.Context, params ListProductsParams) ([]model.Product, error)
	CreateProduct(ctx context.Context, product model.Product) (model.Product, error)
	UpdateProduct(ctx context.Context, params UpdateProductParams) error
	UpdateProductSold(ctx context.Context, ids []int64, amount int64) error
	DeleteProduct(ctx context.Context, id int64) error
	GetResources(ctx context.Context, ownerID int64, resourceType model.ResourceType) ([]string, error)
	AddResources(ctx context.Context, ownerID int64, resourceType model.ResourceType, resources []string) error
	UpdateResources(ctx context.Context, ownerID int64, resourceType model.ResourceType, resources []string) error
	EmptyResources(ctx context.Context, ownerID int64, resourceType model.ResourceType) error
	GetProductByPOPID(ctx context.Context, productOnPaymentID int64) (model.Product, error)

	// Product Model methods
	GetProductModel(ctx context.Context, id int64) (model.ProductModel, error)
	GetProductSerialIDs(ctx context.Context, productID int64) ([]string, error)
	CountProductModels(ctx context.Context, params ListProductModelsParams) (int64, error)
	ListProductModels(ctx context.Context, params ListProductModelsParams) ([]model.ProductModel, error)
	CreateProductModel(ctx context.Context, productModel model.ProductModel) (model.ProductModel, error)
	UpdateProductModel(ctx context.Context, params UpdateProductModelParams) error
	DeleteProductModel(ctx context.Context, id int64) error
	CountProductTypes(ctx context.Context, params ListProductTypesParams) (int64, error)
	ListProductTypes(ctx context.Context, params ListProductTypesParams) ([]model.ProductType, error)

	// Sale methods
	GetSale(ctx context.Context, id int64) (model.Sale, error)
	GetAvailableSales(ctx context.Context, params GetLatestSaleParams) ([]model.Sale, error)
	CountSales(ctx context.Context, params ListSalesParams) (int64, error)
	ListSales(ctx context.Context, params ListSalesParams) ([]model.Sale, error)
	CreateSale(ctx context.Context, sale model.Sale) (model.Sale, error)
	UpdateSale(ctx context.Context, params UpdateSaleParams) error
	DeleteSale(ctx context.Context, id int64) error

	// Cart methods
	ExistsCart(ctx context.Context, userID int64) (bool, error)
	CreateCart(ctx context.Context, userID int64) error
	GetCart(ctx context.Context, params GetCartParams) (model.Cart, error)
	AddCartItem(ctx context.Context, params AddCartItemParams) (int64, error)
	UpdateCartItem(ctx context.Context, params UpdateCartItemParams) (int64, error)
	RemoveCartItem(ctx context.Context, cartID int64, productIDs []int64) error
	ClearCart(ctx context.Context, cartID int64) error

	// Payment methods
	ExistsPayment(ctx context.Context, params GetPaymentParams) (bool, error)
	GetPayment(ctx context.Context, params GetPaymentParams) (model.Payment, error)
	GetPaymentProducts(ctx context.Context, paymentID int64) ([]model.ProductOnPayment, error)
	CountPayments(ctx context.Context, params ListPaymentsParams) (int64, error)
	ListPayments(ctx context.Context, params ListPaymentsParams) ([]model.Payment, error)
	CreatePayment(ctx context.Context, payment model.Payment) (model.Payment, error)
	UpdatePayment(ctx context.Context, params UpdatePaymentParams) error
	DeletePayment(ctx context.Context, params DeletePaymentParams) error

	// Refund methods
	ExistsRefund(ctx context.Context, params ExistsRefundParams) (bool, error)
	GetRefund(ctx context.Context, params GetRefundParams) (model.Refund, error)
	CountRefunds(ctx context.Context, params ListRefundsParams) (int64, error)
	ListRefunds(ctx context.Context, params ListRefundsParams) ([]model.Refund, error)
	CreateRefund(ctx context.Context, refund model.Refund) (model.Refund, error)
	UpdateRefund(ctx context.Context, params UpdateRefundParams) error
	DeleteRefund(ctx context.Context, params DeleteRefundParams) error

	// Brand methods
	GetBrand(ctx context.Context, id int64) (model.Brand, error)
	CountBrands(ctx context.Context, params ListBrandsParams) (int64, error)
	ListBrands(ctx context.Context, params ListBrandsParams) ([]model.Brand, error)
	CreateBrand(ctx context.Context, brand model.Brand) (model.Brand, error)
	UpdateBrand(ctx context.Context, params UpdateBrandParams) error
	DeleteBrand(ctx context.Context, id int64) error

	// Tag methods
	GetTag(ctx context.Context, tag string) (model.Tag, error)
	CountTags(ctx context.Context, params ListTagsParams) (int64, error)
	ListTags(ctx context.Context, params ListTagsParams) ([]model.Tag, error)
	CreateTag(ctx context.Context, tag model.Tag) error
	UpdateTag(ctx context.Context, params UpdateTagParams) error
	DeleteTag(ctx context.Context, tag string) error
	CountProductModelsOnTag(ctx context.Context, tag string) (int64, error)
	GetTags(ctx context.Context, productModelID int64) ([]string, error)
	AddTags(ctx context.Context, productModelID int64, tags []string) error
	RemoveTags(ctx context.Context, productModelID int64, tags []string) error
	UpdateTags(ctx context.Context, productModelID int64, tags []string) error

	// Comment methods
	GetComment(ctx context.Context, id int64) (model.Comment, error)
	CountComments(ctx context.Context, params ListCommentsParams) (int64, error)
	ListComments(ctx context.Context, params ListCommentsParams) ([]model.Comment, error)
	CreateComment(ctx context.Context, comment model.Comment) (model.Comment, error)
	UpdateComment(ctx context.Context, params UpdateCommentParams) error
	DeleteComment(ctx context.Context, params DeleteCommentParams) error

	// Product Serial methods
	GetProductSerial(ctx context.Context, serialID string) (model.ProductSerial, error)
	CountProductSerials(ctx context.Context, params ListProductSerialsParams) (int64, error)
	ListProductSerials(ctx context.Context, params ListProductSerialsParams) ([]model.ProductSerial, error)
	CreateProductSerial(ctx context.Context, serial model.ProductSerial) (model.ProductSerial, error)
	UpdateProductSerial(ctx context.Context, params UpdateProductSerialParams) error
	DeleteProductSerial(ctx context.Context, serialID string) error
	MarkProductSerialsAsSold(ctx context.Context, serialIDs []string) error
}

type ServiceImpl struct {
	db   pgxpool.DBTX
	sqlc *sqlc.Queries
}

func NewService() (Service, error) {
	pool, err := pgxpool.NewPgxpool(pgxpool.PgxpoolOptions{
		Url:             config.GetConfig().Postgres.Url,
		Host:            config.GetConfig().Postgres.Host,
		Port:            config.GetConfig().Postgres.Port,
		Username:        config.GetConfig().Postgres.Username,
		Password:        config.GetConfig().Postgres.Password,
		Database:        config.GetConfig().Postgres.Database,
		MaxConnections:  config.GetConfig().Postgres.MaxConnections,
		MaxConnIdleTime: config.GetConfig().Postgres.MaxConnIdleTime,
	})
	if err != nil {
		return nil, err
	}

	return &ServiceImpl{
		db:   pool,
		sqlc: sqlc.New(pool),
	}, nil
}

func (r *ServiceImpl) Begin(ctx context.Context) (*TxStorage, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &TxStorage{
		// Replace the db with the tx'ed db
		ServiceImpl: &ServiceImpl{
			db:   tx,
			sqlc: sqlc.New(tx),
		},
		tx: tx,
	}, nil
}

type TxStorage struct {
	*ServiceImpl
	tx pgx.Tx
}

func (r *TxStorage) Commit(ctx context.Context) error {
	return r.tx.Commit(ctx)
}

func (r *TxStorage) Rollback(ctx context.Context) error {
	return r.tx.Rollback(ctx)
}

// ===================================================
// UTILITIES
// ===================================================

// PgtypeToPtr converts a pgtype to a pointer value
func PgtypeToPtr[T any](v driver.Valuer) *T {
	data, err := v.Value()
	if err != nil {
		return nil
	}

	if data == nil {
		return nil
	}

	if val, ok := data.(T); ok {
		return &val
	}

	return nil
}

// PtrToPgtype convert a pointer value to a pgtype, panics if an error occurs
//
// Branded types should be converted to string before calling this function (e.g. string(myType)).
// Timestamptz should be converted to time.Time before calling this function
func PtrToPgtype[T sql.Scanner, G any](base T, v *G) T {
	if v == nil {
		return base
	}

	// If branded type is not converted to string before calling this function, it will error since pgx cannot Scan our type
	err := base.Scan(*v)
	if err != nil {
		panic(err)
	}
	return base
}

// ValueToPgtype convert a value to a pgtype, panics if an error occurs
//
// Branded types should be converted to string before calling this function (e.g. string(myType))
// Timestamptz should be converted to time.Time before calling this function
func ValueToPgtype[T sql.Scanner, G any](base T, v G) T {
	// If branded type is not converted to string before calling this function, it will error since pgx cannot Scan our type
	err := base.Scan(v)
	if err != nil {
		panic(err)
	}
	return base
}

// PtrBrandedToPgType convert a pointer branded type to a pgtype, panics if an error occurs
//
// Branded type should use this over *PtrToPgtype
func PtrBrandedToPgType[T sql.Scanner, G ~string](base T, v *G) T {
	if v == nil {
		return base
	}

	err := base.Scan(string(*v))
	if err != nil {
		panic(err)
	}
	return base
}

// BrandedToPgType convert a branded type to a pgtype, panics if an error occurs
//
// Branded type should use this over ToPgtype
func BrandedToPgType[T sql.Scanner, G ~string](base T, v G) T {
	err := base.Scan(string(v))
	if err != nil {
		panic(err)
	}
	return base
}
