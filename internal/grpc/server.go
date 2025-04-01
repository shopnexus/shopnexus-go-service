package grpc

import (
	"fmt"
	"log"
	"net/http"
	"shopnexus-go-service/config"
	pgxutil "shopnexus-go-service/internal/db/pgx"
	"shopnexus-go-service/internal/grpc/handler/account"
	"shopnexus-go-service/internal/grpc/handler/file"
	"shopnexus-go-service/internal/grpc/handler/payment"
	"shopnexus-go-service/internal/grpc/handler/product"
	"shopnexus-go-service/internal/grpc/interceptor/permission"
	"shopnexus-go-service/internal/logger"
	"shopnexus-go-service/internal/repository"
	"shopnexus-go-service/internal/service"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/account/v1/accountv1connect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/file/v1/filev1connect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/payment/v1/paymentv1connect"
	"github.com/shopnexus/shopnexus-protobuf-gen-go/pb/product/v1/productv1connect"
	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	tusd "github.com/tus/tusd/v2/pkg/handler"
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
	s.RegisterTusd()
}

func (s *Server) Start(port int) {
	// Create a CORS middleware handler
	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Connect-Protocol-Version, Upload-Length, Upload-Offset, Tus-Resumable, Upload-Metadata, Connect-Protocol-Version, Tus-Version, Tus-Max-Size, Tus-Extension, X-HTTP-Method-Override, X-Requested-With")
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
		// auth.NewAuthInterceptor(),
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

func (s *Server) RegisterTusd() {
	// store := s3store.New(config.GetConfig().S3.Bucket, s.Services.S3.Client)
	store := filestore.New("./uploads")

	// A locking mechanism helps preventing data loss or corruption from
	// parallel requests to a upload resource. A good match for the disk-based
	// storage is the filelocker package which uses disk-based file lock for
	// coordinating access.
	// More information is available at https://tus.github.io/tusd/advanced-topics/locks/.
	locker := filelocker.New("./uploads")

	// A storage backend for tusd may consist of multiple different parts which
	// handle upload creation, locking, termination and so on. The composer is a
	// place where all those separated pieces are joined together. In this example
	// we only use the file store but you may plug in multiple.
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)
	locker.UseIn(composer)

	// Create a new HTTP handler for the tusd server by providing a configuration.
	// The StoreComposer property must be set to allow the handler to function.
	logger.Log.Info(fmt.Sprintf("Tus url: %s", config.GetConfig().App.TusUrl))
	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              config.GetConfig().App.TusUrl,
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
		// Cors: &tusd.CorsConfig{
		// 	Disable:          false,
		// 	AllowOrigin:      regexp.MustCompile(".*"),
		// 	AllowMethods:     "GET, POST, PUT, DELETE, HEAD, PATCH, OPTIONS",
		// 	AllowHeaders:     "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Upload-Length, Upload-Offset, Tus-Resumable, Upload-Metadata, Connect-Protocol-Version, Tus-Version, Tus-Max-Size, Tus-Extension, X-HTTP-Method-Override, X-Requested-With",
		// 	ExposeHeaders:    "Upload-Offset, Location, Upload-Length, Tus-Version, Tus-Resumable, Tus-Max-Size, Tus-Extension, Upload-Metadata",
		// 	MaxAge:           "3600",
		// 	AllowCredentials: true,
		// },
	})
	if err != nil {
		log.Fatalf("unable to create handler: %s", err)
	}

	// Start another goroutine for receiving events from the handler whenever
	// an upload is completed. The event will contains details about the upload
	// itself and the relevant HTTP request.
	go func() {
		for {
			event := <-handler.CompleteUploads
			log.Printf("Upload %s finished\n", event.Upload.ID)
		}
	}()

	// Right now, nothing has happened since we need to start the HTTP server on
	// our own. In the end, tusd will start listening on and accept request at
	// http://localhost:8080/files
	// http.Handle("/files/", http.StripPrefix("/files/", handler))
	// http.Handle("/files", http.StripPrefix("/files", handler))
	// err = http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	log.Fatalf("unable to listen: %s", err)
	// }
	s.Mux.Handle("/files/", http.StripPrefix("/files/", handler))
}
