package handler

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/hassad/boilerplateSaaS/backend/internal/app/auth"
	"github.com/hassad/boilerplateSaaS/backend/internal/app/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/middleware"
)

var startTime = time.Now()

func NewRouter(
	authSvc *auth.Service,
	userSvc *user.Service,
	jwtSecret string,
	db *sql.DB,
	logger *slog.Logger,
) http.Handler {
	r := chi.NewRouter()

	// Rate limiters
	apiLimiter := middleware.NewRateLimiter(100)  // 100 req/min for API
	authLimiter := middleware.NewRateLimiter(10)   // 10 req/min for auth

	// Global middleware
	r.Use(middleware.StructuredLogging(logger))
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit(apiLimiter))

	// Health check
	r.Get("/health", healthHandler(db))

	// Public routes
	r.Route("/auth", func(r chi.Router) {
		r.Use(middleware.RateLimit(authLimiter))

		authHandler := NewAuthHandler(authSvc)
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(jwtSecret))

		userHandler := NewUserHandler(userSvc)
		r.Get("/users/me", userHandler.GetProfile)
		r.Patch("/users/me", userHandler.UpdateProfile)
		r.Put("/users/me/password", userHandler.ChangePassword)
		r.Delete("/users/me", userHandler.DeleteAccount)
	})

	return r
}

func healthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbStatus := "connected"
		if err := db.PingContext(r.Context()); err != nil {
			dbStatus = "disconnected"
		}

		uptime := time.Since(startTime).Round(time.Second).String()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"db":      dbStatus,
			"uptime":  uptime,
			"version": "1.0.0",
		})
	}
}
