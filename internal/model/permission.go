package model

type Permission int

const (
	PermissionUnknown Permission = iota

	PermissionCreateProduct
	PermissionUpdateProduct
	PermissionDeleteProduct

	PermissionCreateProductModel
	PermissionUpdateProductModel
	PermissionDeleteProductModel

	PermissionCreateSale
	PermissionUpdateSale
	PermissionDeleteSale

	PermissionCreateTag
	PermissionUpdateTag
	PermissionDeleteTag

	PermissionCreateBrand
	PermissionUpdateBrand
	PermissionDeleteBrand

	PermissionUpdatePayment
	PermissionDeletePayment
	PermissionUpdateRefund
	PermissionDeleteRefund
)

func GetAllPermissions() []Permission {
	return []Permission{
		PermissionCreateProduct,
		PermissionUpdateProduct,
		PermissionDeleteProduct,
		PermissionCreateProductModel,
		PermissionUpdateProductModel,
		PermissionDeleteProductModel,
		PermissionCreateSale,
		PermissionUpdateSale,
		PermissionDeleteSale,
		PermissionCreateTag,
		PermissionUpdateTag,
		PermissionDeleteTag,
		PermissionCreateBrand,
		PermissionUpdateBrand,
		PermissionDeleteBrand,
		PermissionUpdatePayment,
		PermissionDeletePayment,
		PermissionUpdateRefund,
		PermissionDeleteRefund,
	}
}
