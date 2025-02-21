package interceptor

import (
	"context"
	"fmt"
	"shopnexus-go-service/internal/util"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
