package connect

import (
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/server/connect/interceptor/permission"
)

var PermissionRoutes = map[string]permission.Options{
	// OTHER will always return permission denied

	// ACCOUNT
	"/account.v1.AccountService/GetUser":       permission.UseOptions(permission.NeedRoles(model.RoleUser)),
	"/account.v1.AccountService/GetAdmin":      permission.UseOptions(permission.NeedRoles(model.RoleAdmin)),
	"/account.v1.AccountService/GetUserPublic": permission.UseOptions(permission.NeedRoles()),
	"/account.v1.AccountService/LoginUser":     permission.UseOptions(permission.NeedRoles()),
	"/account.v1.AccountService/RegisterUser":  permission.UseOptions(permission.NeedRoles()),
	"/account.v1.AccountService/LoginAdmin":    permission.UseOptions(permission.NeedRoles()),
	"/account.v1.AccountService/RegisterAdmin": permission.UseOptions(permission.NeedRoles()),
	"/account.v1.AccountService/UpdateAccount": permission.UseOptions(permission.NeedRoles()),
	"/account.v1.AccountService/UpdateUser":    permission.UseOptions(permission.NeedRoles()),

	// CART
	"/account.v1.AccountService/GetCart":        permission.UseOptions(permission.NeedRoles(model.RoleUser)),
	"/account.v1.AccountService/AddCartItem":    permission.UseOptions(permission.NeedRoles(model.RoleUser)),
	"/account.v1.AccountService/UpdateCartItem": permission.UseOptions(permission.NeedRoles(model.RoleUser)),
	"/account.v1.AccountService/ClearCart":      permission.UseOptions(permission.NeedRoles(model.RoleUser)),

	// ADDRESS
	"/account.v1.AccountService/GetAddress":    permission.UseOptions(permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin)),
	"/account.v1.AccountService/ListAddresses": permission.UseOptions(permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin)),
	"/account.v1.AccountService/CreateAddress": permission.UseOptions(permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin)),
	"/account.v1.AccountService/UpdateAddress": permission.UseOptions(permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin)),
	"/account.v1.AccountService/DeleteAddress": permission.UseOptions(permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin)),

	// PRODUCT MODEL
	"/product.v1.ProductService/ListProductModels": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
	"/product.v1.ProductService/ListProductTypes": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
	"/product.v1.ProductService/GetProductModel": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
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
	"/product.v1.ProductService/ListProducts": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
	"/product.v1.ProductService/GetProduct": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
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

	// PRODUCT SERIAL
	"/product.v1.ProductService/ListProductSerials": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
	"/product.v1.ProductService/GetProductSerial": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
	"/product.v1.ProductService/CreateProductSerial": permission.UseOptions(
		permission.NeedRoles(model.RoleAdmin, model.RoleStaff),
		permission.NeedPermissions(model.PermissionCreateProductSerial),
	),
	"/product.v1.ProductService/UpdateProductSerial": permission.UseOptions(
		permission.NeedRoles(model.RoleAdmin, model.RoleStaff),
		permission.NeedPermissions(model.PermissionUpdateProductSerial),
	),
	"/product.v1.ProductService/DeleteProductSerial": permission.UseOptions(
		permission.NeedRoles(model.RoleAdmin, model.RoleStaff),
		permission.NeedPermissions(model.PermissionDeleteProductSerial),
	),

	// SALE
	"/product.v1.ProductService/ListSales": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
	"/product.v1.ProductService/GetSale": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
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
	"/product.v1.ProductService/ListTags": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
	"/product.v1.ProductService/GetTag": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
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

	// BRAND
	"/product.v1.ProductService/ListBrands": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
	"/product.v1.ProductService/GetBrand": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
	"/product.v1.ProductService/CreateBrand": permission.UseOptions(
		permission.NeedRoles(model.RoleAdmin, model.RoleStaff),
		permission.NeedPermissions(model.PermissionCreateBrand),
	),
	"/product.v1.ProductService/UpdateBrand": permission.UseOptions(
		permission.NeedRoles(model.RoleAdmin, model.RoleStaff),
		permission.NeedPermissions(model.PermissionUpdateBrand),
	),
	"/product.v1.ProductService/DeleteBrand": permission.UseOptions(
		permission.NeedRoles(model.RoleAdmin, model.RoleStaff),
		permission.NeedPermissions(model.PermissionDeleteBrand),
	),

	// PAYMENT
	"/payment.v1.PaymentService/ListPayments": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
	"/payment.v1.PaymentService/GetPayment": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
	"/payment.v1.PaymentService/CreatePayment": permission.UseOptions(
		permission.NeedRoles(model.RoleUser),
	),
	"/payment.v1.PaymentService/UpdatePayment": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
		permission.NeedPermissions(model.PermissionUpdatePayment),
	),
	"/payment.v1.PaymentService/DeletePayment": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
		permission.NeedPermissions(model.PermissionDeletePayment),
	),

	// REFUND
	"/payment.v1.PaymentService/ListRefunds": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
	"/payment.v1.PaymentService/GetRefund": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
	"/payment.v1.PaymentService/CreateRefund": permission.UseOptions(
		permission.NeedRoles(model.RoleUser),
	),
	"/payment.v1.PaymentService/UpdateRefund": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
		permission.NeedPermissions(model.PermissionUpdateRefund),
	),
	"/payment.v1.PaymentService/DeleteRefund": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
		permission.NeedPermissions(model.PermissionDeleteRefund),
	),

	// COMMENT
	"/product.v1.ProductService/ListComments": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
	"/product.v1.ProductService/GetComment": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
	),
	"/product.v1.ProductService/CreateComment": permission.UseOptions(
		permission.NeedRoles(model.RoleUser),
	),
	"/product.v1.ProductService/UpdateComment": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
		permission.NeedPermissions(model.PermissionUpdateComment),
	),
	"/product.v1.ProductService/DeleteComment": permission.UseOptions(
		permission.NeedRoles(model.RoleUser, model.RoleStaff, model.RoleAdmin),
		permission.NeedPermissions(model.PermissionDeleteComment),
	),
}
