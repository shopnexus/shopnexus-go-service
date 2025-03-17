package grpc

import (
	"fmt"
	"log"
	"net/http"
	pgxutil "shopnexus-go-service/internal/db/pgx"
	"shopnexus-go-service/internal/grpc/handler/account"
	"shopnexus-go-service/internal/grpc/handler/payment"
	"shopnexus-go-service/internal/grpc/handler/product"
	"shopnexus-go-service/internal/grpc/interceptor/auth"
	"shopnexus-go-service/internal/grpc/interceptor/permission"
	"shopnexus-go-service/internal/repository"
	"shopnexus-go-service/internal/service"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/account/v1/accountv1connect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/payment/v1/paymentv1connect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/product/v1/productv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	interceptors []connect.Interceptor
	Mux          *http.ServeMux
	Services     *service.Services
}

func NewServer(address string) (*Server, error) {
	pgpool, err := pgxutil.GetPgxPool()
	if err != nil {
		return nil, err
	}

	repo := repository.NewRepository(pgpool)
	services := service.NewServices(repo)

	mux := http.NewServeMux()

	s := &Server{
		Mux:      mux,
		Services: services,
	}
	s.Init()

	return s, nil
}

func (s *Server) Init() {
	s.RegisterInterceptors()
	s.RegisterHandlers()
	s.RegisterReflection()
}

func (s *Server) Start(port int) {
	err := http.ListenAndServe(
		fmt.Sprintf(":%d", port),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(s.Mux, &http2.Server{}),
	)
	if err != nil {
		log.Fatal(err)
	}
}

// RegisterInterceptors registers the interceptors for the server
func (s *Server) RegisterInterceptors() {
	s.interceptors = append(
		s.interceptors,
		auth.NewAuthInterceptor(),
		permission.NewPermissionInterceptor(s.Services.Account, permissionRoutes),
	)
}

// RegisterHandlers registers the handlers for the server
func (s *Server) RegisterHandlers() {
	// Account
	accountPath, accountHandler := accountv1connect.NewAccountServiceHandler(
		account.NewAccountServiceHandler(s.Services.Account),
		connect.WithInterceptors(s.interceptors...),
	)
	s.Mux.Handle(accountPath, accountHandler)

	// Product
	productPath, productHandler := productv1connect.NewProductServiceHandler(
		product.NewProductServiceHandler(s.Services.Product),
		connect.WithInterceptors(s.interceptors...),
	)
	s.Mux.Handle(productPath, productHandler)

	// Payment
	paymentPath, paymentHandler := paymentv1connect.NewPaymentServiceHandler(
		payment.NewPaymentServiceHandler(s.Services.Payment),
		connect.WithInterceptors(s.interceptors...),
	)
	s.Mux.Handle(paymentPath, paymentHandler)
}

// RegisterReflection registers the reflection for the server
func (s *Server) RegisterReflection() {
	reflector := grpcreflect.NewStaticReflector(
		accountv1connect.AccountServiceName,
		paymentv1connect.PaymentServiceName,
		productv1connect.ProductServiceName,
	)
	s.Mux.Handle(grpcreflect.NewHandlerV1(reflector))
	s.Mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
}
