package connect

import (
	"net/http"
	"shopnexus-go-service/internal/server/connect/handler/account"
	"shopnexus-go-service/internal/server/connect/handler/file"
	"shopnexus-go-service/internal/server/connect/handler/payment"
	"shopnexus-go-service/internal/server/connect/handler/product"
	"shopnexus-go-service/internal/server/connect/interceptor/permission"
	"shopnexus-go-service/internal/service"

	"connectrpc.com/connect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/file/v1/filev1connect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/payment/v1/paymentv1connect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/product/v1/productv1connect"
)

func Init(mux *http.ServeMux, services *service.Services, withReflector bool) error {
	if withReflector {
		InitReflector(mux)
	}

	connectOpts := []connect.HandlerOption{
		connect.WithInterceptors(
			permission.NewPermissionInterceptor(services.Account, PermissionRoutes),
		),
	}

	// File
	filePath, fileHandler := filev1connect.NewFileServiceHandler(
		file.NewFileServiceHandler(services.S3),
		connectOpts...,
	)
	mux.Handle(filePath, fileHandler)

	// Account
	accountPath, accountHandler := account.NewAccountServiceHandler(
		services.Account,
		connectOpts...,
	)
	mux.Handle(accountPath, accountHandler)

	// Product
	productPath, productHandler := productv1connect.NewProductServiceHandler(
		product.NewProductServiceHandler(services.Product),
		connectOpts...,
	)
	mux.Handle(productPath, productHandler)

	// Payment
	paymentPath, paymentHandler := paymentv1connect.NewPaymentServiceHandler(
		payment.NewPaymentServiceHandler(services.Payment),
		connectOpts...,
	)
	mux.Handle(paymentPath, paymentHandler)

	return nil
}
