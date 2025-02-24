package interceptor

import (
	"context"
	"fmt"
	"shopnexus-go-service/internal/util"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// TokenAuth is a server interceptor for authentication
func TokenAuth(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	// Extract metadata from incoming request
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	// Get the token from the "authorization" header
	authHeader, exists := md["authorization"]
	if !exists || len(authHeader) == 0 {
		return nil, fmt.Errorf("authorization token required")
	}

	tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")

	claims, err := util.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, CtxKeyUser, claims)

	// Proceed with the request
	return handler(ctx, req)
}

func Auth(ctx context.Context) (context.Context, error) {
	token, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	claims, err := util.ValidateAccessToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid auth token")
	}

	ctx = context.WithValue(ctx, CtxKeyUser, claims)

	return ctx, nil
}
