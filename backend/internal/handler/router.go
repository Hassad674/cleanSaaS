package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/hassad/boilerplateSaaS/backend/internal/app/auth"
	"github.com/hassad/boilerplateSaaS/backend/internal/app/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/middleware"
)

func NewRouter(
	authSvc *auth.Service,
	userSvc *user.Service,
	jwtSecret string,
) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(middleware.CORS())

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Public routes
	r.Route("/auth", func(r chi.Router) {
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
		r.Delete("/users/me", userHandler.DeleteAccount)
	})

	return r
}
