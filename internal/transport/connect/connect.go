package connect

import (
	"net/http"
	"shopnexus-go-service/internal/logger"
	"shopnexus-go-service/internal/service"
	"shopnexus-go-service/internal/transport/connect/handler/account"
	"shopnexus-go-service/internal/transport/connect/handler/payment"
	"shopnexus-go-service/internal/transport/connect/handler/product"
	"shopnexus-go-service/internal/transport/connect/interceptor/permission"

	"connectrpc.com/connect"
)

func Init(mux *http.ServeMux, services *service.Services, withReflector bool) error {
	if withReflector {
		InitReflector(mux)
		logger.Log.Info("Connect reflector enabled")
	}

	connectOpts := []connect.HandlerOption{
		connect.WithInterceptors(
			permission.NewPermissionInterceptor(services.Account, PermissionRoutes),
		),
	}

	// Account
	accountPath, accountHandler := account.NewAccountServiceHandler(
		services.Account,
		connectOpts...,
	)
	mux.Handle(accountPath, accountHandler)

	// Product
	productPath, productHandler := product.NewProductServiceHandler(
		services.Product,
		connectOpts...,
	)
	mux.Handle(productPath, productHandler)

	// Payment
	paymentPath, paymentHandler := payment.NewPaymentServiceHandler(
		services.Payment,
		connectOpts...,
	)
	mux.Handle(paymentPath, paymentHandler)

	logger.Log.Info("Connect server initialized")

	return nil
}
