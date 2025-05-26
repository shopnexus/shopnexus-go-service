package server

import (
	"fmt"
	"log"
	"net/http"
	"shopnexus-go-service/internal/service"
	"shopnexus-go-service/internal/transport/connect"
	serverHttp "shopnexus-go-service/internal/transport/http"
	"shopnexus-go-service/internal/transport/tus"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type server struct {
	mux      *http.ServeMux
	services *service.Services
}

func NewServer(address string) (*server, error) {
	services, err := service.NewServices()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	return &server{
		mux:      mux,
		services: services,
	}, nil
}

func (s *server) Start(port int) {
	var err error

	// Register server components
	if err = connect.Init(s.mux, s.services, true); err != nil {
		log.Fatalf("failed to initialize connect server: %v", err)
	}

	if err = tus.Init(s.mux, s.services.S3.BaseClient()); err != nil {
		log.Fatalf("failed to initialize tus server: %v", err)
	}

	if err = serverHttp.Init(s.mux, s.services); err != nil {
		log.Fatalf("failed to initialize http server: %v", err)
	}

	// Create a CORS middleware handler
	if err = http.ListenAndServe(
		fmt.Sprintf(":%d", port),
		CorsHandler(h2c.NewHandler(s.mux, &http2.Server{})),
	); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
