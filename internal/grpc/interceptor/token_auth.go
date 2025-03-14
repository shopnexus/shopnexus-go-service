package interceptor

import (
	"context"
	"errors"
	"shopnexus-go-service/internal/util"

	"connectrpc.com/connect"
)

type (
	ctxServerKey string
	ctxClientKey string
)

const (
	tokenHeader                = "authorization"
	ctxServerUser ctxServerKey = "server-user"
	ctxToken      ctxClientKey = "client-user"
)

// NewAuthInterceptor returns a new auth interceptor.
// This interceptor checks for a token in the request headers or sends a token with client requests.
// If no token is provided, it returns an unauthenticated error.
func NewAuthInterceptor(methods ...string) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if len(methods) > 0 {
				for _, method := range methods {
					if req.Spec().Procedure == method {
						break
					}
				}
			}

			if req.Spec().IsClient {
				// Send a token with client requests.
				token, ok := ctx.Value(ctxToken).(string)
				if !ok {
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						errors.New("no token provided"),
					)
				}
				req.Header().Set(tokenHeader, "Bearer "+token)
			} else {
				// Check token in handlers.
				if req.Header().Get(tokenHeader) == "" {
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						errors.New("no token provided"),
					)
				}

				// Check token in headers.
				token := req.Header().Get(tokenHeader)
				userClaim, err := util.ValidateAccessToken(token)
				if err != nil {
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						err,
					)
				}

				ctx = context.WithValue(ctx, ctxServerUser, userClaim)
			}
			return next(ctx, req)
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}
