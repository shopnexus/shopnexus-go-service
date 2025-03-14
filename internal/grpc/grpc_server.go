package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	pgxutil "shopnexus-go-service/internal/db/pgx"
	"shopnexus-go-service/internal/grpc/account"
	"shopnexus-go-service/internal/grpc/interceptor"
	"shopnexus-go-service/internal/grpc/payment"
	"shopnexus-go-service/internal/grpc/product"
	"shopnexus-go-service/internal/repository"
	"shopnexus-go-service/internal/service"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	accountv1 "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/account/v1"
	accountv1connect "github.com/shopnexus/shopnexus-protobuf-gen-go/pb/account/v1/accountv1connect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/payment/v1/paymentv1connect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/product/v1/productv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func NewServer(address string) error {
	db := pgxutil.GetPgxPool()
	repo := repository.NewRepository(db)
	services := service.NewServices(repo)

	mux := http.NewServeMux()

	// Setup account service
	accountPath, accountHandler := accountv1connect.NewAccountServiceHandler(
		account.NewAccountServer(services.Account),
		connect.WithInterceptors(interceptor.NewAuthInterceptor(
			accountv1connect.AccountServiceGetCartProcedure,
			accountv1connect.AccountServiceAddCartItemProcedure,
			accountv1connect.AccountServiceUpdateCartItemProcedure,
		)),
	)
	mux.Handle(accountPath, accountHandler)

	// Setup payment service
	paymentPath, paymentHandler := paymentv1connect.NewPaymentServiceHandler(
		payment.NewPaymentServer(services.Payment),
		connect.WithInterceptors(interceptor.NewAuthInterceptor()),
	)
	mux.Handle(paymentPath, paymentHandler)

	// Setup product service
	productPath, productHandler := productv1connect.NewProductServiceHandler(
		product.NewProductServer(services.Product),
		connect.WithInterceptors(interceptor.NewAuthInterceptor()),
	)
	mux.Handle(productPath, productHandler)

	// Setup reflection for postman service discovery
	reflector := grpcreflect.NewStaticReflector(
		accountv1connect.AccountServiceName,
		paymentv1connect.PaymentServiceName,
		productv1connect.ProductServiceName,
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	go func() {
		time.Sleep(time.Second * 3)
		NewClient()
	}()

	return http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

func NewClient() {
	client := accountv1connect.NewAccountServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		connect.WithGRPC(),
	)

	resp, err := client.LoginAdmin(context.Background(), &connect.Request[accountv1.LoginAdminRequest]{
		Msg: &accountv1.LoginAdminRequest{
			Username: "admin",
			Password: "admin",
		},
	})
	js, _ := json.Marshal(resp)
	fmt.Println(string(js), err)

}
