package middleware

import "net/http"

// SecurityHeaders adds essential HTTP security headers to every response.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent MIME-type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// Disable XSS auditor (modern browsers don't need it, can cause issues)
		w.Header().Set("X-XSS-Protection", "0")

		// Force HTTPS (1 year, include subdomains)
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Control referrer information
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Restrict browser features
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		// Prevent caching of authenticated responses
		if r.Header.Get("Authorization") != "" {
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
			w.Header().Set("Pragma", "no-cache")
		}

		next.ServeHTTP(w, r)
	})
}
