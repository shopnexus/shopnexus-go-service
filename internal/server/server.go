package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	pgxutil "shopnexus-go-service/internal/db/pgx"
	"shopnexus-go-service/internal/logger"
	"shopnexus-go-service/internal/repository"
	"shopnexus-go-service/internal/server/connect"
	serverHttp "shopnexus-go-service/internal/server/http"
	"shopnexus-go-service/internal/server/tus"
	"shopnexus-go-service/internal/service"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type server struct {
	mux      *http.ServeMux
	services *service.Services
}

func NewServer(address string) (*server, error) {
	// Initialize the database connection pool
	pgpool, err := pgxutil.GetPgxPool()
	if err != nil {
		return nil, err
	}

	// Check if the connection pool is valid
	if err = pgpool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	logger.Log.Info("Connected to PostgreSQL database")

	repo := repository.NewRepository(pgpool)
	services := service.NewServices(repo)

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

	if err = tus.Init(s.mux, s.services.S3.Client); err != nil {
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
