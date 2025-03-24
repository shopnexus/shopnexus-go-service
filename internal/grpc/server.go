package grpc

import (
	"fmt"
	"log"
	"net/http"
	pgxutil "shopnexus-go-service/internal/db/pgx"
	"shopnexus-go-service/internal/grpc/handler/account"
	"shopnexus-go-service/internal/grpc/handler/file"
	"shopnexus-go-service/internal/grpc/handler/payment"
	"shopnexus-go-service/internal/grpc/handler/product"
	"shopnexus-go-service/internal/grpc/interceptor/auth"
	"shopnexus-go-service/internal/grpc/interceptor/permission"
	"shopnexus-go-service/internal/repository"
	"shopnexus-go-service/internal/service"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/account/v1/accountv1connect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/file/v1/filev1connect"
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
	// Create a CORS middleware handler
	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Connect-Protocol-Version")
			w.Header().Set("Access-Control-Max-Age", "3600")

			// Handle preflight OPTIONS requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// Pass the request to the next handler
			h.ServeHTTP(w, r)
		})
	}

	err := http.ListenAndServe(
		fmt.Sprintf(":%d", port),
		// Apply CORS middleware to the h2c handler
		corsHandler(h2c.NewHandler(s.Mux, &http2.Server{})),
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
	// File
	filePath, fileHandler := filev1connect.NewFileServiceHandler(
		file.NewFileServiceHandler(s.Services.S3),
		connect.WithInterceptors(s.interceptors...),
	)
	s.Mux.Handle(filePath, fileHandler)

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
		filev1connect.FileServiceName,
		accountv1connect.AccountServiceName,
		paymentv1connect.PaymentServiceName,
		productv1connect.ProductServiceName,
	)
	s.Mux.Handle(grpcreflect.NewHandlerV1(reflector))
	s.Mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
}
