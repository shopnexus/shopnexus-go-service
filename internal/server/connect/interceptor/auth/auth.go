package auth

import (
	"context"
	"errors"
	"shopnexus-go-service/internal/model"
	"shopnexus-go-service/internal/util"
	"shopnexus-go-service/pkg/cache"
	"strings"
	"time"

	"slices"

	"connectrpc.com/connect"
)

type ctxKey string

const (
	tokenHeader        = "authorization"
	CtxClaims   ctxKey = "ctx-claims" // Storing model.Claims in context
	CtxToken    ctxKey = "ctx-token"  // Storing token in context
)

var (
	claimsCache = cache.NewCache[string, model.Claims]()
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
				found := slices.Contains(methods, req.Spec().Procedure)
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
				claims, err := util.ValidateAccessToken(token)
				if err == nil {
					ctx = context.WithValue(ctx, CtxClaims, claims)
				}
			}
			return next(ctx, req)
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}

func GetClaims(req connect.AnyRequest) (claims model.Claims, err error) {
	token := req.Header().Get(tokenHeader)

	claims, ok := claimsCache.Get(token)
	if ok {
		return claims, nil
	}

	claims, err = util.ValidateAccessToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return model.Claims{}, connect.NewError(connect.CodeUnauthenticated, err)
	}

	claimsCache.Set(token, claims, 5*60*time.Second)
	return claims, nil
}
