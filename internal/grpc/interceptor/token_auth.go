package interceptor

import (
	"context"
	"shopnexus-go-service/internal/util"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
