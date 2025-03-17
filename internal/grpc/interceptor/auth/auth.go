package auth

import (
	"context"
	"errors"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/util"
	"strings"

	"connectrpc.com/connect"
)

type (
	ctxServerKey string
	ctxClientKey string
)

const (
	tokenHeader                   = "authorization"
	CtxServerAccount ctxServerKey = "server-account" // Storing model.Claims in context
	CtxToken         ctxClientKey = "client-account" // Storing token in context
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
				found := false
				for _, method := range methods {
					if req.Spec().Procedure == method {
						found = true
						break
					}
				}
				if !found {
					return next(ctx, req)
				}
			}

			if req.Spec().IsClient {
				// Send a token with client requests.
				token, ok := ctx.Value(CtxToken).(string)
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
					return next(ctx, req)
					// return nil, connect.NewError(
					// 	connect.CodeUnauthenticated,
					// 	errors.New("no token provided"),
					// )
				}

				// Check token in headers.
				token := strings.TrimPrefix(req.Header().Get(tokenHeader), "Bearer ")
				accountClaim, err := util.ValidateAccessToken(token)
				if err != nil {
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						err,
					)
				}

				ctx = context.WithValue(ctx, CtxServerAccount, accountClaim)
			}
			return next(ctx, req)
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}

func GetAccountFromContext(ctx context.Context) (model.Claims, error) {
	account, ok := ctx.Value(CtxServerAccount).(model.Claims)
	if !ok {
		return model.Claims{}, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}
	return account, nil
}
