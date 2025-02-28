package httpmiddleware

import (
	"net/http"
	"shopnexus-go-service/internal/http/response"

	"google.golang.org/grpc/metadata"
)

// GrpcAuthorization middleware automatically
func GrpcAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			response.FromMessage(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Forward the auth token to gRPC metadata
		ctx := metadata.AppendToOutgoingContext(r.Context(), "authorization", token)

		// Call the next handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
