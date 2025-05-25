package server

import "net/http"

func CorsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Connect-Protocol-Version, Upload-Length, Upload-Offset, Tus-Resumable, Upload-Metadata, Connect-Protocol-Version, Tus-Version, Tus-Max-Size, Tus-Extension, X-HTTP-Method-Override, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Pass the request to the next handler
		h.ServeHTTP(w, r)
	})
}
