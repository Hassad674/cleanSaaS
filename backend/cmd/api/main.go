package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hassad/boilerplateSaaS/backend/internal/adapter/postgres"
	"github.com/hassad/boilerplateSaaS/backend/internal/adapter/resend"
	adaptstripe "github.com/hassad/boilerplateSaaS/backend/internal/adapter/stripe"
	appauth "github.com/hassad/boilerplateSaaS/backend/internal/app/auth"
	appbilling "github.com/hassad/boilerplateSaaS/backend/internal/app/billing"
	appuser "github.com/hassad/boilerplateSaaS/backend/internal/app/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/config"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
	"github.com/hassad/boilerplateSaaS/backend/pkg/jwt"
)

func main() {
	// Structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := config.Load()

	// Database
	db := postgres.NewDB(cfg.DatabaseURL)
	defer db.Close()

	// Repositories
	userRepo := postgres.NewUserRepository(db)
	passwordResetRepo := postgres.NewPasswordResetRepository(db)
	emailVerificationRepo := postgres.NewEmailVerificationRepository(db)

	// JWT
	jwtMaker := jwt.NewMaker(cfg.JWTSecret)

	// External services
	var emailSvc service.EmailService
	if cfg.ResendKey != "" {
		emailSvc = resend.NewEmailService(cfg.ResendKey)
	}

	// Billing repositories + service (optional — only if Stripe key set)
	var billingSvc *appbilling.Service
	if cfg.StripeKey != "" {
		subscriptionRepo := postgres.NewSubscriptionRepository(db)
		planRepo := postgres.NewPlanRepository(db)
		invoiceRepo := postgres.NewInvoiceRepository(db)
		paymentSvc := adaptstripe.NewPaymentService(cfg.StripeKey, cfg.StripeWebhookSecret)
		billingSvc = appbilling.NewService(userRepo, subscriptionRepo, planRepo, invoiceRepo, paymentSvc, cfg.FrontendURL)
	}

	// App services
	authSvc := appauth.NewService(userRepo, passwordResetRepo, emailVerificationRepo, emailSvc, jwtMaker, cfg.FrontendURL)
	userSvc := appuser.NewService(userRepo)

	// Router
	router := handler.NewRouter(authSvc, userSvc, billingSvc, cfg.JWTSecret, db, logger)

	// HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("API server starting", slog.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info("shutdown signal received", slog.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", slog.String("error", err.Error()))
	}

	db.Close()
	logger.Info("server stopped gracefully")
}
