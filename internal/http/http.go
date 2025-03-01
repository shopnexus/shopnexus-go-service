package http

import (
	"net/http"
	"shopnexus-go-service/gen/pb"
	httpmiddleware "shopnexus-go-service/internal/http/middleware"
	"shopnexus-go-service/internal/http/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	router     *chi.Mux
	grpcClient *grpc.ClientConn
	product    pb.ProductClient
	account    pb.AccountClient
	cart       pb.CartClient
	payment    pb.PaymentClient
	refund     pb.RefundClient
}

var decode = schema.NewDecoder()

func NewServer(grpcAddr string) (*Server, error) {
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	s := &Server{
		router:     chi.NewRouter(),
		grpcClient: conn,
		product:    pb.NewProductClient(conn),
		account:    pb.NewAccountClient(conn),
		cart:       pb.NewCartClient(conn),
		payment:    pb.NewPaymentClient(conn),
		refund:     pb.NewRefundClient(conn),
	}

	s.setupMiddlewares()
	s.setupRoutes()

	return s, nil
}

func (s *Server) setupMiddlewares() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	s.router.Use(middleware.AllowContentType("application/json"))
	s.router.Use(middleware.SetHeader("Content-Type", "application/json"))

	// Add 404 handler with logging
	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		println("404 Not Found:", r.URL.Path)
		response.FromMessage(w, http.StatusNotFound, "404 Not Found")
	})

}

func (s *Server) setupRoutes() {
	s.router.Route("/api", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Mount("/account", NewAccountHandler(s.account))
			r.Mount("/ipn", NewIPNHandler(s.payment))
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(httpmiddleware.GrpcAuthorization)
			r.Mount("/products", NewProductHandler(s.product))
			r.Mount("/cart", NewCartHandler(s.cart))
			r.Mount("/payment", NewPaymentHandler(s.payment))
			r.Mount("/refund", NewRefundHandler(s.refund))
		})
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Close() error {
	return s.grpcClient.Close()
}

func (s *Server) Start(addr string) error {
	// Log server starting
	println("HTTP server starting on", addr)
	return http.ListenAndServe(addr, s)
}
