package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	adaptgemini "github.com/hassad/boilerplateSaaS/backend/internal/adapter/gemini"
	"github.com/hassad/boilerplateSaaS/backend/internal/adapter/postgres"
	adaptr2 "github.com/hassad/boilerplateSaaS/backend/internal/adapter/r2"
	"github.com/hassad/boilerplateSaaS/backend/internal/adapter/resend"
	adaptstripe "github.com/hassad/boilerplateSaaS/backend/internal/adapter/stripe"
	appai "github.com/hassad/boilerplateSaaS/backend/internal/app/ai"
	appauth "github.com/hassad/boilerplateSaaS/backend/internal/app/auth"
	appbilling "github.com/hassad/boilerplateSaaS/backend/internal/app/billing"
	appblog "github.com/hassad/boilerplateSaaS/backend/internal/app/blog"
	appnotif "github.com/hassad/boilerplateSaaS/backend/internal/app/notification"
	appstorage "github.com/hassad/boilerplateSaaS/backend/internal/app/storage"
	appteam "github.com/hassad/boilerplateSaaS/backend/internal/app/team"
	appuser "github.com/hassad/boilerplateSaaS/backend/internal/app/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/config"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler"
	"github.com/hassad/boilerplateSaaS/backend/internal/port/service"
	"github.com/hassad/boilerplateSaaS/backend/pkg/jobs"
	"github.com/hassad/boilerplateSaaS/backend/pkg/jwt"
	"github.com/hassad/boilerplateSaaS/backend/pkg/ws"
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
		processedEventRepo := postgres.NewProcessedEventRepository(db)
		paymentSvc := adaptstripe.NewPaymentService(cfg.StripeKey, cfg.StripeWebhookSecret)
		billingSvc = appbilling.NewService(userRepo, subscriptionRepo, planRepo, invoiceRepo, processedEventRepo, paymentSvc, cfg.FrontendURL)
	}

	// Storage (optional — only if R2 keys set)
	var storageSvc *appstorage.Service
	if cfg.R2AccessKey != "" {
		r2Client := adaptr2.NewClient(cfg.R2AccountID, cfg.R2AccessKey, cfg.R2SecretKey)
		r2Storage := adaptr2.NewStorageService(r2Client, cfg.R2BucketName, cfg.R2PublicURL)
		fileRepo := postgres.NewFileRepository(db)
		storageSvc = appstorage.NewService(r2Storage, fileRepo)
	}

	// AI Chat (optional — only if Gemini key set)
	var aiSvc *appai.Service
	var demoAI service.AIService
	if cfg.GeminiKey != "" {
		geminiClient, err := adaptgemini.NewClient(context.Background(), cfg.GeminiKey)
		if err != nil {
			logger.Error("failed to create Gemini client", slog.String("error", err.Error()))
		} else {
			geminiAI := adaptgemini.NewAIService(geminiClient)
			demoAI = geminiAI // expose for the public demo endpoint
			conversationRepo := postgres.NewConversationRepository(db)
			aiSvc = appai.NewService(conversationRepo, geminiAI)
		}
	}

	// Notifications
	notificationRepo := postgres.NewNotificationRepository(db)
	notifSvc := appnotif.NewService(notificationRepo)

	// WebSocket hub (real-time communication)
	wsHub := ws.NewHub()
	go wsHub.Run()

	// Wire WebSocket broadcaster into notification service
	notifSvc.SetBroadcaster(wsHub)

	// Blog
	blogRepo := postgres.NewBlogRepository(db)
	blogSvc := appblog.NewService(blogRepo)

	// Teams (optional)
	teamRepo := postgres.NewTeamRepository(db)
	memberRepo := postgres.NewTeamMemberRepository(db)
	txManager := postgres.NewTxManager(db)
	teamSvc := appteam.NewService(teamRepo, memberRepo, txManager)

	// App services
	authSvc := appauth.NewService(userRepo, passwordResetRepo, emailVerificationRepo, emailSvc, jwtMaker, cfg.FrontendURL)
	userSvc := appuser.NewService(userRepo)

	// Router
	router := handler.NewRouter(authSvc, userSvc, billingSvc, storageSvc, aiSvc, notifSvc, blogSvc, teamSvc, wsHub, cfg.JWTSecret, cfg.FrontendURL, db, logger, demoAI)

	// HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 120 * time.Second, // Increased for SSE streaming endpoints
		IdleTimeout:  60 * time.Second,
	}

	// Background job scheduler
	scheduler := jobs.NewScheduler(logger)
	scheduler.Register(jobs.Job{
		Name:     "clean-expired-password-resets",
		Interval: 1 * time.Hour,
		Fn: func(ctx context.Context) error {
			return passwordResetRepo.DeleteExpired(ctx)
		},
	})
	scheduler.Register(jobs.Job{
		Name:     "clean-expired-email-verifications",
		Interval: 1 * time.Hour,
		Fn: func(ctx context.Context) error {
			return emailVerificationRepo.DeleteExpired(ctx)
		},
	})
	scheduler.Register(jobs.Job{
		Name:     "log-system-stats",
		Interval: 5 * time.Minute,
		Fn: func(_ context.Context) error {
			stats := db.Stats()
			logger.Info("system stats",
				slog.Int("goroutines", runtime.NumGoroutine()),
				slog.Int("db_open_connections", stats.OpenConnections),
				slog.Int("db_in_use", stats.InUse),
				slog.Int("db_idle", stats.Idle),
			)
			return nil
		},
	})
	scheduler.Start(context.Background())

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

	scheduler.Stop()
	wsHub.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", slog.String("error", err.Error()))
	}

	db.Close()
	logger.Info("server stopped gracefully")
}
