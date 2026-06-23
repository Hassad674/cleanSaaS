package handler

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	appai "github.com/hassad/boilerplateSaaS/backend/internal/app/ai"
	"github.com/hassad/boilerplateSaaS/backend/internal/app/auth"
	appbilling "github.com/hassad/boilerplateSaaS/backend/internal/app/billing"
	appblog "github.com/hassad/boilerplateSaaS/backend/internal/app/blog"
	appnotif "github.com/hassad/boilerplateSaaS/backend/internal/app/notification"
	appstorage "github.com/hassad/boilerplateSaaS/backend/internal/app/storage"
	appteam "github.com/hassad/boilerplateSaaS/backend/internal/app/team"
	"github.com/hassad/boilerplateSaaS/backend/internal/app/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/middleware"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
	"github.com/hassad/boilerplateSaaS/backend/pkg/jwt"
	"github.com/hassad/boilerplateSaaS/backend/pkg/ws"
)

var startTime = time.Now()

func NewRouter(
	authSvc *auth.Service,
	userSvc *user.Service,
	billingSvc *appbilling.Service,
	storageSvc *appstorage.Service,
	aiSvc *appai.Service,
	notifSvc *appnotif.Service,
	blogSvc *appblog.Service,
	teamSvc *appteam.Service,
	wsHub *ws.Hub,
	jwtSecret string,
	frontendURL string,
	db *sql.DB,
	logger *slog.Logger,
	demoAI service.AIService,
) http.Handler {
	r := chi.NewRouter()

	// Rate limiters
	apiLimiter := middleware.NewRateLimiter(100) // 100 req/min for API
	authLimiter := middleware.NewRateLimiter(10) // 10 req/min for auth

	// Global middleware
	r.Use(middleware.StructuredLogging(logger))
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.CORS(frontendURL))
	r.Use(middleware.MaxBodySize(1 << 20)) // 1MB default for JSON endpoints
	r.Use(middleware.RateLimit(apiLimiter))

	// Health check
	r.Get("/health", healthHandler(db))

	// WebSocket endpoint (optional — only if hub is provided)
	if wsHub != nil {
		wsHandler := NewWSHandler(wsHub, jwt.NewMaker(jwtSecret), frontendURL)
		r.Get("/ws", wsHandler.Upgrade)
	}

	// Public routes
	r.Route("/auth", func(r chi.Router) {
		r.Use(middleware.RateLimit(authLimiter))

		authHandler := NewAuthHandler(authSvc)
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)
		r.Post("/logout", authHandler.Logout)
		r.Post("/forgot-password", authHandler.ForgotPassword)
		r.Post("/reset-password", authHandler.ResetPassword)
		r.Post("/verify-email", authHandler.VerifyEmail)
	})

	// Public billing routes
	if billingSvc != nil {
		billingHandler := NewBillingHandler(billingSvc)
		r.Get("/billing/plans", billingHandler.GetPlans)
		r.Post("/webhooks/stripe", billingHandler.HandleWebhook)

		// Public demo billing (no auth required, real Stripe in test mode)
		demoBillingLimiter := middleware.NewRateLimiter(20)
		r.Route("/demo/billing", func(r chi.Router) {
			r.Use(middleware.RateLimit(demoBillingLimiter))
			r.Post("/checkout", billingHandler.DemoCheckout)
			r.Get("/session", billingHandler.DemoSession)
			r.Post("/portal", billingHandler.DemoPortal)
		})
	}

	// Public blog routes
	blogHandler := NewBlogHandler(blogSvc)
	r.Get("/blog/posts", blogHandler.ListPublished)
	r.Get("/blog/posts/{slug}", blogHandler.GetBySlug)
	r.Get("/blog/tags", blogHandler.ListTags)

	// Public demo AI chat (no auth required, no DB persistence)
	if demoAI != nil {
		demoHandler := NewDemoHandler(demoAI)
		demoLimiter := middleware.NewRateLimiter(20) // 20 req/min for demo
		r.Route("/demo/ai", func(r chi.Router) {
			r.Use(middleware.RateLimit(demoLimiter))
			r.Post("/chat", demoHandler.StreamChat)
		})
	}

	// Public demo storage (no auth required, uses fixed demo user ID)
	if storageSvc != nil {
		demoStorageHandler := NewDemoStorageHandler(storageSvc)
		demoStorageLimiter := middleware.NewRateLimiter(20) // 20 req/min for demo
		r.Route("/demo/storage", func(r chi.Router) {
			r.Use(middleware.RateLimit(demoStorageLimiter))
			r.With(middleware.MaxBodySize(10<<20)).Post("/upload", demoStorageHandler.Upload)
			r.Get("/files", demoStorageHandler.List)
			r.Delete("/files/{id}", demoStorageHandler.Delete)
		})
	}

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(jwtSecret))

		// Auth actions requiring authentication
		authHandler := NewAuthHandler(authSvc)
		r.Post("/auth/resend-verification", authHandler.ResendVerification)

		userHandler := NewUserHandler(userSvc)
		r.Get("/users/me", userHandler.GetProfile)
		r.Patch("/users/me", userHandler.UpdateProfile)
		r.Put("/users/me/password", userHandler.ChangePassword)
		r.Delete("/users/me", userHandler.DeleteAccount)

		// Billing (authenticated)
		if billingSvc != nil {
			billingHandler := NewBillingHandler(billingSvc)
			r.Post("/billing/checkout", billingHandler.CreateCheckout)
			r.Get("/billing/subscription", billingHandler.GetSubscription)
			r.Post("/billing/cancel", billingHandler.CancelSubscription)
			r.Post("/billing/portal", billingHandler.CreatePortalSession)
			r.Get("/billing/invoices", billingHandler.GetInvoices)
		}

		// Storage (authenticated)
		if storageSvc != nil {
			storageHandler := NewStorageHandler(storageSvc)
			r.With(middleware.MaxBodySize(50<<20)).Post("/files/upload", storageHandler.Upload)
			r.Get("/files", storageHandler.List)
			r.Delete("/files/{id}", storageHandler.Delete)
		}

		// Notifications (always available)
		notifHandler := NewNotificationHandler(notifSvc)
		r.Get("/notifications", notifHandler.List)
		r.Get("/notifications/count", notifHandler.UnreadCount)
		r.Put("/notifications/{id}/read", notifHandler.MarkRead)
		r.Put("/notifications/read-all", notifHandler.MarkAllRead)

		// AI Chat (authenticated)
		if aiSvc != nil {
			aiHandler := NewAIHandler(aiSvc)
			r.Get("/ai/conversations", aiHandler.ListConversations)
			r.Post("/ai/conversations", aiHandler.CreateConversation)
			r.Get("/ai/conversations/{id}/messages", aiHandler.GetMessages)
			r.Post("/ai/conversations/{id}/messages", aiHandler.SendMessage)
			r.Post("/ai/conversations/{id}/stream", aiHandler.StreamMessage)
			r.Delete("/ai/conversations/{id}", aiHandler.DeleteConversation)
		}

		// Teams (optional — only if team service is provided)
		if teamSvc != nil {
			teamHandler := NewTeamHandler(teamSvc)
			r.Post("/teams", teamHandler.Create)
			r.Get("/teams", teamHandler.List)
			r.Post("/teams/invite/accept", teamHandler.AcceptInvite)
			r.Post("/teams/invite/decline", teamHandler.DeclineInvite)
			r.Route("/teams/{id}", func(r chi.Router) {
				r.Get("/", teamHandler.Get)
				r.Put("/", teamHandler.Update)
				r.Delete("/", teamHandler.Delete)
				r.Post("/invite", teamHandler.InviteMember)
				r.Get("/members", teamHandler.ListMembers)
				r.Put("/members/{userId}/role", teamHandler.UpdateMemberRole)
				r.Delete("/members/{userId}", teamHandler.RemoveMember)
				r.Post("/leave", teamHandler.Leave)
			})
		}

		// Admin routes (require admin role)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAdmin)

			adminHandler := NewAdminHandler(userSvc, blogSvc)
			r.Get("/admin/stats", adminHandler.DashboardStats)
			r.Get("/admin/users", adminHandler.ListUsers)
			r.Put("/admin/users/{id}/role", adminHandler.UpdateUserRole)

			r.Get("/admin/blog/posts", blogHandler.AdminList)
			r.Post("/admin/blog/posts", blogHandler.Create)
			r.Put("/admin/blog/posts/{id}", blogHandler.Update)
			r.Delete("/admin/blog/posts/{id}", blogHandler.Delete)
		})
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
