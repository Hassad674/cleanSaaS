package middleware

import "net/http"

// MaxBodySize limits the size of incoming request bodies to prevent DoS attacks.
// size is in bytes. Defaults to 1MB if size <= 0.
func MaxBodySize(size int64) func(http.Handler) http.Handler {
	if size <= 0 {
		size = 1 << 20 // 1MB
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, size)
			next.ServeHTTP(w, r)
		})
	}
}
