package server

import (
	"fmt"
	"log"
	"net/http"
	pgxutil "shopnexus-go-service/internal/db/pgx"
	"shopnexus-go-service/internal/repository"
	"shopnexus-go-service/internal/server/connect"
	serverHttp "shopnexus-go-service/internal/server/http"
	"shopnexus-go-service/internal/server/tus"
	"shopnexus-go-service/internal/service"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type serverImpl struct {
	mux      *http.ServeMux
	services *service.Services
}

func NewServer(address string) (*serverImpl, error) {
	pgpool, err := pgxutil.GetPgxPool()
	if err != nil {
		return nil, err
	}

	repo := repository.NewRepository(pgpool)
	services := service.NewServices(repo)

	mux := http.NewServeMux()

	return &serverImpl{
		mux:      mux,
		services: services,
	}, nil
}

func (s *serverImpl) Start(port int) {
	var err error

	// Register server components
	if err = connect.Init(s.mux, s.services, true); err != nil {
		log.Fatalf("failed to initialize connect server: %v", err)
	}

	if err = tus.Init(s.mux); err != nil {
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
