package interceptor

import (
	"context"
	"fmt"
	"shopnexus-go-service/internal/model"

	"google.golang.org/grpc"
)

func PermissionAuth(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	//! Must be called after TokenAuth
	user, ok := ctx.Value(CtxKeyUser).(model.Claims)
	if !ok {
		return nil, model.ErrTokenInvalid
	}

	fmt.Println("User:", user)

	// Proceed with the request
	return handler(ctx, req)
}
