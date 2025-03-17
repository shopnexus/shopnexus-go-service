package grpc

import (
	"shopnexus-go-service/internal/grpc/interceptor/permission"
	"shopnexus-go-service/internal/model"
)

var permissionRoutes = map[string]permission.Options{
	// OTHER will skip validate roles & permissions

	// ACCOUNT
	"/account.v1.AccountService/GetCart":        permission.UseOptions(permission.NeedRoles(model.RoleUser)),
	"/account.v1.AccountService/AddCartItem":    permission.UseOptions(permission.NeedRoles(model.RoleUser)),
	"/account.v1.AccountService/UpdateCartItem": permission.UseOptions(permission.NeedRoles(model.RoleUser)),
	"/account.v1.AccountService/ClearCart":      permission.UseOptions(permission.NeedRoles(model.RoleUser)),

	// PRODUCT MODEL
	"/product.v1.ProductService/CreateProductModel": permission.UseOptions(
		permission.NeedRoles(model.RoleAdmin, model.RoleStaff),
		permission.NeedPermissions(model.PermissionCreateProductModel),
	),
	"/product.v1.ProductService/UpdateProductModel": permission.UseOptions(
		permission.NeedRoles(model.RoleAdmin, model.RoleStaff),
		permission.NeedPermissions(model.PermissionUpdateProductModel),
	),
	"/product.v1.ProductService/DeleteProductModel": permission.UseOptions(
		permission.NeedRoles(model.RoleAdmin, model.RoleStaff),
		permission.NeedPermissions(model.PermissionDeleteProductModel),
	),

	// PRODUCT
	"/product.v1.ProductService/CreateProduct": permission.UseOptions(
		permission.NeedRoles(model.RoleAdmin, model.RoleStaff),
		permission.NeedPermissions(model.PermissionCreateProduct),
	),
	"/product.v1.ProductService/UpdateProduct": permission.UseOptions(
		permission.NeedRoles(model.RoleAdmin, model.RoleStaff),
		permission.NeedPermissions(model.PermissionUpdateProduct),
	),
	"/product.v1.ProductService/DeleteProduct": permission.UseOptions(
		permission.NeedRoles(model.RoleAdmin, model.RoleStaff),
		permission.NeedPermissions(model.PermissionDeleteProduct),
	),

	// SALE
	"/product.v1.ProductService/CreateSale": permission.UseOptions(
		permission.NeedRoles(model.RoleAdmin, model.RoleStaff),
		permission.NeedPermissions(model.PermissionCreateSale),
	),
	"/product.v1.ProductService/UpdateSale": permission.UseOptions(
		permission.NeedRoles(model.RoleAdmin, model.RoleStaff),
		permission.NeedPermissions(model.PermissionUpdateSale),
	),
	"/product.v1.ProductService/DeleteSale": permission.UseOptions(
		permission.NeedRoles(model.RoleAdmin, model.RoleStaff),
		permission.NeedPermissions(model.PermissionDeleteSale),
	),

	// TAG
	"/product.v1.ProductService/CreateTag": permission.UseOptions(
		permission.NeedRoles(model.RoleAdmin, model.RoleStaff),
		permission.NeedPermissions(model.PermissionCreateTag),
	),
	"/product.v1.ProductService/UpdateTag": permission.UseOptions(
		permission.NeedRoles(model.RoleAdmin, model.RoleStaff),
		permission.NeedPermissions(model.PermissionUpdateTag),
	),
	"/product.v1.ProductService/DeleteTag": permission.UseOptions(
		permission.NeedRoles(model.RoleAdmin, model.RoleStaff),
		permission.NeedPermissions(model.PermissionDeleteTag),
	),
}
