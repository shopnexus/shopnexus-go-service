package server

import (
	"context"
	"net"
	"shopnexus-go-service/gen/pb"
	pgxutil "shopnexus-go-service/internal/db/pgx"
	"shopnexus-go-service/internal/interceptor"
	"shopnexus-go-service/internal/logger"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/repository"
	"shopnexus-go-service/internal/service"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"google.golang.org/grpc"
)

func NewServer(address string) error {
	// Setup auth matcher.
	matchFunc := func(ctx context.Context, callMeta interceptors.CallMeta) bool {
		return pb.Account_ServiceDesc.ServiceName != callMeta.Service
		// return pb.Product_ServiceDesc.ServiceName == callMeta.Service
	}

	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			selector.UnaryServerInterceptor(
				auth.UnaryServerInterceptor(interceptor.Auth),
				selector.MatchFunc(matchFunc),
			),
		),
		// grpc.ChainStreamInterceptor(
		// 	selector.StreamServerInterceptor(auth.StreamServerInterceptor(interceptor.Auth), selector.MatchFunc(matchFunc)),
		// ),
	)

	db := pgxutil.GetPgxPool()
	repo := repository.NewRepository(db)
	services := service.NewServices(repo)

	pb.RegisterProductServer(grpcSrv, NewProductServer(services.Product))
	pb.RegisterAccountServer(grpcSrv, NewAccountServer(services.Account))

	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	logger.Log.Info("starting gRPC server " + l.Addr().String())
	return grpcSrv.Serve(l)
}

func modelToPaginationResponse[T any](data model.PaginateResult[T]) *pb.PaginationResponse {
	return &pb.PaginationResponse{
		Total:      data.Total,
		NextPage:   data.NextPage,
		NextCursor: data.NextCursor,
	}
}
