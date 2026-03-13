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
	appauth "github.com/hassad/boilerplateSaaS/backend/internal/app/auth"
	appuser "github.com/hassad/boilerplateSaaS/backend/internal/app/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/config"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler"
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

	// JWT
	jwtMaker := jwt.NewMaker(cfg.JWTSecret)

	// External services
	var emailSvc *resend.EmailService
	if cfg.ResendKey != "" {
		emailSvc = resend.NewEmailService(cfg.ResendKey)
	}

	// App services
	authSvc := appauth.NewService(userRepo, emailSvc, jwtMaker)
	userSvc := appuser.NewService(userRepo)

	// Router
	router := handler.NewRouter(authSvc, userSvc, cfg.JWTSecret, db, logger)

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
