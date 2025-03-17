package interceptor

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
)

func NewCacheInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			// Check if the method is in the list of methods to cache
			res, err := next(ctx, req)
			if err != nil {
				return res, err
			}

			if req.HTTPMethod() == http.MethodGet {
				res.Header().Set("Cache-Control", "max-age=604800")
			}

			return res, nil
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}
